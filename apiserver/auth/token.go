package auth

import (
	"encoding/csv"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/zhaozf-zhiming/suneee/apiserver/common/log"
	"github.com/zhaozf-zhiming/suneee/apiserver/etc/apiconfig"
	"github.com/zhaozf-zhiming/suneee/apiserver/handler"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

var (
	ErrTokenNotFound     = errors.New("请求未携带token，无权限访问")
	ErrTokenInvalid      = errors.New("无效的token")
	ErrTokenCanNotAccess = errors.New("无权访问API分组")
	TokenAuthMap         = make(map[TokenContent]TokenCsv)
)

type TokenContent string
type TokenLogin struct {
	Content TokenContent
}
type TokenCsv struct {
	Content  TokenContent
	UserName string
	UserId   string
	ApiGrp   []string
}
type TokenAuth struct {
	Map map[TokenContent]TokenCsv
}

func (t *TokenAuth) CanAccess(token TokenLogin, grp string) bool {
	csv, err := t.GetTokenCsv(token)
	if err != nil {
		log.Logger.Error(err.Error())
		return false
	}
	for _, v := range csv.ApiGrp {
		if v == grp {
			return true
		}
	}
	return false
}

func (t *TokenAuth) IsTokenExist(token TokenLogin) bool {
	_, ok := t.Map[token.Content]
	return ok
}

func (t *TokenAuth) GetTokenCsv(token TokenLogin) (*TokenCsv, error) {
	csv, ok := t.Map[token.Content]
	if !ok {
		return nil, ErrTokenNotFound
	}
	return &csv, nil
}

func init() {
	file, err := ioutil.ReadFile(filepath.Join(apiconfig.GetServerDir(), apiconfig.GetTokenAuthPath()))
	if err != nil {
		log.Logger.Fatal(err.Error())
	}
	reader := csv.NewReader(strings.NewReader(string(file)))
	rows, _ := reader.ReadAll()
	for _, v := range rows {
		if len(v) != 4 {
			log.Logger.Fatal("未知的csv行", zap.Strings("csv", v))
		}
		tokenCsv := TokenCsv{
			Content:  TokenContent(v[0]),
			UserName: v[1],
			UserId:   v[2],
			ApiGrp:   strings.Split(v[3], ","),
		}
		TokenAuthMap[TokenContent(v[0])] = tokenCsv
	}
}

func NewTokenAuth() *TokenAuth {
	auth := new(TokenAuth)
	auth.Map = TokenAuthMap
	return auth
}

func (t *TokenAuth) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/swagger") {
			c.Next()
			return
		}
		content := c.Request.Header.Get("token")
		if content == "" {
			log.Logger.Error(ErrTokenNotFound.Error(), zap.String("url", c.Request.URL.String()))
			c.JSON(http.StatusUnauthorized, gin.H{"code": handler.FailCode, "msg": ErrTokenNotFound.Error()})
			c.Abort()
			return
		}
		token := TokenLogin{Content: TokenContent(content)}
		if !t.IsTokenExist(token) {
			c.JSON(http.StatusUnauthorized, gin.H{"code": handler.FailCode, "msg": ErrTokenInvalid.Error()})
			log.Logger.Error(ErrTokenInvalid.Error(), zap.String("url", c.Request.URL.String()))
			c.Abort()
			return
		}
		grp := GetApiGrpFromContext(c)
		log.Sugar.Debugf("path => %s, grp => %s", c.Request.URL.Path, grp)
		if !t.CanAccess(token, grp) {
			c.JSON(http.StatusForbidden, gin.H{"code": handler.FailCode, "msg": ErrTokenCanNotAccess.Error()})
			log.Logger.Error(ErrTokenCanNotAccess.Error(), zap.String("url", c.Request.URL.String()))
			c.Abort()
			return
		}
		csv, err := t.GetTokenCsv(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": handler.FailCode, "msg": err.Error()})
			log.Logger.Error(err.Error(), zap.String("url", c.Request.URL.String()))
			c.Abort()
			return
		}
		// 继续交由下一个路由处理,并将解析出的信息传递下去
		c.Set("token_csv", csv)
	}
}

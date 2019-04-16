package auth

import (
	"../etc/apiconfig"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

func GetApiGrpFromContext(c *gin.Context) string {
	path := c.Request.URL.Path
	var prefix string
	withNamespacePrefix := fmt.Sprintf("/api/%s/%s/", apiconfig.GetApiDefaultVersion(), "namespace")
	withoutNamespacePrefix := fmt.Sprintf("/api/%s/", apiconfig.GetApiDefaultVersion())

	if strings.HasPrefix(path, withNamespacePrefix) {
		prefix = withNamespacePrefix
		pathCutPrefix := path[strings.Index(path, prefix)+len(prefix):]
		grp := pathCutPrefix[strings.Index(pathCutPrefix, "/")+1:]
		if strings.Index(grp, "/") >= 0 {
			grp = grp[:strings.Index(grp, "/")]
		}
		return grp
	} else {
		prefix = withoutNamespacePrefix
		pathCutPrefix := path[strings.Index(path, prefix)+len(prefix):]
		grp := pathCutPrefix
		if strings.Index(pathCutPrefix, "/") >= 0 {
			grp = pathCutPrefix[:strings.Index(pathCutPrefix, "/")]
		}
		return grp
	}
}

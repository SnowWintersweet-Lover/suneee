package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zhaozf-zhiming/suneee/apiserver/common/log"
	"github.com/zhaozf-zhiming/suneee/apiserver/common/types"
	"github.com/zhaozf-zhiming/suneee/apiserver/k8s_cli"
	"net/http"
	"strconv"
)

//查询合约列表信息
func HandlerGetDeployment(c *gin.Context) {
	var queryInfo types.QueryDeployment
	queryInfo.Namespace = c.Query("namespace")
	//if queryInfo.Namespace == "" {
	//	err := errors.New("namespace can not empty")
	//	log.Logger.Error(err.Error())
	//	c.JSON(http.StatusOK, gin.H{"status": FailCode, "reason": err.Error()})
	//	return
	//}
	queryInfo.Name = c.Query("name")
	queryInfo.Start, _ = strconv.Atoi(c.Query("start"))
	queryInfo.Limit, _ = strconv.Atoi(c.Query("limit"))

	rtVal, err := QueryDeploymentList(queryInfo)
	if err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{"status": FailCode, "reason": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": SuccCode, "data": *rtVal})
}

func QueryDeploymentList(queryInfo types.QueryDeployment) (*types.QueryOut, error) {
	return k8s_cli.QueryK8sInfo(queryInfo)
}

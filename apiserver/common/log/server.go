package log

import (
	"../../etc/apiconfig"
	"fmt"
	"net/http"
)

// 可以通过设置Atom.SetLevel(zap.DebugLevel)动态调节日志级别
// 可以通过http接口动态设置日志级别和查看当前日志级别http://localhost:9000/level
func InitLogServer() {
	http.HandleFunc("/level", Atom.ServeHTTP)
	address := fmt.Sprintf("%s:%d", apiconfig.GetLogHostAddress(), apiconfig.GetLogHostPort())
	Logger.Fatal(http.ListenAndServe(address, nil).Error())
}

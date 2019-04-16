package main

import (
	"./auth"
	"./common/log"
	"./etc/apiconfig"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func InitServer() {
	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.Use(gin.Recovery())
	// 默认设置logger，但启用logger会导致吞吐量大幅度降低
	if os.Getenv("GIN_LOG") != "off" {
		r.Use(gin.Logger())
	}
	r.MaxMultipartMemory = 10 << 20 // 10 MB
	r.GET("/ping", handler.PingHandler)
	r.Use(auth.NewTokenAuth().Middleware())
	v1Deployment := r.Group(fmt.Sprintf("/api/%s/deployment", apiconfig.GetApiDefaultVersion()))
	{
		v1Deployment.GET("", handler.HandlerGetDeployment)

	}
	s := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", apiconfig.Config.Server.Host.Address, apiconfig.Config.Server.Host.Port),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	if apiconfig.GetServerTlsEnable() {
		pool := x509.NewCertPool()
		caCertPath := apiconfig.GetServerPkiCa()
		caCrt, err := ioutil.ReadFile(caCertPath)
		if err != nil {
			log.Logger.Fatal(err.Error())
		}
		pool.AppendCertsFromPEM(caCrt)
		tlsConfig := &tls.Config{
			ClientCAs: pool,
		}
		if apiconfig.GetServerTlsVerifyPeer() {
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		} else {
			tlsConfig.ClientAuth = tls.NoClientCert
		}
		go func() {
			if err := s.ListenAndServeTLS(apiconfig.GetServerPkiCert(), apiconfig.GetServerPkiKey()); err != nil {
				log.Logger.Fatal(err.Error())
			}
		}()
		log.Sugar.Infof("apiserver[https] listening on %s", s.Addr)
	} else {
		go func() {
			if err := s.ListenAndServe(); err != nil {
				log.Logger.Fatal(err.Error())
			}
		}()
		log.Sugar.Infof("apiserver[http] listening on %s", s.Addr)
	}
	// 监听log server，动态调整log级别
	go func() {
		log.InitLogServer()
		log.Sugar.Infof("logserver[http] listening on %s", s.Addr)
	}()
	// apiserver发生错误后延时五秒钟，优雅关闭
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Sugar.Info("收到操作系统退出信号量，即将停止apiserver ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Logger.Fatal(err.Error())
	}
	log.Sugar.Info("成功停止apiserver")
}

func main() {
	InitServer()
}

package main

import (
	"context"
	"fmt"
	_ "ginpro/boot"
	"ginpro/config"
	"ginpro/internal/router"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

//var g errgroup.Group

// @title gin系统
// @version 1.0
// @description gin开发的系统
// @termsOfService
func main() {
	appCfg := config.Conf.App

	gin.SetMode(appCfg.RunMode)
	routers := router.ApiRouter()
	var s *http.Server
	go func() {
		s = &http.Server{
			Addr:           ":" + appCfg.Port,
			Handler:        routers,
			ReadTimeout:    appCfg.ReadTimeout * time.Second,
			WriteTimeout:   appCfg.WriteTimeout * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("s.ListenAndServe error: %v", err)
		}
	}()
	//等待中断信号
	quit := make(chan os.Signal)
	//接收syscall.SIGINT和syscall.SIGTERM信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		fmt.Println("Server forced shutdown", err)
	}
	fmt.Println("Server existing")
}

package main

import (
	"context"
	"fmt"
	_ "gin-api/boot"
	"gin-api/config"
	"gin-api/internal/router"
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
	cfg := config.Conf.App
	gin.SetMode(cfg.RunMode)
	routers := router.ApiRouter()
	var s *http.Server
	go func() {
		s = &http.Server{
			Addr:           ":" + cfg.Port,
			Handler:        routers,
			ReadTimeout:    cfg.ReadTimeout * time.Second,
			WriteTimeout:   cfg.WriteTimeout * time.Second,
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

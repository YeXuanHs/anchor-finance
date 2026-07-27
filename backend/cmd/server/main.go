package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/anchor-finance/backend/internal/api"
	"github.com/anchor-finance/backend/internal/config"
	"github.com/anchor-finance/backend/internal/job"
	"github.com/anchor-finance/backend/pkg/db"
	"github.com/anchor-finance/backend/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	logger.Init(cfg.Log)
	if err := db.Init(cfg.Database); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	if err := db.InitRedis(cfg.Redis); err != nil {
		log.Fatalf("Redis初始化失败: %v", err)
	}
	go job.Start(cfg.Job)
	srv := api.NewServer(cfg)
	go func() {
		addr := fmt.Sprintf(":%d", cfg.Server.Port)
		logger.Infof("锚点财务服务启动: http://localhost%s", addr)
		if err := srv.Run(addr); err != nil {
			log.Fatalf("服务启动失败: %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("服务正在关闭...")
}

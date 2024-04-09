package main

import (
	"log"

	"github.com/Grey0520/s3proxy/internal/config"
	"github.com/Grey0520/s3proxy/internal/server"
	"github.com/Grey0520/s3proxy/internal/server/routes"
)

func main() {
	// 1. 加载配置
	err := config.LoadConfig("")
	if err != nil {
		log.Fatal("failed to load configuration:", err)
	}
	log.Print(config.Cfg)

	// 2. 资源初始化
	app := server.NewServer(&config.Cfg)

	routes.ConfigureRoutes(app)
	err = app.Start(config.Cfg.S3Proxy.Endpoint)
	if err != nil {
		log.Fatal("Port already used")
	}
}

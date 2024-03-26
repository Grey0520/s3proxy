package main

import (
	"log"

	"github.com/Grey0520/s3proxy/internal/config"
)

func main() {
	// 1. 加载配置
	err := config.LoadConfig("")
	if err != nil {
		log.Fatal("failed to load configuration:", err)
	}
	log.Print(config.Cfg)

	// 2. 资源初始化

	return
}

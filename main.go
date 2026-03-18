package main

import (
	"datasources/handlers"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建Gin引擎
	r := gin.Default()

	// 设置静态文件目录
	r.Static("/static", "./static")

	// 设置路由
	r.GET("/", handlers.Index)
	r.POST("/api/datasource/test", handlers.TestConnectivity)
	r.POST("/api/datasource/save", handlers.SaveDatasource)
	r.GET("/api/datasource/list", handlers.GetDatasourceList)
	r.GET("/api/datasource/:id", handlers.GetDatasourceDetail)

	// 获取端口号，默认为8080
	port := ":8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = ":" + envPort
	}

	// 启动服务
	fmt.Printf("服务启动成功: http://localhost%s\n", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}

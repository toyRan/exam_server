package main

import (
	"exam_server/config"
	"exam_server/routes"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// 初始化日志文件
var logFile *os.File

func init() {
	var err error
	logFile, err = os.OpenFile("error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	log.SetOutput(logFile)
}

func main() {

	fmt.Println("Application started...")

	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// 初始化数据库
	config.InitDB()
	log.Println("初始化数据库成功")

	// 初始化 OSS
	config.InitOSS()
	log.Println("初始化 OSS 成功")

	// 初始化 Gin 路由
	r := gin.Default()
	log.Println("路由注册成功")

	// 配置 CORS，允许来自指定来源的请求
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"https://catalog.glassesshare.com",
			"https://admin.glassesshare.com",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 加载 HTML 模板
	//r.LoadHTMLGlob("templates/*")

	// 加载路由
	routes.RegisterRoutes(r)

	// 运行 Gin 服务
	err := r.Run(":8082")
	if err != nil {
		log.Printf("Gin 服务器启动失败: %v\n", err)
	}
}

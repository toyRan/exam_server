package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// 配置数据库连接池
	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
		// 启用连接池
		DefaultStringSize: 256,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)                  // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)                 // 最大打开连接数
	sqlDB.SetConnMaxLifetime(59 * time.Second) // 连接最大生命周期
	sqlDB.SetConnMaxIdleTime(time.Hour)        // 空闲连接最大生命周期

	// 配置自动重连
	DB = DB.Session(&gorm.Session{
		// 启用预处理语句
		PrepareStmt: true,
		// 设置操作超时
		Context: context.Background(),
	})

	// 添加连接状态检查中间件
	DB.Callback().Query().Before("gorm:query").Register("check_connection", checkConnection)
	DB.Callback().Create().Before("gorm:create").Register("check_connection", checkConnection)
	DB.Callback().Update().Before("gorm:update").Register("check_connection", checkConnection)
	DB.Callback().Delete().Before("gorm:delete").Register("check_connection", checkConnection)
}

// 连接检查中间件
func checkConnection(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("获取数据库实例失败: %v", err)
		reconnect()
		return
	}

	// 检查连接是否有效
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		log.Printf("数据库连接检查失败: %v", err)
		reconnect()
		return
	}
}

// 重连函数
func reconnect() {
	log.Println("尝试重新连接数据库...")
	InitDB()
}

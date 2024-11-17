package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func maintest() {
	// 连接数据库
	dsn := "jx_ebook:jx_ebook@tcp(127.0.0.1:3306)/jx_ebook?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	// 创建代码生成器
	g := gen.NewGenerator(gen.Config{
		OutPath: "./models", // 生成代码的输出路径
	})

	// 使用数据库对象
	g.UseDB(db)

	// 生成模型代码
	g.GenerateAllTable()

	// 执行生成
	g.Execute()
}

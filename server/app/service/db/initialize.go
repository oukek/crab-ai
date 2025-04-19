package db

import (
	"log"
	"oukek/crab-ai/app/model"
)

// InitializeDatabase 初始化数据库
func InitializeDatabase() {
	// 连接数据库
	Init()

	// 自动迁移表结构
	err := AutoMigrate(
		&model.Provider{},
		// 这里可以添加其他模型
	)

	if err != nil {
		log.Fatalf("自动迁移数据库失败: %v", err)
	}

	log.Println("数据库迁移完成")
} 
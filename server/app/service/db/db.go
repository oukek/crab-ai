package db

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"oukek/crab-ai/app/common"

	"github.com/joho/godotenv"
	"github.com/ncruces/go-sqlite3/gormlite"
	"gorm.io/gorm"
)

var (
	instance *gorm.DB
	once     sync.Once
)

// 初始化数据库连接
func Init() {
	once.Do(func() {
		// 加载环境变量
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file, using default settings")
		}

		// 从环境变量获取数据库文件路径
		dbPath := os.Getenv("DB_FILE_PATH")
		if dbPath == "" {
			dbPath = "./data.db" // 默认路径
		}

		// 连接到SQLite数据库
		db, err := gorm.Open(gormlite.Open(dbPath), &gorm.Config{
			Logger: &common.GormLogger{Log: common.BaseLogger},
		})

		if err != nil {
			panic(fmt.Sprintf("无法连接到数据库: %v", err))
		}

		// 设置连接池参数
		sqlDB, err := db.DB()
		if err != nil {
			panic(fmt.Sprintf("获取数据库连接失败: %v", err))
		}

		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)

		instance = db
		log.Println("数据库连接成功：", dbPath)
	})
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	if instance == nil {
		Init()
	}
	return instance
}

// AutoMigrate 自动迁移数据库结构
func AutoMigrate(models ...interface{}) error {
	return GetDB().AutoMigrate(models...)
}

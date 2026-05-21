package database

import (
	"dinosaur-emoji-service/config"
	"dinosaur-emoji-service/internal/model"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitMySQL 初始化MySQL连接
func InitMySQL(cfg config.MySQLConfig) {
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}
	fmt.Println("数据库连接成功")

	// 设置连接池
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)

	// 自动迁移：根据User结构体创建或更新表结构
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Emoji{})
	db.AutoMigrate(&model.EmojiPack{})
	db.AutoMigrate(&model.EmojiPack_Emoji{})

	DB = db
}

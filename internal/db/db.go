package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func OpenDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("text.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // ログを無効化
	})
	if err != nil {
		log.Fatalln("DB接続失敗")
	}

	err = db.AutoMigrate(&Project{})
	if err != nil {
		log.Panicln(err)
	}

	err = db.AutoMigrate(&WatchRepository{})
	if err != nil {
		log.Panicln(err)
	}

	err = db.AutoMigrate(&Notification{})
	if err != nil {
		log.Panicln(err)
	}

	return db
}

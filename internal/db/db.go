package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func OpenDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("text.db"), &gorm.Config{})
	if err != nil {
		log.Fatalln("DB接続失敗")
	}

	err = db.AutoMigrate(&WatchRepository{})
	if err != nil {
		log.Panicln(err)
	}

	err = db.AutoMigrate(&Project{})
	if err != nil {
		log.Panicln(err)
	}

	err = db.AutoMigrate(&Notification{})
	if err != nil {
		log.Panicln(err)
	}

	return db
}

func GetWatchRepositories(db *gorm.DB) ([]WatchRepository, error) {
	var watchRepos []WatchRepository
	result := db.Find(&watchRepos)
	if result.Error != nil {
		return nil, result.Error
	}
	return watchRepos, nil
}

package db

import (
	"time"

	"gorm.io/gorm"
)

type WatchRepository struct {
	gorm.Model
	Owner       string
	Name        string
	PublishedAt time.Time
}

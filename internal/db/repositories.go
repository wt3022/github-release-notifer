package db

import (
	"time"

	"gorm.io/gorm"
)

type WatchRepository struct {
	gorm.Model
	Owner                string    `json:"owner" label:"Githubユーザー名"`
	Name                 string    `json:"name" label:"リポジトリ名"`
	LastNotificationDate time.Time `json:"last_notification_date" label:"最終通知日"`
}

package db

import (
	"time"

	"gorm.io/gorm"
)

type WatchRepository struct {
	gorm.Model
	Owner                string    `json:"owner" label:"Githubユーザー名" binding:"required"`
	Name                 string    `json:"name" label:"リポジトリ名" binding:"required"`
	LastNotificationDate time.Time `json:"last_notification_date" label:"最終通知日" gorm:"default:CURRENT_TIMESTAMP"`
	ProjectID            uint      `json:"project_id" label:"プロジェクトID" binding:"required"`
}

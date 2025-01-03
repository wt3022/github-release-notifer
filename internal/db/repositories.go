package db

import (
	"time"
)

type WatchType string

const (
	WatchTypeTag     WatchType = "tag"
	WatchTypeRelease WatchType = "release"
)

type WatchRepository struct {
	ID                   uint `gorm:"primarykey"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	Owner                string    `json:"owner" label:"Githubユーザー名" binding:"required"`
	Name                 string    `json:"name" label:"リポジトリ名" binding:"required"`
	WatchType            WatchType `json:"watch_type" label:"監視タイプ" binding:"required"`
	LastNotificationDate time.Time `json:"last_notification_date" label:"最終通知日" gorm:"default:CURRENT_TIMESTAMP"`
	LastPublishedAt      time.Time `json:"last_published_at" label:"最終公開日" gorm:"default:CURRENT_TIMESTAMP"`
	ProjectID            uint      `json:"project_id" label:"プロジェクトID" binding:"required"`
}

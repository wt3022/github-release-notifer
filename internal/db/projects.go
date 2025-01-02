package db

import (
	"time"
)

type Project struct {
	ID                uint `gorm:"primarykey"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Name              string            `json:"name" label:"プロジェクト名" gorm:"size:256;unique" binding:"required"`
	Description       *string           `json:"description" label:"プロジェクト説明"`
	Notification      *Notification     `json:"notification" label:"通知設定" binding:"required"`
	WatchRepositories []WatchRepository `json:"watch_repositories" label:"監視リポジトリ"`
}

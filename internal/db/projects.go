package db

import (
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Name              string            `json:"name" label:"プロジェクト名" gorm:"size:256" binding:"required"`
	Description       *string           `json:"description" label:"プロジェクト説明"`
	Notification      *Notification     `json:"notification" label:"通知設定"`
	WatchRepositories []WatchRepository `json:"watch_repositories" label:"監視リポジトリ"`
}

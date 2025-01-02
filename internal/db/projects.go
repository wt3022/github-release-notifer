package db

import (
	"errors"
	"time"

	"gorm.io/gorm"
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

func (p *Project) Validate(tx *gorm.DB) (err error) {
	var count int64
	tx.Model(&Project{}).Where("name = ?", p.Name).Count(&count)
	if count > 0 {
		return errors.New("このプロジェクト名は既に存在します")
	}
	return
}

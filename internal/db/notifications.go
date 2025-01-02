package db

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type NotificationType string

const (
	Slack NotificationType = "slack"
	Email NotificationType = "email"
)

type Notification struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	ProjectID uint             `json:"project" label:"プロジェクトID" gorm:"unique;onDelete:CASCADE"`
	Type      NotificationType `json:"type" label:"通知方法" gorm:"type:varchar(20);not null" binding:"required"`
}

func (n *Notification) BeforeSave(tx *gorm.DB) (err error) {
	if n.Type != Slack && n.Type != Email {
		return errors.New("通知方法は slack か email を選択してください")
	}

	// プロジェクトは一意である必要がある
	var count int64
	tx.Model(&Notification{}).Where("project_id = ?", n.ProjectID).Count(&count)
	if count > 0 {
		return errors.New("このプロジェクトは既に通知設定が存在します")
	}

	return
}

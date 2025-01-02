package db

import (
	"time"
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

package db

import (
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Name        string  `gorm:"size:256" json:"name" label:"プロジェクト名"`
	Description *string `json:"description" label:"プロジェクト説明"`
}

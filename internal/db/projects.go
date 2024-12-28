package db

import (
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Name        string  `gorm:"size:256" json:"name"`
	Description *string `json:"description"`
}

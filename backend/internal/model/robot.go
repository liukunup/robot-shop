package model

import (
	"time"

	"gorm.io/gorm"
)

type Robot struct {
	Id        uint   `gorm:"primarykey"`
	RobotId   string `gorm:"unique;not null"`
	Name      string
	Desc      string
	Webhook   string `gorm:"unique;not null"`
	Callback  string `gorm:"unique;not null"`
	Options   string `gorm:"type:text"` // JSON string for options
	Enabled   bool   `gorm:"default:true"`
	Owner     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m *Robot) TableName() string {
	return "robots"
}

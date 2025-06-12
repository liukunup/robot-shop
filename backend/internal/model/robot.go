package model

import (
	"gorm.io/gorm"
	"time"
)

type Robot struct {
	Id        uint   `gorm:"primarykey"`
	RobotId   string `gorm:"unique;not null"`
	Name      string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m *Robot) TableName() string {
    return "robots"
}

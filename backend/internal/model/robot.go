package model

import "gorm.io/gorm"

type Robot struct {
	gorm.Model
	RobotId  string `gorm:"unique;not null"`
	Name     string `gorm:"type:varchar(255);not null;comment:'名称'"`
	Desc     string `gorm:"type:varchar(255);not null;comment:'描述'"`
	Webhook  string `gorm:"unique;not null"`
	Callback string `gorm:"unique;not null"`
	Options  string `gorm:"type:text"` // JSON string for options
	Enabled  bool   `gorm:"default:true"`
	Owner    string `gorm:"type:varchar(255);not null;comment:'所有者'"`
}

func (m *Robot) TableName() string {
	return "robot"
}

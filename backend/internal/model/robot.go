package model

import "gorm.io/gorm"

type Robot struct {
	gorm.Model
}

func (m *Robot) TableName() string {
    return "robot"
}

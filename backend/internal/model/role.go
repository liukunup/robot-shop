package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model

	Name string `gorm:"type:varchar(255);uniqueIndex;comment:Custom Name"`
	Role string `gorm:"type:varchar(255);uniqueIndex;comment:Casbin Role"`
}

func (m *Role) TableName() string {
	return "role"
}

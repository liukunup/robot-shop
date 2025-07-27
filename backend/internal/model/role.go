package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model

	Name       string `gorm:"column:name;type:varchar(255);uniqueIndex;comment:Custom Name"`
	CasbinRole string `gorm:"column:casbin_role;type:varchar(255);uniqueIndex;comment:Casbin Role"`
}

func (m *Role) TableName() string {
	return "role"
}

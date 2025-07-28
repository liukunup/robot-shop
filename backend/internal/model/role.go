package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model

	Name       string `gorm:"column:name;type:varchar(64);uniqueIndex:uniq_name;comment:'角色名称(用于显示)'"`
	CasbinRole string `gorm:"column:casbin_role;type:varchar(128);uniqueIndex:uniq_casbin_role;comment:'Casbin角色标识'"`
}

func (m *Role) TableName() string {
	return "role"
}

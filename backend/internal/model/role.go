package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name string `json:"name" gorm:"column:name;type:varchar(100);uniqueIndex;comment:角色名"`
	Sid  string `json:"sid" gorm:"column:sid;type:varchar(100);uniqueIndex;comment:角色标识"`
}

func (m *Role) TableName() string {
    return "role"
}

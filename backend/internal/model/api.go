package model

import "gorm.io/gorm"

type Api struct {
	gorm.Model

	Group  string `gorm:"column:group;type:varchar(255);not null;comment:'分组'"`
	Name   string `gorm:"column:name;type:varchar(255);not null;comment:'名称'"`
	Path   string `gorm:"column:path;type:varchar(255);not null;comment:'路径'"`
	Method string `gorm:"column:method;type:varchar(255);not null;comment:'方法'"`
}

func (m *Api) TableName() string {
	return "api"
}

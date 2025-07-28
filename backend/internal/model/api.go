package model

import "gorm.io/gorm"

type Api struct {
	gorm.Model

	Group  string `gorm:"column:group;type:varchar(100);not null;comment:'API分组'"`
	Name   string `gorm:"column:name;type:varchar(255);not null;comment:'API名称'"`
	Path   string `gorm:"column:path;type:varchar(512);not null;comment:'API路径'"`
	Method string `gorm:"column:method;type:varchar(10);not null;index;check:method IN ('GET','POST','PUT','DELETE','PATCH','HEAD','OPTIONS');comment:'HTTP请求方法'"`
}

func (m *Api) TableName() string {
	return "api"
}

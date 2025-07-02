package model

import "gorm.io/gorm"

type Robot struct {
	gorm.Model

	Name     string `gorm:"column:name;type:varchar(255);not null;comment:'名称'"`
	Desc     string `gorm:"column:desc;type:varchar(255);comment:'描述'"`
	Webhook  string `gorm:"column:webhook;type:varchar(255);comment:'通知地址'"`
	Callback string `gorm:"column:callback;type:varchar(255);comment:'回调地址'"`
	Enabled  bool   `gorm:"column:enabled;type:tinyint(1);not null;comment:'是否启用'"`
	Owner    string `gorm:"column:owner;type:varchar(255);comment:'所有者'"`
}

func (m *Robot) TableName() string {
	return "robot"
}

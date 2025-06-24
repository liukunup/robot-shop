package model

import "gorm.io/gorm"

type Robot struct {
	gorm.Model

	Name     string `gorm:"type:varchar(255);not null;comment:'名称'"`
	Desc     string `gorm:"type:varchar(255);comment:'描述'"`
	Webhook  string `gorm:"type:varchar(255);comment:'回调地址'"`
	Callback string `gorm:"type:varchar(255);comment:'通知地址'"`
	Enabled  bool   `gorm:"type:tinyint(1);default:true;comment:'是否启用'"`
	Owner    string `gorm:"type:varchar(255);comment:'所有者'"`
}

func (m *Robot) TableName() string {
	return "robot"
}

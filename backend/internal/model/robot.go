package model

import "gorm.io/gorm"

type Robot struct {
	gorm.Model

	Name     string `gorm:"column:name;type:varchar(100);not null;index;comment:'机器人名称'"`
	Desc     string `gorm:"column:desc;type:varchar(512);comment:'详细描述'"`
	Webhook  string `gorm:"column:webhook;type:varchar(512);not null;comment:'Webhook 通知地址'"`
	Callback string `gorm:"column:callback;type:varchar(512);comment:'RobotShop 回调地址'"`
	Enabled  bool   `gorm:"column:enabled;type:tinyint(1);not null;default:0;comment:'启用状态:0-禁用,1-启用'"`
	Owner    string `gorm:"column:owner;type:varchar(64);index;comment:'所有者ID/标识'"`
}

func (m *Robot) TableName() string {
	return "robot"
}

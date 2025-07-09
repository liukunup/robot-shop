package model

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Email    string `gorm:"column:email;type:varchar(255);not null;unique;index;comment:'电子邮件'"`
	Username string `gorm:"column:username;type:varchar(255);not null;unique;index;comment:'用户名'"`
	Password string `gorm:"column:password;type:varchar(255);not null;comment:'密码哈希值'"`
	Avatar   string `gorm:"column:avatar;type:varchar(255);comment:'头像'"`
	Nickname string `gorm:"column:nickname;type:varchar(255);comment:'昵称'"`
	Bio      string `gorm:"column:bio;type:text;comment:'个人简介'"`
	Language string `gorm:"column:language;type:varchar(255);comment:'语言'"`
	Timezone string `gorm:"column:timezone;type:varchar(255);comment:'时区'"`
	Theme    string `gorm:"column:theme;type:varchar(255);comment:'主题'"`
	Status   int    `gorm:"column:status;type:int;comment:'状态 0:待激活 1:正常 2:禁用'"`
}

func (m *User) TableName() string {
	return "user"
}

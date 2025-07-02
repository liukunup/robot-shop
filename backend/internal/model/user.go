package model

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Username string `gorm:"column:username;type:varchar(255);not null;unique;index;comment:'用户名'"`
	Password string `gorm:"column:password;type:varchar(255);not null;comment:'密码哈希值'"`
	Nickname string `gorm:"column:nickname;type:varchar(255);comment:'昵称'"`
	Avatar   string `gorm:"column:avatar;type:varchar(255);comment:'头像'"`
	Email    string `gorm:"column:email;type:varchar(255);not null;unique;index;comment:'电子邮件'"`
	Phone    string `gorm:"column:phone;type:varchar(255);comment:'手机号'"`
	Status   int    `gorm:"column:status;type:int;comment:'状态 0:待激活 1:正常 2:禁用'"`
}

func (m *User) TableName() string {
	return "user"
}

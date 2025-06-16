package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(255);not null;unique;index;comment:'用户名'"`
	Password string `gorm:"type:varchar(255);not null;comment:'密码哈希值'"`
	Nickname string `gorm:"type:varchar(255);comment:'昵称'"`
	Email    string `gorm:"type:varchar(255);not null;comment:'电子邮件'"`
	Phone    string `gorm:"type:varchar(255);comment:'手机号'"`
	Avatar   string `gorm:"type:varchar(255);comment:'头像'"`
}

func (m *User) TableName() string {
	return "user"
}

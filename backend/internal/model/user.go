package model

import "gorm.io/gorm"

type User struct {
	gorm.Model

	// 登录凭证
	Email    string `gorm:"column:email;type:varchar(128);not null;uniqueIndex:uniq_email;comment:'电子邮箱(登录用)'"`
	Username string `gorm:"column:username;type:varchar(64);not null;uniqueIndex:uniq_username;comment:'用户名(登录用)'"`
	Password string `gorm:"column:password;type:varchar(255);not null;comment:'BCrypt密码哈希值'"`

	// 个人信息
	Avatar   string `gorm:"column:avatar;type:varchar(512);comment:'头像地址'"`
	Nickname string `gorm:"column:nickname;type:varchar(64);comment:'昵称'"`
	Bio      string `gorm:"column:bio;type:varchar(500);comment:'个人简介'"`

	// 偏好设置
	Language string `gorm:"column:language;type:varchar(32);default:'zh-CN';comment:'语言'"`
	Timezone string `gorm:"column:timezone;type:varchar(64);default:'Asia/Shanghai';comment:'时区'"`
	Theme    string `gorm:"column:theme;type:varchar(32);default:'light';comment:'主题配色: light/dark'"`

	// 账户状态
	Status int `gorm:"column:status;type:tinyint;not null;default:0;comment:'账户状态:0-待激活,1-正常,2-禁用'"`
}

func (m *User) TableName() string {
	return "user"
}

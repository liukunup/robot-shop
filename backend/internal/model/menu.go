package model

import "gorm.io/gorm"

type Menu struct {
	gorm.Model

	ParentID  uint   `gorm:"column:parent_id;index;comment:父级菜单ID"`
	Path      string `gorm:"column:path;type:varchar(255);comment:地址"`
	Component string `gorm:"column:component;type:varchar(255);comment:组件"`
	Name      string `gorm:"column:name;type:varchar(100);comment:名称"`
	Icon      string `gorm:"column:icon;type:varchar(100);comment:图标"`
	Weight    int    `gorm:"column:weight;default:0;comment:权重"`
}

func (m *Menu) TableName() string {
	return "menu"
}

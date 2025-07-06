package model

import "gorm.io/gorm"

type Menu struct {
	gorm.Model

	ParentID           uint   `gorm:"column:parent_id;index;comment:父级菜单"`
	Icon               string `gorm:"column:icon;type:varchar(255);comment:图标"`
	Name               string `gorm:"column:name;type:varchar(255);comment:名称"`
	Path               string `gorm:"column:path;type:varchar(255);comment:路由"`
	Component          string `gorm:"column:component;type:varchar(255);comment:组件"`
	Access             string `gorm:"column:access;type:varchar(255);comment:可见性"`
	Locale             string `gorm:"column:locale;type:varchar(255);comment:国际化"`
	Redirect           string `gorm:"column:redirect;type:varchar(255);comment:重定向"`
	Target             string `gorm:"column:target;type:varchar(255);comment:指定外链打开形式"`
	HideChildrenInMenu bool   `gorm:"column:hide_children_in_menu;default:false;comment:隐藏子节点"`
	HideInMenu         bool   `gorm:"column:hide_in_menu;default:false;comment:隐藏自身和子节点"`
	FlatMenu           bool   `gorm:"column:flat_menu;default:false;comment:隐藏自身+子节点提升并打平"`
	Disabled           bool   `gorm:"column:disabled;default:false;"`
	Tooltip            string `gorm:"column:tooltip;type:varchar(255);"`
	DisabledTooltip    bool   `gorm:"column:disabled_tooltip;default:false;"`
	Key                string `gorm:"column:key;type:varchar(255);"`
	ParentKeys         string `gorm:"column:parent_keys;type:varchar(255);"`
}

func (m *Menu) TableName() string {
	return "menu"
}

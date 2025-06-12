package model

import "gorm.io/gorm"

type Skill struct {
	gorm.Model
}

func (m *Skill) TableName() string {
    return "skill"
}

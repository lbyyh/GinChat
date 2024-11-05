package models

import "gorm.io/gorm"

type GroupBasic struct {
	gorm.Model
	Name    string
	OwnerId uint //谁的关系信息
	Icon    string
	Type    int // 0 1 2
	Desc    string
}

func (table *GroupBasic) TableName() string {
	return "group_basic"
}

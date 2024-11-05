package models

import (
	"GinChat/utils"
	"fmt"
	"gorm.io/gorm"
)

type Community struct {
	gorm.Model
	Name    string
	OwnerId uint
	Img     string
	Desc    string
}

func CreateCommunity(community *Community) (int, string) {
	if len(community.Name) == 0 {
		return -1, "群昵称不能为空"
	}
	if community.OwnerId == 0 {
		return -1, "请先登录"
	}
	if err := utils.DB.Create(&community).Error; err != nil {
		fmt.Println(err)
		return -1, "建群失败"
	}
	return 0, "建群成功"
}

func SearchGroup(ownerId uint) ([]Community, string) {
	contacts := make([]Contact, 0)
	objIds := make([]uint64, 0)
	utils.DB.Where("owner_id = ? and type=2", ownerId).Find(&contacts)
	for _, v := range contacts {
		fmt.Println(">>>>>>>>>>>>>>", v)
		objIds = append(objIds, uint64(v.TargetId))
	}
	community := make([]Community, 0)
	utils.DB.Where("id in ?", objIds).Find(&community)
	return community, ""
}

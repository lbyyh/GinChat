package models

import (
	"GinChat/utils"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type UserBasic struct {
	gorm.Model
	Name          string
	PassWord      string
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email         string `valid:"email"`
	Identity      string
	ClientIp      string
	ClientPort    string
	Salt          string
	loginTime     time.Time
	HeartbeatTime time.Time
	LoginOutTime  time.Time `gorm:"column:login_out_time" json:"login_out_time"`
	IsLogout      bool
	DeviceInfo    string
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}

func FindUserByNameAndPwd(name, password string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name=? and pass_word=?", name, password).First(&user)

	//token加密
	str := fmt.Sprintf("%d", time.Now().Unix())
	temp := utils.MD5Encode(str)
	utils.DB.Model(&user).Where("id = ?", user.ID).Update("identity", temp)
	return user
}

func FindUserByName(name string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name=?", name).First(&user)
	return user
}

func FindUserByID(id uint) UserBasic {
	user := UserBasic{}
	utils.DB.Where("id=?", id).First(&user)
	return user
}

// FindCommunityByID 根据ID找群聊
func FindCommunityByID(id uint) Community {
	community := Community{}
	utils.DB.Where("id=?", id).First(&community)
	return community
}

// FindGroupIdByName 根据群聊名找群聊ID
func FindGroupIdByName(name string) (uint, error) {
	community := Community{}
	err := utils.DB.Where("name=?", name).First(&community).Error
	if err != nil {
		return community.ID, err
	}
	return community.ID, nil
}

func FindUserByPhone(phone string) *gorm.DB {
	user := UserBasic{}
	return utils.DB.Where("phone=?", phone).First(&user)
}

func FindUserByEmail(email string) *gorm.DB {
	user := UserBasic{}
	return utils.DB.Where("email=?", email).First(&user)
}

func CreateUser(user UserBasic) *gorm.DB {
	return utils.DB.Create(&user)
}

func DeleteUser(user UserBasic) *gorm.DB {
	return utils.DB.Delete(&user)
}

func UpdaterUser(user UserBasic) *gorm.DB {
	return utils.DB.Model(&user).Updates(UserBasic{Name: user.Name, PassWord: user.PassWord, Phone: user.Phone, Email: user.Email})
}

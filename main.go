package main

import (
	"GinChat/models"
	"GinChat/router"
	"GinChat/utils"
)

func main() {
	utils.InitConfig()
	utils.InitMySql()
	utils.InitRedis()
	utils.InitMongo()          // 初始化MongoDB连接
	go models.HandleMessages() //群聊消息提交

	r := router.Router()
	r.Run(":8081")
}

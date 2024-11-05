package router

import (
	"GinChat/docs"
	"GinChat/models"
	"GinChat/service"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	//生成API文档
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler)) //注册了一个路由，使得能够通过/swagger/*any这个路径来访问Swagger生成的API文档。
	//静态资源
	r.Static("/asset", "asset/")
	r.LoadHTMLGlob("views/**/*")
	//首页
	r.GET("/", service.GetIndex)
	r.GET("/index", service.GetIndex)
	r.GET("/toRegister", service.ToRegister)
	r.GET("/toChat", service.ToChat)
	r.GET("/chat", service.Chat)
	r.GET("/chatUI", service.ChatUI)
	r.GET("/chatGroupUI", service.ChatGroupUI)

	r.GET("/user/getUserList", service.GetUserList)
	r.POST("/user/createUser", service.CreateUser)
	r.POST("/user/deleteUser", service.DeleteUser)
	r.POST("/user/updateUser", service.UpdateUser)
	r.POST("/user/login", service.Login)

	//发送消息
	r.GET("/user/sendMsg", service.SendMsg)
	r.GET("/user/sendUserMsg", service.SendUserMsg)
	r.POST("/user/searchFriends", service.SearchFriends)
	r.POST("/user/addFriend", service.AddFriend) //加好友
	r.POST("/user/addGroup", service.AddGroup)   //加群
	r.POST("/attach/upload", service.Upload)
	r.POST("/user/audioUploadHandler", models.AudioUploadHandler)
	r.POST("/user/createCommunity", service.CreateCommunity)       //创建群聊
	r.POST("/contact/loadCommunity", service.LoadCommunity)        //群聊列表
	r.POST("/contact/saveMessage", service.SaveMessage)            //保存前端发送来的消息数据到mongodb
	r.GET("/contact/getRecentMessages", service.GetRecentMessages) // 查询最近的消息记录

	// 设置群聊连接的路由
	r.GET("/groupChat", models.GroupChatConnection)

	// 设置消息广播的路由
	r.POST("/sendGroupMessage", models.SendGroupMessage)

	return r
}

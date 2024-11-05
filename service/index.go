package service

import (
	"GinChat/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"strconv"
)

// GetIndex godoc
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	// 使用gin框架的模板渲染功能
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "聊天室登录",
	})
}

// ToRegister godoc
// @Tags 注册页面
// @Success 200 {string} html
// @Router /index [get]
func ToRegister(c *gin.Context) {
	// 使用gin框架的模板渲染功能
	c.HTML(http.StatusOK, "register.html", gin.H{
		"title": "聊天室登录",
	})
}

// ToChat godoc
// @Tags 主页面
// @Success 200 {string} html
// @Router /index [get]
func ToChat(c *gin.Context) {
	ind, err := template.ParseFiles("D:\\code\\GO\\GinChat\\views\\chat\\index.html",
		"D:\\code\\GO\\GinChat\\views\\chat\\head.html",
		"D:\\code\\GO\\GinChat\\views\\chat\\foot.html",
		"D:\\code\\GO\\GinChat\\views\\chat\\tabmenu.html",
		"D:\\code\\GO\\GinChat\\views\\chat\\concat.html",
		"D:\\code\\GO\\GinChat\\views\\chat\\group.html",
		"D:\\code\\GO\\GinChat\\views\\chat\\profile.html",
		"D:\\code\\GO\\GinChat\\views\\chat\\main.html",
	)
	if err != nil {
		panic(err)
	}
	userId, _ := strconv.Atoi(c.Query("userId"))
	fmt.Println("userId:", userId)
	token := c.Query("token")
	user := models.UserBasic{}
	user.ID = uint(userId)
	user.Identity = token
	fmt.Println("user:--------", user, token)
	ind.Execute(c.Writer, user)
}

// ChatUI godoc
// @Tags 聊天页面
// @Success 200 {string} html
// @Router /index [get]
func ChatUI(c *gin.Context) {
	//userId, _ := strconv.Atoi(c.Query("userId"))
	//fmt.Println("userId:", userId)
	// 使用gin框架的模板渲染功能
	c.HTML(http.StatusOK, "main.html", gin.H{
		"title": "聊天登录",
		//"data":  userId,
	})
}

// ChatGroupUI godoc
// @Tags 群聊天页面
// @Success 200 {string} html
// @Router /index [get]
func ChatGroupUI(c *gin.Context) {
	//userId, _ := strconv.Atoi(c.Query("userId"))
	//fmt.Println("userId:", userId)
	// 使用gin框架的模板渲染功能
	c.HTML(http.StatusOK, "groupChat.html", gin.H{
		"title": "群聊天登录",
		//"data":  userId,
	})
}

func Chat(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}

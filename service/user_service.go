package service

import (
	"GinChat/models"
	"GinChat/tools"
	"GinChat/utils"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// GetUserList godoc
// @Summary 获取用户列表
// @Description 获取所有用户的列表
// @Tags 用户
// @Produce json
// @Success 200 {object} tools.ECode{data}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	userList := models.GetUserList()
	c.JSON(200, tools.ECode{
		Code:    0,
		Message: "获取用户列表成功",
		Data:    userList,
	})
}

type CUser struct {
	Name       string `json:"name"`
	Password   string `json:"password"`
	RePassword string `json:"rePassword"`
}

// CreateUser godoc
// @Summary 创建用户
// @Description 创建一个新的用户
// @Tags 用户
// @Accept json
// @Produce json
// @Param user body models.UserBasic true "用户基本信息"
// @Success 200 {object} tools.ECode{code=int,message=string,data=models.UserBasic}
// @Failure 200 {object} tools.ECode{code=int,message=string}
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var user CUser
	err := c.ShouldBind(&user)
	fmt.Println(user)
	if err != nil {
		c.JSON(200, tools.ECode{
			Code:    1,
			Message: err.Error(),
		})
		return
	}
	data := models.FindUserByName(user.Name)
	if user.Name == "" || user.Password == "" || user.RePassword == "" {
		c.JSON(200, tools.ECode{
			Code:    1,
			Message: "用户名或密码不为空",
			Data:    nil,
		})
		return
	}
	if data.Name != "" {
		c.JSON(200, tools.ECode{
			Code:    1,
			Message: "用户名已存在",
			Data:    nil,
		})
		return
	}
	if user.Password != user.RePassword {
		c.JSON(200, tools.ECode{
			Code:    1,
			Message: "两次密码不一致",
			Data:    nil,
		})
		return
	}
	slate := fmt.Sprintf("%06b", rand.Int31())
	fmt.Printf("Slate: %s\n", slate)
	NewUser := models.UserBasic{
		Name:     user.Name,
		Salt:     slate,
		PassWord: utils.MakePassword(user.Password, slate),
	}
	models.CreateUser(NewUser)
	c.JSON(200, tools.ECode{
		Code:    0,
		Message: "创建用户成功",
	})
}

// DeleteUser godoc
// @Summary 删除用户
// @Description 通过用户 ID 删除用户
// @Tags 用户
// @Produce json
// @Param id query int true "用户 ID"
// @Success 200 {object} tools.ECode{code=int,message=string}
// @Router /users [delete]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(200, tools.ECode{
		Code:    0,
		Message: "删除用户成功",
	})
}

// UpdateUser godoc
// @Summary 更新用户
// @Description 更新用户的基本信息
// @Tags 用户
// @Accept json
// @Produce json
// @Param user body models.UserBasic true "用户基本信息"
// @Success 200 {object} tools.ECode{code=int,message=string,data=models.UserBasic}
// @Router /users [put]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	c.ShouldBind(&user)
	_, err := govalidator.ValidateStruct(&user)
	if err != nil {
		c.JSON(200, tools.ECode{
			Code:    1,
			Message: err.Error(),
		})
		return
	}
	models.UpdaterUser(user)
	c.JSON(200, tools.ECode{
		Code:    0,
		Message: "更新用户成功",
	})
}

type loginUser struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	user := models.UserBasic{}
	loginUser := models.UserBasic{}
	err := c.ShouldBind(&loginUser)
	if err != nil {
		fmt.Println("111111111err----", err)
		return
	}

	user = models.FindUserByName(loginUser.Name)
	fmt.Println("data:-----------11111----------", user)
	if user.Name == "" {
		c.JSON(200, tools.ECode{
			Code:    1,
			Message: "用户名不存在",
		})
		return
	}
	if user.PassWord != utils.MakePassword(loginUser.PassWord, user.Salt) {
		c.JSON(200, tools.ECode{
			Code:    1,
			Message: "密码错误",
		})
		return
	}
	token, err := tools.GetJwt(int64(user.ID), user.Name)
	user.Identity = token
	fmt.Println("data----------", user)
	if err != nil {
		c.JSON(200, tools.ECode{
			Code:    1,
			Message: "token生成失败",
			Data:    nil,
		})
		return
	}
	c.JSON(200, tools.ECode{
		Code:    0,
		Message: "登录成功",
		Data:    user,
	})
}

// 防止跨域站点伪造请求
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		//origin := r.Header.Get("Origin")
		//// 检查是否为允许的来源
		//if origin == "https://your-allowed-origin.com" {
		//	return true
		//}
		//return false
		return true
	},
}

func SendMsg(c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	MsgHandler(ws, c)
}

func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	msg, err := utils.Subscribe(c, utils.PublishKey)
	if err != nil {
		fmt.Println(err)
		return
	}
	tm := time.Now().Format("2006-01-02 15:04:05")
	m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
	err = ws.WriteMessage(1, []byte(m))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}

func SearchFriends(c *gin.Context) {
	var requestData struct {
		UserID uint `json:"userId"` // 定义一个结构体匹配前端发送的JSON格式
	}

	// 从请求体中解析JSON数据到requestData
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 使用requestData中的UserID来调用模型层的搜索函数
	users := models.SearchFriend(requestData.UserID)

	fmt.Println("users=======", users)

	// 使用自定义的RespOKList函数返回查询到的好友列表
	utils.RespOKList(c.Writer, users, len(users))
}

func AddFriend(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	targetId, _ := strconv.Atoi(c.Request.FormValue("targetId"))
	code, msg := models.AddFriend(uint(userId), uint(targetId))
	if code == 1 {
		utils.RespOK(c.Writer, code, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

// CreateCommunity 创建群聊
func CreateCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.PostForm("ownerId"))
	name := c.PostForm("name")
	desc := c.PostForm("description")

	// 获取上传的文件
	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		utils.RespFail(c.Writer, "获取图片失败："+err.Error())
		return
	}

	// 建立保存文件的目标目录
	targetDir := filepath.Join("asset", "groupAvatar")
	err = os.MkdirAll(targetDir, os.ModePerm)
	if err != nil {
		utils.RespFail(c.Writer, "创建目录失败："+err.Error())
		return
	}

	// 生成唯一文件名
	fileNewName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileHeader.Filename) // 使用时间戳和原始文件名生成新文件名
	filePath := filepath.Join(targetDir, fileNewName)

	// 打开上传的文件
	file, err := fileHeader.Open()
	if err != nil {
		utils.RespFail(c.Writer, "打开上传文件失败："+err.Error())
		return
	}
	defer file.Close()

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		utils.RespFail(c.Writer, "创建文件失败："+err.Error())
		return
	}
	defer dst.Close()

	// 将上传的文件复制到目标文件
	if _, err := io.Copy(dst, file); err != nil {
		utils.RespFail(c.Writer, "保存图片失败："+err.Error())
		return
	}

	community := models.Community{
		OwnerId: uint(ownerId),
		Name:    name,
		Desc:    desc,
		Img:     "/asset/groupAvatar/" + fileNewName, // 存储相对路径方便后续使用
	}

	// 调用模型层的函数，创建群聊并保存至数据库
	code, msg := models.CreateCommunity(&community)
	if code == 0 {
		utils.RespOK(c.Writer, nil, "群聊创建成功")
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

func LoadCommunity(c *gin.Context) {
	var requestData struct {
		UserID uint `json:"userId"` // 定义一个结构体匹配前端发送的JSON格式
	}

	// 从请求体中解析JSON数据到requestData
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 使用requestData中的UserID来调用模型层的搜索函数
	users, _ := models.SearchGroup(requestData.UserID)

	fmt.Println("users=======", users)

	// 使用自定义的RespOKList函数返回查询到的好友列表
	utils.RespOKList(c.Writer, users, len(users))
}

func AddGroup(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	groupInfo := c.Request.FormValue("groupInfo")
	var groupId uint
	// 尝试将groupInfo解析为int型，以判断是否为群聊ID
	if id, err := strconv.Atoi(groupInfo); err == nil {
		groupId = uint(id) // 如果是群聊ID，直接使用
	} else {
		// 如果不是群聊ID，则假设是群聊名称并查询获取群聊ID
		groupId, err = models.FindGroupIdByName(groupInfo)
		if err != nil {
			// 如果找不到群聊或查询出错，返回失败信息
			utils.RespFail(c.Writer, "找不到指定的群聊")
			return
		}
	}

	// 调用函数加入群聊
	code, msg := models.AddGroup(uint(userId), groupId)
	if code == 1 {
		utils.RespOK(c.Writer, code, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

// GetRecentMessages 查询最近的消息记录
func GetRecentMessages(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Query("userId"))
	targetId, _ := strconv.Atoi(c.Query("targetId"))

	// 假设你的models层有一个函数可以获取最近的记录
	messages, err := models.GetRecentMessages(uint(userId), uint(targetId))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// SaveMessage 处理保存消息的请求
func SaveMessage(c *gin.Context) {
	var message models.MongoMessage
	if err := c.ShouldBindJSON(&message); err != nil {
		utils.RespFail(c.Writer, "请求参数错误")
		return
	}
	message.CreatedAt = time.Now() // 设置消息创建时间
	err := models.SaveMessage(&message)
	if err != nil {
		fmt.Println("保存消息失败:", err)
		utils.RespFail(c.Writer, "保存消息失败")
	} else {
		utils.RespOK(c.Writer, nil, "消息保存成功")
	}
}

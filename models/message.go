package models

import (
	"GinChat/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/fatih/set"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type Message struct {
	gorm.Model
	FormId   int64  //发送者
	TargetId int64  //接收者
	Type     int    // 发送消息类型 1私聊 2群聊	3广播
	Media    int    //消息类型	1文字 2表情包 3图片 4音频
	Content  string // 消息内容
	Pic      string
	Url      string
	Desc     string
	Amount   int //其他数字统计
}

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

// 映射关系
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

func Chat(writer http.ResponseWriter, request *http.Request) {
	//token := query.Get("token)
	query := request.URL.Query()

	Id := query.Get("userId")
	userId, _ := strconv.ParseInt(Id, 10, 64)
	isvalida := true //checkToken()
	conn, err := (&websocket.Upgrader{
		//token校验
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}
	//userid 和node绑定 并加锁
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()
	//发送逻辑
	go sendProc(node)
	//接收逻辑
	go recvProc(node)
}

func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			fmt.Println("[ws] sendProc >>> msg:", string(data))
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		//dispatch(data)
		broadMsg(data)
		fmt.Println("[ws] recvProc <<<<", string(data))
	}
}

var udpsendChan chan []byte = make(chan []byte, 1024)

func broadMsg(data []byte) {
	udpsendChan <- data
}

func init() {
	go udpSendProc()
	go udpRecvProc()
	fmt.Println("init goroutine")
}

func udpSendProc() {
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 56, 1),
		Port: viper.GetInt("port.udp"),
	})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
	}
	for {
		select {
		case data := <-udpsendChan:
			_, err := con.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

// udp数据接收协程
func udpRecvProc() {
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: viper.GetInt("port.udp"),
	})
	if err != nil {
		fmt.Println(err)
	}
	defer con.Close()
	for {
		var buf [512]byte
		n, err := con.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		dispatch(buf[0:n])
	}
}

func dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {

	case 1:

		fmt.Println("dispatch <<<<", string(data))
		sendMsg(msg.TargetId, data)

	}

}

func sendMsg(userId int64, msg []byte) {
	rwLocker.RLock()
	Node, ok := clientMap[userId]
	rwLocker.RUnlock()
	if ok {
		Node.DataQueue <- msg
	}
}

// Emoji 表情包结构
type Emoji struct {
	Name string `json:"name"`
	Char string `json:"char"`
}

// AudioUploadHandler 语音上传处理
func AudioUploadHandler(c *gin.Context) {
	file, err := c.FormFile("audio")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 生成文件名：当前时间戳 + 原始文件名
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	savePath := filepath.Join("asset/uploadedAudios", filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 假设服务器的静态文件被配置到 /static 下，需要返回静态文件对应的可访问URL
	url := fmt.Sprintf("http://127.0.0.1:8081/%s", savePath)

	// 示例中此处仅返回URL，实际项目中可能包括持续时间等其他信息
	c.JSON(http.StatusOK, gin.H{"url": url})
}

// GetRecentMessages - 获取与特定用户或群组的最近的消息列表
func GetRecentMessages(userId uint, targetId uint) ([]MongoMessage, error) {
	var messages []MongoMessage

	// 获取数据库和集合的引用
	messagesCollection := utils.MongoClient.Database("chat").Collection("messages")

	// 构建查询
	// 请根据你的数据结构和需求调整查询条件和逻辑
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"createdAt", -1}}) // 按创建时间降序排序
	findOptions.SetLimit(10)                       // 限制结果数量为最近的10条记录

	filter := bson.D{
		{"$or", bson.A{
			bson.D{{"FormId", userId}, {"TargetId", targetId}},
			bson.D{{"FormId", targetId}, {"TargetId", userId}},
		}},
	}

	// 执行查询
	cur, err := messagesCollection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Println("Error finding recent messages:", err)
		return nil, err
	}
	defer cur.Close(context.Background())

	// 遍历结果
	for cur.Next(context.Background()) {
		var message MongoMessage
		err := cur.Decode(&message)
		if err != nil {
			log.Println("Error decoding message:", err)
			return nil, err
		}
		messages = append([]MongoMessage{message}, messages...)
	}

	if err := cur.Err(); err != nil {
		log.Println("Cursor error:", err)
		return nil, err
	}

	return messages, nil
}

package models

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

var broadcast = make(chan Message)                          // 广播频道
var groupClients = make(map[int64]map[*websocket.Conn]bool) // 群组ID到客户端集合的映射
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 允许来自所有域的连接，生产环境中应更严格
		return true
	},
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	groupIdStr := r.URL.Query().Get("groupId")
	if groupIdStr == "" {
		http.Error(w, "GroupId is required", http.StatusBadRequest)
		return
	}
	groupId, err := strconv.ParseInt(groupIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid GroupId", http.StatusBadRequest)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %s\n", err.Error())
		return
	}
	//defer ws.Close()

	clientGroupMap, exists := groupClients[groupId]
	if !exists {
		clientGroupMap = make(map[*websocket.Conn]bool)
		groupClients[groupId] = clientGroupMap
	}
	clientGroupMap[ws] = true

	go handleClientMessages(ws, groupId, clientGroupMap)
}

func handleClientMessages(ws *websocket.Conn, groupId int64, clientGroupMap map[*websocket.Conn]bool) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Recovered in handleClientMessages: %v", err)
		}
		delete(clientGroupMap, ws)
		ws.Close()
	}()

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading JSON: %s\n", err.Error())
			break
		}

		// 广播消息到同一群组的其他客户端
		broadcastToGroup(groupId, msg)
	}
}

func broadcastToGroup(groupId int64, msg Message) {
	clientsInGroup, exists := groupClients[groupId]
	if !exists {
		log.Printf("No clients in group: %d\n", groupId)
		return
	}

	for client := range clientsInGroup {
		err := client.WriteJSON(msg)
		if err != nil {
			log.Printf("Error writing JSON: %s\n", err.Error())
			client.Close()
			delete(clientsInGroup, client)
		}
	}
}

func HandleMessages() {
	for {
		msg := <-broadcast
		// 处理接收到的消息，例如广播给群组中的其他成员
		broadcastToGroup(msg.TargetId, msg)
	}
}

// GroupChatConnection 处理WebSocket连接
func GroupChatConnection(c *gin.Context) {
	// ... 这里使用上面提供的HandleConnections函数处理WebSocket连接
	HandleConnections(c.Writer, c.Request)
}

// SendGroupMessage 接收前端发送的消息并将其发送到broadcast通道
func SendGroupMessage(c *gin.Context) {
	var msg Message
	if err := c.BindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message format"})
		return
	}
	broadcast <- msg // 将消息发送到broadcast通道
	c.JSON(http.StatusOK, gin.H{"status": "Message sent"})
}

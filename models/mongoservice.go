package models

import (
	"GinChat/utils"
	"context"
	"gorm.io/gorm"
	"time"
)

// MongoMessage 在 MongoDB 中代表消息记录的结构体
type MongoMessage struct {
	gorm.Model
	FormId    uint      `bson:"FormId"`
	TargetId  uint      `bson:"TargetId"`
	Type      int       `bson:"type"`      // 消息类型
	Media     int       `bson:"media"`     // 媒体类型
	Content   string    `bson:"content"`   // 消息内容
	Duration  string    `bson:"duration"`  // 语音消息的持续时间
	CreatedAt time.Time `bson:"createdAt"` // 消息时间戳
}

// SaveMessage 保存消息到 MongoDB
func SaveMessage(message *MongoMessage) error {
	collection := utils.MongoClient.Database("chat").Collection("messages")
	_, err := collection.InsertOne(context.TODO(), message)
	return err
}

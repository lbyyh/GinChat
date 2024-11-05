package utils

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

var (
	DB          *gorm.DB
	Rdb         *redis.Client
	MongoClient *mongo.Client
)

func InitMySql() {
	//自定义日志输出模板，打印sql语句
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags|log.Lshortfile),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	dbUser := viper.GetString("mysql.username")
	dbPass := viper.GetString("mysql.password")
	dbHost := viper.GetString("mysql.host")
	dbPort := viper.GetString("mysql.port")
	dbName := viper.GetString("mysql.database")
	my := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(my), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic(err)
	}
	DB = db
}

func InitRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     "",
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})
	Rdb = rdb
	pong, err := rdb.Ping(context.TODO()).Result()
	if err != nil {
		fmt.Println("init redis ................", err)
	} else {
		fmt.Println("redis inited ................", pong)
	}
}

const (
	PublishKey = "websocket"
)

// Publish 发起Publish 发消息到redis
func Publish(c context.Context, channel string, msg string) error {
	var err error
	err = Rdb.Publish(c, channel, msg).Err()
	return err
}

// Subscribe Subscribe订阅redis消息
func Subscribe(c context.Context, channel string) (string, error) {
	sub := Rdb.Subscribe(c, channel)
	fmt.Println("Subscribe。。。", sub)
	msg, err := sub.ReceiveMessage(c)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("msg---------", msg)
	return msg.Payload, err
}

// InitMongo mongodb初始化链接
func InitMongo() {
	mongoURI := viper.GetString("mongodb.uri")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to create new MongoDB client: %v", err)
	}

	// 设置连接超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// 可选的，检测连接
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB")
	MongoClient = client
}

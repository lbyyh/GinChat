package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func Upload(c *gin.Context) {
	srcFile, head, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer srcFile.Close() // 确保在返回前关闭源文件

	fileExt := filepath.Ext(head.Filename) // 直接获取文件扩展名
	if fileExt == "" {
		fileExt = ".jpg" // 默认扩展名
	}

	fileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31n(10000), fileExt)
	savePath := filepath.Join("./asset/upload", fileName)

	// 确保文件夹存在
	if _, err := os.Stat("./asset/upload"); os.IsNotExist(err) {
		os.MkdirAll("./asset/upload", os.ModePerm)
	}

	// 创建目标文件
	dstFile, err := os.Create(savePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer dstFile.Close() // 确保在返回前关闭目标文件

	// 复制文件内容
	if _, err = io.Copy(dstFile, srcFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 这里需要返回给客户端能访问的完整URL，假设是http://yourserver.com/asset/upload/
	fullURL := "http://127.0.0.1:8081/" + savePath
	c.JSON(http.StatusOK, gin.H{"url": fullURL, "message": "发送图片成功"})
}

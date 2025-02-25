package main

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建一个默认的 Gin 路由器
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 定义一个处理所有请求的路由
	r.Any("/echo", func(c *gin.Context) {
		// 读取请求的 Body
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to read request body")
			return
		}

		log.Printf("Request received:\nMethod:%s\nURL:%s\nHeaders:%s\nBody:%s", c.Request.Method, c.Request.URL.String(), c.Request.Header, string(body))

		// 设置响应头与请求头一致
		for key, values := range c.Request.Header {
			for _, value := range values {
				c.Header(key, value)
			}
		}

		// 返回请求的 Body
		c.Data(http.StatusOK, c.ContentType(), body)
	})

	// 启动 HTTP 服务器
	err := r.Run(":48080")
	if err != nil {
		return
	}
}

package main

import (
	"github.com/gin-gonic/gin"
	"terraform-service/handler"
)

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 测试 schema Attribute 的 computed 属性
	computed := r.Group("/computed")
	{
		computed.GET("/detail", handler.ComputedDetail)
	}

	err := r.Run(":29999")
	if err != nil {
		panic(err)
	}
}

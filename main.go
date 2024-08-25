package main

import (
	"github.com/gin-gonic/gin"
	"socket/dao"
)

func main() {
	dao.InitDB()
	r := gin.Default()
	r.POST("/register", Register)
	r.GET("/ws/:username", Wsfc)
	r.Run(":1226")
}

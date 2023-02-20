package server

import (
	"chatbot/internal/conf"
	"chatbot/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
)

func ginWeb(c *conf.Server, service *service.GreeterService) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	chatPath := fmt.Sprintf("/chat/%s", c.Chat.Path)
	router.LoadHTMLGlob("html/templates/*")
	router.Static("/html/js", "html/js")
	router.Static("/html/css", "html/css")
	router.GET("/", service.Index)
	router.GET(chatPath, service.Chat)
	router.GET("/ws", service.Ws)
	router.GET("/favicon.ico", service.Test2)
	return router
}

package service

import (
	"github.com/gin-gonic/gin"
	ohttp "net/http"
)

func (s *GreeterService) Chat(ctx *gin.Context) {
	ctx.HTML(ohttp.StatusOK, "chat.html", gin.H{})
}
func (s *GreeterService) Index(ctx *gin.Context) {
	ctx.HTML(ohttp.StatusOK, "index.html", gin.H{})
}
func (s *GreeterService) Test2(ctx *gin.Context) {
	ctx.JSON(200, nil)
}

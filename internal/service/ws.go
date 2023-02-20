package service

import (
	"chatbot/internal/service/ws"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (s *GreeterService) Ws(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		s.log.Error(err)
		return
	}
	ws.NewClient(s.hub, conn, s.log)
}

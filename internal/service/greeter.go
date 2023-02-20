package service

import (
	v1 "chatbot/api/helloworld/v1"
	"chatbot/internal/biz"
	"chatbot/internal/conf"
	"chatbot/internal/service/ws"
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServer
	uc  *biz.GreeterUsecase
	hub *ws.Hub
	log *log.Helper
}

// NewGreeterService new a greeter service.
func NewGreeterService(c *conf.Server, uc *biz.GreeterUsecase, logger log.Logger) *GreeterService {
	gs := &GreeterService{
		uc:  uc,
		hub: ws.NewHub(uc, logger, c.Chat.Mode),
		log: log.NewHelper(logger),
	}
	go gs.hub.Run()
	return gs
}

// SayHello implements helloworld.GreeterServer.
func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	g, err := s.uc.CreateGreeter(ctx, &biz.Greeter{Hello: in.Name})
	if err != nil {
		return nil, err
	}
	return &v1.HelloReply{Message: "Hello " + g.Hello}, nil
}
func (s *GreeterService) Question(ctx context.Context, in *v1.QuestionRequest) (*v1.QuestionReply, error) {

	g, err := s.uc.Question(ctx, in.Q, in.Code)
	if err != nil {
		return nil, err
	}
	b, _ := json.Marshal(g)
	return &v1.QuestionReply{Res: string(b)}, nil
}

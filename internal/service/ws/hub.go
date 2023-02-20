package ws

import (
	"chatbot/internal/biz"
	"chatbot/pkg/define"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"strconv"
	"sync"
	"time"
)

// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan define.HubMsg

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	log  *log.Helper
	uc   *biz.GreeterUsecase
	mode string
	no   int64
	lock *sync.Mutex
}

const (
	modeGroup = "group"
	modeOne   = "one"
)

func NewHub(uc *biz.GreeterUsecase, logger log.Logger, mode string) *Hub {
	if mode != modeGroup {
		mode = modeOne
	}
	return &Hub{
		broadcast:  make(chan define.HubMsg),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		log:        log.NewHelper(logger),
		uc:         uc,
		mode:       mode,
		lock:       &sync.Mutex{},
	}
}
func (h *Hub) genNo() string {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.no++
	a, _ := strconv.Atoi(time.Now().Format("0405"))
	return fmt.Sprintf("%s%s", define.DecimalTo26(h.no), define.DecimalTo26(int64(a)))
}
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				client.receive(message)
			}
		}
	}
}
func (h *Hub) toOther(client *Client, message []byte) {
	// to send other like bot
}
func (h *Hub) askBot(client *Client, jos define.Send) {
	ch := make(chan define.Message, 10)
	// ask bot flag
	token := "@bot "
	tokenL := len(token)
	isCall := len(jos.Msg) > tokenL && jos.Msg[0:tokenL] == token
	if !isCall {
		return
	}
	jos.Msg = jos.Msg[tokenL:]
	// start task
	go func() {
		h.log.Info("start request askBot")
		h.uc.Stream()(jos, ch)
	}()
	for {
		msg, ok := <-ch
		if !ok {
			return
		}
		if client.transfer(msg) {
			return
		}
	}
}

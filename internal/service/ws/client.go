// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ws

import (
	"bytes"
	"chatbot/pkg/define"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gorilla/websocket"
	"strings"
	"sync"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
	// logger
	log *log.Helper
	// send state
	state bool
	// lock for send
	lock *sync.RWMutex
	id   string
}

func NewClient(hub *Hub, conn *websocket.Conn, logger *log.Helper) {
	client := &Client{
		hub:   hub,
		conn:  conn,
		send:  make(chan []byte, 256),
		log:   logger,
		state: true,
		lock:  &sync.RWMutex{},
		id:    hub.genNo(),
	}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

func (c *Client) Close() {
	c.lock.Lock()
	defer c.lock.Unlock()
	close(c.send)
	log.Info("be close ...")
	c.state = false
}
func (c *Client) transfer(msg define.Message) (stop bool) {
	msg.From = c.id
	b, _ := json.Marshal(msg)
	if c.hub.mode == modeOne {
		c.receive(define.HubMsg{
			Data: b,
			From: msg.From,
		})
	} else {
		c.hub.broadcast <- define.HubMsg{
			From: msg.From,
			Data: b,
		}
	}
	return false
}
func (c *Client) receive(b define.HubMsg) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if c.state == false {
		log.Info("broke ...")
		return true
	}
	c.send <- b.Data
	return false
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		err := c.conn.Close()
		if err != nil {
			log.Error("conn,close\t", err)
		}
	}()
	c.conn.SetReadLimit(maxMessageSize)
	err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Error("conn,SetReadDeadline\t", err)
	}
	c.conn.SetPongHandler(func(string) error {
		err = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			log.Error("SetReadDeadline\t", err)
		}
		return nil
	})
	for {
		var message []byte
		_, message, err = c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Errorf("error.IsUnexpectedCloseError: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		messageF, err1 := c.parse(message)
		if err1 != nil {
			continue
		}
		messageF.Msg = strings.TrimSpace(messageF.Msg)
		c.parseRead(messageF)
	}
}
func (c *Client) parseRead(messageF define.Send) {
	if messageF.Types == SendTypePing {
		b, _ := json.Marshal(define.Message{
			Types: define.TypesUserPong,
			Msg:   c.id,
		})
		c.receive(define.HubMsg{
			From: c.id,
			Data: b,
		})
	} else {
		if c.hub.mode == modeGroup {
			b := c.format(messageF)
			if b != nil {
				c.hub.broadcast <- define.HubMsg{
					From: c.id,
					Data: b,
				}
			}
		}
		c.hub.askBot(c, messageF)
	}
}

const (
	SendTypeNormal = 1
	SendTypePing   = 2
)

func (c *Client) parse(message []byte) (jso define.Send, err error) {
	err = json.Unmarshal(message, &jso)
	return
}
func (c *Client) format(jso define.Send) []byte {
	m := define.Message{
		Id:       define.GenID(),
		Err:      nil,
		Eof:      true,
		Msg:      jso.Msg,
		Time:     0,
		YourName: jso.YourName,
		Types:    define.TypesUser,
		From:     c.id,
	}
	b, _ := json.Marshal(m)
	return b
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		err := c.conn.Close()
		if err != nil {
			log.Error("writePump.conn.close\t", err)
		}
	}()
	for {
		select {
		case message, ok := <-c.send:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Error("writePump.SetWriteDeadline\t", err)
			}
			if !ok {
				// The hub closed the channel.
				err = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Error("writePump.WriteMessage\t", err)
				}
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Error("writePump.NextWriter\t", err)
				return
			}
			_, err = w.Write(message)
			if err != nil {
				log.Error("writePump.Write\t", err)
			}

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				_, err = w.Write(newline)
				if err != nil {
					log.Error("writePump.Write,for\t", err)
				}
				res := <-c.send
				_, err = w.Write(res)
				if err != nil {
					log.Error("writePump.Write,for2\t", err)
				}
			}

			if err = w.Close(); err != nil {
				log.Error("writePump.Close\t", err)
				return
			}
		case <-ticker.C:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Error("writePump.SetWriteDeadline\t", err)
			}
			if err = c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Error("writePump.ticker.WriteMessage\t", err)
				return
			}
		}
	}
}

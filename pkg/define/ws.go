package define

import (
	"time"
)

type Message struct {
	Id       int64  `json:"id"`
	Err      error  `json:"err"`
	Eof      bool   `json:"eof"`
	Msg      string `json:"msg"`
	Time     int64  `json:"time"`
	YourName string `json:"your_name"`
	Types    int64  `json:"types"`
	From     string `json:"from"`
}
type Send struct {
	BotCode  string `json:"bot_code"`
	YourName string `json:"your_name"`
	Msg      string `json:"msg"`
	Types    int64  `json:"types"`
}
type HubMsg struct {
	From string `json:"from"`
	Data []byte `json:"data"`
}

type MessagePipe func(q Send, ch chan Message)

func GenID() int64 {
	return time.Now().UnixNano()
}

const (
	TypesSys      = 0
	TypesUser     = 1
	TypesUserPong = 2
)

func DecimalTo26(num int64) string {
	var result string
	for num > 0 {
		temp := num % 26
		if temp == 0 {
			temp = 26
		}
		result = string(rune(temp+64)) + result
		num = (num - temp) / 26
	}

	return result
}

package send

import (
	"log"
	"os"
	"shiro_bot/internal/models/api"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Sender 请求发送
type Sender struct {
	conn      *websocket.Conn
	once      *sync.Once
	done      *chan interface{}
	interrupt *chan os.Signal
	privates  []int64
	groups    []int64
}

// NewSender 新建发送者
func NewSender(conn *websocket.Conn, done *chan interface{}, interrupt *chan os.Signal) Sender {
	once := sync.Once{}

	privates := func() []int64 {
		privateI64s := []int64{}
		priStrs := strings.Split(os.Getenv("PRIVATE"), ",")
		for _, one := range priStrs {
			i64, err := strconv.ParseInt(one, 10, 64)
			if err != nil {
				log.Fatalf("ParseInt err: %v", err)
				return privateI64s
			}
			privateI64s = append(privateI64s, i64)
		}

		return privateI64s
	}()

	groups := func() []int64 {
		groupI64s := []int64{}
		groStrs := strings.Split(os.Getenv("GROUP"), ",")
		for _, one := range groStrs {
			i64, err := strconv.ParseInt(one, 10, 64)
			if err != nil {
				log.Fatalf("ParseInt err: %v", err)
				return groupI64s
			}
			groupI64s = append(groupI64s, i64)
		}

		return groupI64s
	}()

	// TODO: 参数option化
	return Sender{
		conn:      conn,
		once:      &once,
		done:      done,
		interrupt: interrupt,
		privates:  privates,
		groups:    groups,
	}
}

// Handler 处理
func (s *Sender) Handler() {
	// TODO: 所有发送消息的部分在这里处理

	// 启动时执行一次
	s.once.Do(func() {
		s.doOnce()
	})

	for {
		select {
		case <-*s.done:
			return
		// case <-time.After(time.Duration(1) * time.Minute):
		// 定时每分钟发送消息
		// message := ``
		// err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		// if err != nil {
		// 	log.Println("write:", err)
		// 	return
		// }
		case <-*s.interrupt:
			log.Println("Received SIGINT interrupt signal. Closing all pending connections")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := s.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error during closing websocket:", err)
				return
			}
			select {
			case <-*s.done:
				log.Println("Receiver Channel Closed! Exiting....")
			case <-time.After(time.Second):
				log.Println("Timeout in closing receiving channel. Exiting....")
			}
			return
		}
	}
}

// doOnce 执行一次
func (s *Sender) doOnce() {
	//
	message := "shiro酱堂堂复活!!!\n\n" + time.Now().Format(time.RFC3339) + "\n\nDo you wanna build a snowman?"

	s.PrivateMsg(message, s.privates)
	s.GroupMsg(message, s.groups)
}

// PrivateMsg 私聊消息
func (s *Sender) PrivateMsg(message string, privates []int64) {
	// TODO: 拆出来
	sentPrivateMsg := (func() []struct {
		UserID  int64  `json:"user_id"`
		Message string `json:"message"`
	} {
		sent := []struct {
			UserID  int64  `json:"user_id"`
			Message string `json:"message"`
		}{}
		for _, one := range privates {
			sent = append(sent, struct {
				UserID  int64  `json:"user_id"`
				Message string `json:"message"`
			}{
				UserID:  one,
				Message: message,
			})
		}
		return sent
	})()

	for _, one := range sentPrivateMsg {
		err := s.conn.WriteJSON(api.NewRequest("send_private_msg", one, "echo"))
		if err != nil {
			log.Println("write:", err)
			return
		}
	}
}

// GroupMsg 群发消息
func (s *Sender) GroupMsg(message string, groups []int64) {
	// TODO: 拆出来
	sentGroupMsg := (func() []struct {
		GroupID int64  `json:"group_id"`
		Message string `json:"message"`
	} {
		sent := []struct {
			GroupID int64  `json:"group_id"`
			Message string `json:"message"`
		}{}
		for _, one := range groups {
			sent = append(sent, struct {
				GroupID int64  `json:"group_id"`
				Message string `json:"message"`
			}{
				GroupID: one,
				Message: message,
			})
		}
		return sent
	})()

	for _, one := range sentGroupMsg {
		err := s.conn.WriteJSON(api.NewRequest("send_group_msg", one, ""))
		if err != nil {
			log.Println("write:", err)
			return
		}
	}
}

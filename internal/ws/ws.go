package ws

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"shiro_bot/internal/service/receive"
	"shiro_bot/internal/service/send"

	"github.com/gorilla/websocket"
)

// TODO use default option
var (
	// ADDR ws host address
	ADDR = "127.0.0.1:8080"

	// PATH ws route path
	PATH = "/"
)

// Client websocket client
type Client struct {
	u         url.URL // url ws服务器地址
	done      *chan interface{}
	interrupt *chan os.Signal
}

// NewClient ws new client
func NewClient() Client {
	// TODO: 参数option化
	done := make(chan interface{})
	interrupt := make(chan os.Signal, 1)

	u := url.URL{Scheme: "ws", Host: ADDR, Path: PATH}

	wsaddr := os.Getenv("WSADDR")
	if wsaddr != "" {
		u.Host = wsaddr
	}

	wspath := os.Getenv("WSPATH")
	if wspath != "" {
		u.Path = wspath
	}

	client := Client{
		u:         u,
		done:      &done,
		interrupt: &interrupt,
	}

	return client
}

// Connect ws connect
func (client *Client) Connect() {
	// notify system interrupt
	signal.Notify(*client.interrupt, os.Interrupt)

	// ws拨号
	log.Printf("connecting to %s", client.u.String()) // url
	c, _, err := websocket.DefaultDialer.Dial(client.u.String(), nil)
	if err != nil {
		log.Fatal("dial err:", err)
	}
	defer c.Close()

	// receiver.Handler 接收消息处理
	receiver := receive.NewReceiver(c, client.done, client.interrupt)
	go receiver.Handler()

	// sender.Handler 发送消息处理
	sender := send.NewSender(c, client.done, client.interrupt)
	sender.Handler()
}

package receive

import (
	"log"
	"os"

	"github.com/gorilla/websocket"
)

// Receiver 消息接收
type Receiver struct {
	conn      *websocket.Conn
	done      *chan interface{}
	interrupt *chan os.Signal
}

// NewReceiver 新建接收者
func NewReceiver(conn *websocket.Conn, done *chan interface{}, interrupt *chan os.Signal) Receiver {
	// TODO: 参数option化
	return Receiver{
		conn:      conn,
		done:      done,
		interrupt: interrupt,
	}
}

// Handler 处理
func (r *Receiver) Handler() {
	// TODO: 所有接收消息在这里处理
	defer close(*r.done)
	for {
		_, msg, err := r.conn.ReadMessage()
		if err != nil {
			log.Println("Error in receive:", err)
			return
		}
		// TODO: 根据具体消息转到对应事件的处理
		log.Printf("Received: %s\n", msg)
	}
}

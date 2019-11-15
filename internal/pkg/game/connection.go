package game

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Connection interface {
	RunReceive(ws *websocket.Conn, wg *sync.WaitGroup) error
	RunSend(ws *websocket.Conn, wg *sync.WaitGroup) error
	Stop()

	GetUserID() int
	GetSendChan() chan Event
	GetReceiveChan() chan Event
	GetStopSendChan() chan bool
	GetStopReceiveChan() chan bool
}

package connection

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/game"
	"github.com/gorilla/websocket"
	"github.com/op/go-logging"
)

type gameConnection struct {
	userID      int
	sendChan    chan game.Event
	receiveChan chan game.Event
	stopSend    chan bool
	stopReceive chan bool
}

var log = logging.MustGetLogger("connection")

func NewConnection(
	userID int,
	sendChan, receiveChan chan game.Event,
	stopSend, stopReceive chan bool,
) game.Connection {
	return &gameConnection{
		userID:      userID,
		sendChan:    sendChan,
		receiveChan: receiveChan,
		stopSend:    stopSend,
		stopReceive: stopReceive,
	}
}

func (conn *gameConnection) RunReceive(ws *websocket.Conn, wg *sync.WaitGroup) error {
	defer wg.Done()

	for {
		select {
		case <-conn.stopReceive:
			close(conn.receiveChan)
			log.Info("Stopped reading from websocket manually")
			return nil

		default:
			_, msg, err := ws.ReadMessage()
			if err != nil {
				close(conn.receiveChan)
				return err
			}

			var event game.Event
			json.Unmarshal(msg, &event)

			conn.receiveChan <- event
		}
	}
}

func (conn *gameConnection) RunSend(ws *websocket.Conn, wg *sync.WaitGroup) error {
	defer wg.Done()

	ticker := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-conn.stopSend:
			err := ws.WriteMessage(websocket.CloseNormalClosure, []byte{})
			if err != nil {
				return err
			}
			log.Info("Stopped writing into websocket manually")
			return nil

		case event := <-conn.sendChan:
			msg, _ := json.Marshal(event)
			err := ws.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				return err
			}

		case <-ticker.C:
			err := ws.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				return err
			}
		}
	}
}

func (conn *gameConnection) Stop() {
	conn.stopSend <- true
	conn.stopReceive <- true
}

func (conn *gameConnection) GetUserID() int {
	return conn.userID
}

func (conn *gameConnection) GetSendChan() chan game.Event {
	return conn.sendChan
}

func (conn *gameConnection) GetReceiveChan() chan game.Event {
	return conn.receiveChan
}

func (conn *gameConnection) GetStopSendChan() chan bool {
	return conn.stopSend
}

func (conn *gameConnection) GetStopReceiveChan() chan bool {
	return conn.stopReceive
}

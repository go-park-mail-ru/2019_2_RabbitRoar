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
	running     bool
	ws          *websocket.Conn
	wg          sync.WaitGroup
	sendChan    chan game.Event
	receiveChan chan game.EventWrapper
	stop        chan bool
}

var log = logging.MustGetLogger("connection")

func NewConnectionWrapper(
	ws *websocket.Conn,
	sendChan chan game.Event,
	stop chan bool,
) game.ConnectionWrapper {
	return &gameConnection{
		running:  false,
		ws:       ws,
		wg:       sync.WaitGroup{},
		sendChan: sendChan,
		stop:     stop,
	}
}

func (conn *gameConnection) RunReceive(senderID int) {
	conn.wg.Add(1)
	defer conn.wg.Done()

	conn.running = true
	defer func() {
		conn.running = false
	}()

	log.Infof("starting receive goroutine for user %d", senderID)

	conn.receiveChan <- game.EventWrapper{
		SenderID: senderID,
		Event:    &game.Event{
			Type: game.WsRun,
		},
	}

	for {
		select {
		case <- conn.stop:
			return
		default:
			mt, msg, err := conn.ws.ReadMessage()
			if err != nil {
				log.Error("Error reading msg: ", err)
				return
			}

			if mt == websocket.PongMessage {
				continue
			}

			log.Info("Got msg: ", string(msg))

			eventWrap := game.EventWrapper{
				SenderID: senderID,
			}

			err = json.Unmarshal(msg, &eventWrap.Event)
			if err != nil {
				log.Info("Invalid event json received")
			}
			log.Info("Unmarshalled event: ", eventWrap)

			conn.receiveChan <- eventWrap
		}
	}
}

func (conn *gameConnection) RunSend() {
	conn.wg.Add(1)
	defer conn.wg.Done()

	conn.running = true
	defer func() {
		conn.running = false
	}()

	ticker := time.NewTicker(10 * time.Second)

	log.Info("Starting send goroutine for user")

	for {
		select {
		case <-conn.stop:
			err := conn.ws.WriteMessage(websocket.CloseMessage, []byte{})
			if err != nil {
				log.Info("Error sending msg: ", err)
				return
			}

			log.Info("Stopped writing into websocket manually")
			return

		case event := <-conn.sendChan:
			log.Info("Got to send event: ", event)
			msg, err := json.Marshal(event)
			if err != nil {
				log.Info("Error marshalling event: ", err)
			}

			err = conn.ws.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Error("Error sending event: ", err)
				return
			}
			log.Info("Event sent: ", string(msg))

		case <-ticker.C:
			log.Info("Got ticker event, sending ping.")
			err := conn.ws.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				return
			}
		}
	}
}

func (conn *gameConnection) Stop() {
	conn.stop <- true
	conn.stop <- true

	conn.wg.Wait()

	conn.ws.Close()
}

func (conn *gameConnection) IsRunning() bool {
	return conn.running
}

func (conn *gameConnection) SetReceiveChan(rc chan game.EventWrapper) {
	conn.receiveChan = rc
}

func (conn gameConnection) GetSendChan() chan game.Event {
	return conn.sendChan
}

func (conn gameConnection) GetReceiveChan() chan game.EventWrapper {
	return conn.receiveChan
}

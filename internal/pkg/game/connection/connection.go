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
	ws          *websocket.Conn
	wg          sync.WaitGroup
	sendChan    chan game.EventWrapper
	receiveChan chan game.EventWrapper
	stopSend    chan bool
	stopReceive chan bool
}

var log = logging.MustGetLogger("connection")

func NewConnectionWrapper(
	ws *websocket.Conn,
	sendChan, receiveChan chan game.EventWrapper,
	stopSend, stopReceive chan bool,
) game.ConnectionWrapper {
	return &gameConnection{
		ws:          ws,
		wg:          sync.WaitGroup{},
		sendChan:    sendChan,
		receiveChan: receiveChan,
		stopSend:    stopSend,
		stopReceive: stopReceive,
	}
}

func (conn *gameConnection) RunReceive(senderID int) error {
	conn.wg.Add(1)
	defer conn.wg.Done()

	log.Infof("starting receive goroutine for user %d", senderID)

	for {
		log.Info("RECV Loop start")
		select {
		case <-conn.stopReceive:
			close(conn.receiveChan)
			log.Info("Stopped reading from websocket manually")
			return nil

		default:
			_, msg, err := conn.ws.ReadMessage()
			if err != nil {
				log.Error("Error reading msg: ", err)
				close(conn.receiveChan)
				return err
			}
			log.Info("Got msg: ", msg)

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
		log.Info("RECV Loop end")
	}
}

func (conn *gameConnection) RunSend() error {
	conn.wg.Add(1)
	defer conn.wg.Done()

	ticker := time.NewTicker(10 * time.Second)

	log.Info("starting send goroutine for user")

	for {
		log.Info("SEND Loop start")
		select {
		case <-conn.stopSend:
			err := conn.ws.WriteMessage(websocket.CloseMessage, []byte{})
			if err != nil {
				log.Info("Error sending msg: ", err)
				return err
			}

			log.Info("Stopped writing into websocket manually")
			return nil

		case event := <-conn.sendChan:
			log.Info("Got to send event: ", event)
			msg, err := json.Marshal(event)
			log.Info("Error marshalling event: ", err)
			err = conn.ws.WriteMessage(websocket.TextMessage, msg)
			log.Info("Event sent: ", msg)
			if err != nil {
				log.Error("Error sending event: ", err)
				return err
			}

		case <-ticker.C:
			log.Info("Got ticker event, sending ping.")
			err := conn.ws.WriteMessage(websocket.TextMessage, []byte{'p', 'i', 'n', 'g'})
			if err != nil {
				return err
			}
		}
		log.Info("SEND Loop end")
	}
}

func (conn *gameConnection) Stop() {
	conn.stopSend <- true
	conn.stopReceive <- true

	conn.wg.Wait()

	conn.ws.Close()
}

func (conn gameConnection) GetSendChan() chan game.EventWrapper {
	return conn.sendChan
}

func (conn gameConnection) GetReceiveChan() chan game.EventWrapper {
	return conn.receiveChan
}

func (conn gameConnection) GetStopSendChan() chan bool {
	return conn.stopSend
}

func (conn gameConnection) GetStopReceiveChan() chan bool {
	return conn.stopReceive
}

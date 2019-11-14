package connection

import "github.com/go-park-mail-ru/2019_2_RabbitRoar/internal/pkg/game"

type gamePlayerConnection struct {
	sendChan    chan []byte
	receiveChan chan []byte
	stopSend    chan bool
	stopReceive chan bool
}

func NewConnection(
	sendChan, receiveChan chan []byte,
	stopSend, stopReceive chan bool,
) game.PlayerConnection {
	return &gamePlayerConnection{
		sendChan:    sendChan,
		receiveChan: receiveChan,
		stopSend:    stopSend,
		stopReceive: stopReceive,
	}
}

func (conn *gamePlayerConnection) GetSendChan() chan []byte {
	return conn.sendChan
}

func (conn *gamePlayerConnection) GetReceiveChan() chan []byte {
	return conn.receiveChan
}

func (conn *gamePlayerConnection) GetStopSendChan() chan bool {
	return conn.stopSend
}

func (conn *gamePlayerConnection) GetStopReceiveChan() chan bool {
	return conn.stopReceive
}

func (conn *gamePlayerConnection) Stop() {
	conn.stopSend <- true
	conn.stopReceive <- true
}

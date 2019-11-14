package game

type PlayerConnection interface {
	GetSendChan() chan []byte
	GetReceiveChan() chan []byte
	GetStopSendChan() chan bool
	GetStopReceiveChan() chan bool

	Stop()
}

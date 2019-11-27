package game

type ConnectionWrapper interface {
	RunReceive(senderID int)
	RunSend()
	Stop()
	IsRunning() bool

	SetReceiveChan(rc chan EventWrapper)
	GetSendChan() chan Event
	GetReceiveChan() chan EventWrapper
}

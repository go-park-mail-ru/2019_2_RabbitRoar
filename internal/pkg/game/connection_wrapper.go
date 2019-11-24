package game

type ConnectionWrapper interface {
	RunReceive(senderID int)
	RunSend()
	Stop()

	SetReceiveChan(rc chan EventWrapper)
	GetSendChan() chan Event
	GetReceiveChan() chan EventWrapper
}

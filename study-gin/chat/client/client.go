package client

import (
	"github.com/ScottAI/chatserver/protocol"
)

type Client interface {
	Dial(address string) error
	Start()
	Close()
	Send(commond interface{}) error
	SetName(name string) error
	SendMess(message string) error
	InComming() chan protocol.MessCmd
}


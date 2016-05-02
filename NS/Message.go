package NS

import (
	"net"
	"errors"
)

type Msg struct {
	Length uint32
	Body   []byte
}

type MsgSender func(net.Conn, Msg) error
type MsgProcessor func(net.Conn, []byte) error
type OnConnected func(TCPConnect) error
type OnDisconnect func(TCPConnect) error

var NoProcessorError = errors.New("No Processor")
var NoSenderError = errors.New("No Sender")
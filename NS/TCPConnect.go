package NS

import (
	"io"
	"net"
	"sync"
	"encoding/binary"
)

type OnStart func()error
type OnEnd func()error

type TCPConnect interface {
	io.Closer
	RegisterMsgSender(MsgSender)
	RegisterMsgProcessor(MsgProcessor)
	RegisterOnStart(OnStart)
	RegisterOnEnd(OnEnd)
	RegisterOnConnected(OnConnected)
	RegisterOnDisconnected(OnDisconnect)

	Send([]byte) error
	RunLoop()
}

type tcpConnect_impl struct {
	connect net.TCPConn
	sender MsgSender
	processor MsgProcessor
	onStart OnStart
	onEnd OnEnd
	isRunning bool
	lock sync.RWMutex
}

func (self *tcpConnect_impl) RegisterMsgSender(sender MsgSender) {
	self.sender = sender
}
func (self *tcpConnect_impl) RegisterMsgProcessor(processor MsgProcessor) {
	self.processor = processor
}
func (self *tcpConnect_impl) RegisterOnStart(onStart OnStart){
	self.onStart = onStart
}
func (self *tcpConnect_impl) RegisterOnEnd(onEnd OnEnd){
	self.onEnd = onEnd
}

func (self *tcpConnect_impl) Send(body []byte) error {
	if nil == self.sender {
		return NoSenderError
	}
	return self.sender(self.connect, body)
}
func (self *tcpConnect_impl) RunLoop() {
	self.lock.Lock()
	self.isRunning = true
	self.lock.Unlock()

	if nil != self.onStart {
		err := self.onStart()
		if nil != err {
			panic(err)
		}
	}
	if nil != self.onEnd {
		defer func() {
			err := self.onEnd()
			if nil != err {
				panic(err)
			}
		}()
	}
	if nil == self.processor {
		panic(NoProcessorError)
	}
	for {
		self.lock.RLock()
		if !self.isRunning {
			break
		}
		self.lock.RUnlock()

		var length uint32
		binary.Read(self.connect, binary.BigEndian, &length)

		body := make([]byte, length)
		_, err := io.ReadFull(self.connect, body)
		if nil != err {
			break
		}
		self.processor(self.connect, body)
	}
}
func (self *tcpConnect_impl) Close() error {
	self.lock.Lock()
	self.lock.Unlock()
	self.isRunning = false
	return nil
}

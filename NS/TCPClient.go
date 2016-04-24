package NS

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

type TCPClient interface {
	io.Closer
	Connect(string, int) error
	RunLoop() error
	Send(uint16, []byte) error
	RegisterProcessor(MessageProcessor) error
	RegisterOnConnected(func() error)
	RegisterOnDisconnect(func() error)
}

func CreateTCPClient(ip string, port int) (TCPClient, error) {
	client := new(tcpClient_impl)
	err := client.Connect(ip, port)
	return client, err
}

type tcpClient_impl struct {
	remoteAddress *net.TCPAddr
	isRunning     bool
	connect       *net.TCPConn

	processor    MessageProcessor
	onConnected  func() error
	onDisconnect func() error
}

func (self *tcpClient_impl) Connect(ip string, port int) (err error) {
	self.remoteAddress, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ip, port))
	return err
}
func (self *tcpClient_impl) RunLoop() (err error) {
	self.connect, err = net.DialTCP(self.remoteAddress.Network(), nil, self.remoteAddress)
	if nil != err {
		return err
	}
	defer self.connect.Close()
	self.isRunning = true
	self.onConnected()
	defer func() {
		if nil != self.onDisconnect {
			self.onDisconnect()
		}
	}()
	for self.isRunning && nil == err {
		err = processConnect(self.connect, self.processor)
	}
	return err
}
func (self *tcpClient_impl) Close() error {
	self.isRunning = false
	return nil
}
func (self *tcpClient_impl) RegisterProcessor(processor MessageProcessor) error {
	self.processor = processor
	return nil
}
func (self *tcpClient_impl) RegisterOnConnected(onConnected func() error) {
	self.onConnected = onConnected
}
func (self *tcpClient_impl) RegisterOnDisconnect(onDisconnect func() error) {
	self.onDisconnect = onDisconnect
}
func (self *tcpClient_impl) Send(msgType uint16, msgBytes []byte) error {
	var msg Message
	msg.Body = msgBytes
	msg.TypeID = msgType
	msg.Length = uint16(len(msg.Body) + 8)
	binary.Write(self.connect, binary.BigEndian, msg.Length)
	binary.Write(self.connect, binary.BigEndian, msg.TypeID)
	self.connect.Write(msg.Body)
	//SendMessage(self.connect, msg)
	return nil
}

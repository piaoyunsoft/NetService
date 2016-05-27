package NS

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
)

type IParent interface {
	OnConnectClosed(connect TCPConnect)
}

type TCPConnect interface {
	SetOnConnectedCallback(callback func(connect TCPConnect) (err error))                   //设置连接成功回调
	SetOnReceiveMsgCallback(callback func(connect TCPConnect, msgBytes []byte) (err error)) //设置收到消息回调（会传递一个完整的消息体）
	SetOnClosedCallback(callback func(connect TCPConnect) (err error))                      //设置连接主动关闭回调
	Send(msgBytes []byte) error
	io.Closer
}

type tcpConnect_impl struct {
	parent    IParent
	connect   *net.TCPConn
	isRunning bool

	onConnect func(connect TCPConnect) (err error)
	onReceive func(connect TCPConnect, msgBytes []byte) (err error)
	onClosed  func(connect TCPConnect) (err error)
}

func (self *tcpConnect_impl) SetOnConnectedCallback(callback func(connect TCPConnect) (err error)) {
	self.onConnect = callback
}
func (self *tcpConnect_impl) SetOnReceiveMsgCallback(callback func(connect TCPConnect, msgBytes []byte) (err error)) {
	self.onReceive = callback
}
func (self *tcpConnect_impl) SetOnClosedCallback(callback func(connect TCPConnect) (err error)) {
	self.onClosed = callback
}
func (self *tcpConnect_impl) Send(msgBytes []byte) error {
	buffer := bytes.NewBuffer(nil)
	length := uint32(len(msgBytes))
	binary.Write(buffer, binary.LittleEndian, &length)
	buffer.Write(msgBytes)
	_, err := self.connect.Write(buffer.Bytes())
	return err
}
func (self *tcpConnect_impl) Close() error {
	self.isRunning = false
	return self.connect.Close()
}

func (self *tcpConnect_impl) init(parent IParent, connect *net.TCPConn) {
	self.parent = parent
	self.connect = connect
	self.isRunning = true
}
func (self *tcpConnect_impl) runOnce() (err error) {
	var length uint32
	err = binary.Read(self.connect, binary.LittleEndian, &length)
	if nil != err {
		return err
	}
	buffer := make([]byte, length)
	_, err = io.ReadFull(self.connect, buffer)
	if nil != err {
		return err
	}
	return self.onReceive(self, buffer)
}
func (self *tcpConnect_impl) run() {
	if nil != self.onConnect {
		err := self.onConnect(self)
		if nil != err {
			return
		}
	}
	if nil != self.onClosed {
		defer self.onClosed(self)
	}
	if nil != self.parent {
		defer self.parent.OnConnectClosed(self)
	}

	for self.isRunning {
		err := self.runOnce()
		if nil != err {
			break
		}
	}
}

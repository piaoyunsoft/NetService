package NS

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/SailorKGame/SimpleLog/SLog"
	"net"
	"sync"
)

type ReceiveCallback func(*TCPConnect, []byte) error
type OnConnectedCallback func(*TCPConnect) error

var NoConnectError = errors.New("Haven't created tcp connect!")
var NotReadEnoughBytesError = errors.New("Not read enough bytes!")
var NotSetReceiveCallbackError = errors.New("Not set receive callback!")

type TCPConnect struct {
	lock                sync.RWMutex
	connect             *net.TCPConn
	security            ISecurity
	onConnectedCallback OnConnectedCallback
	receiveCallback     ReceiveCallback
}

func (self *TCPConnect) setConnect(connect *net.TCPConn) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.connect = connect
}
func (self TCPConnect) GetConnect() *net.TCPConn {
	return self.connect
}

func (self *TCPConnect) SetOnConnectedCallback(callback OnConnectedCallback) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.onConnectedCallback = callback
}

func (self *TCPConnect) SetReceiveCallback(callback ReceiveCallback) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.receiveCallback = callback
}

func (self TCPConnect) getReceiveCallback() ReceiveCallback {
	return self.receiveCallback
}

func (self *TCPConnect) SetSecurity(security ISecurity) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.security = security
}

func (self TCPConnect) getSecurity() ISecurity {
	return self.security
}

func (self *TCPConnect) RunOnce() error {
	connect := self.GetConnect()
	if nil == connect {
		return NoConnectError
	}
	var length uint32
	err := binary.Read(connect, binary.LittleEndian, &length)
	if nil != err {
		return err
	}
	buffer := make([]byte, length)
	n, err := connect.Read(buffer)
	if nil != err {
		return err
	}
	if n != int(length) {
		return NotReadEnoughBytesError
	}
	security := self.getSecurity()
	if nil != security {
		buffer, err = security.Decrypt(buffer)
		if nil != err {
			return err
		}
	}
	callback := self.getReceiveCallback()
	if nil != callback {
		return callback(self, buffer)
	} else {
		return NotSetReceiveCallbackError
	}
}

func (self *TCPConnect) Write(msgBytes []byte) (n int, err error) {
	connect := self.GetConnect()
	if nil == connect {
		return 0, NoConnectError
	}
	security := self.getSecurity()
	if nil != security {
		msgBytes, err = security.Encrypt(msgBytes)
		if nil != err {
			return 0, err
		}
	}
	buffer := bytes.NewBuffer(nil)
	err = binary.Write(buffer, binary.LittleEndian, uint32(len(msgBytes)))
	if nil != err {
		return 0, err
	}
	n, err = buffer.Write(msgBytes)
	if nil != err {
		return 0, err
	}
	return connect.Write(buffer.Bytes())
}

func RunTCPConnect(connect *TCPConnect) {
	err := connect.onConnectedCallback(connect)
	if nil != err {
		SLog.E("TCPConnect", err)
		return
	}
	for {
		err := connect.RunOnce()
		if nil == err {
			continue
		}
		SLog.W("TCPConnect", err)
		connect.setConnect(nil)
		break
	}
}

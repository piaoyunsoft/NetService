package NS

import (
	"net"
	"encoding/binary"
	"bytes"
	"errors"
	"sync"
)

var NoConnectError = errors.New("Haven't created tcp connect!")
var NotReadEnoughBytesError = errors.New("Not read enough bytes!")
var NotSetReceiveCallbackError = errors.New("Not set receive callback!")

type TCPConnect struct {
	lock sync.RWMutex
	connect  net.TCPConn
	security ISecurity
	receiveCallback ReceiveCallback
}

func (self *TCPConnect)SetConnect(connect net.TCPConn) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.connect = connect
}
func (self *TCPConnect)getConnect() net.TCPConn {
	self.lock.RLock()
	defer self.lock.RUnlock()
	return self.connect
}

func (self *TCPConnect)SetReceiveCallback(callback ReceiveCallback) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.receiveCallback = callback
}

func (self *TCPConnect)getReceiveCallback() ReceiveCallback {
	self.lock.RLock()
	defer self.lock.RUnlock()
	return self.receiveCallback
}

func (self *TCPConnect) SetSecurity(security ISecurity) {
	self.lock.Lock()
	defer self.lock.Unlock()
	self.security = security
}

func (self TCPConnect)getSecurity() ISecurity {
	self.lock.RLock()
	defer self.lock.RUnlock()
	return self.security
}

func (self *TCPConnect) RunOnce() error {
	connect := self.getConnect()
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
	if n != length {
		return NotReadEnoughBytesError
	}
	security := self.getSecurity()
	if nil != security {
		buffer = security.Decrypt(buffer)
	}
	callback := self.getReceiveCallback()
	if nil != callback {
		return callback(buffer)
	} else {
		return NotSetReceiveCallbackError
	}
}

func (self *TCPConnect)Write(msgBytes []byte) (n int, err error) {
	connect := self.getConnect()
	if nil == connect {
		return 0, NoConnectError
	}
	security := self.getSecurity()
	if nil != security {
		msgBytes = security.Encrypt(msgBytes)
	}
	buffer := bytes.NewBuffer(nil)
	err = binary.Write(buffer, binary.LittleEndian, len(msgBytes))
	if nil != err {
		return 0, err
	}
	n, err = buffer.Write(msgBytes)
	if nil != err {
		return 0, err
	}
	return connect.Write(buffer)
}
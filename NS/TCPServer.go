package NS

import (
	"net"
	"errors"
	"sync"
)

var NoListenerError = errors.New("No listener has created!")

type TCPServer struct {
	lock sync.RWMutex
	address net.TCPAddr
	listener net.TCPListener
}

func (self *TCPServer)SetAddress(address string) (err error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.address, err = net.ResolveTCPAddr("tcp", address)
	return err
}

func (self *TCPServer)Listen() (err error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.listener, err = net.ListenTCP("tcp", self.address)
	return err
}

func (self *TCPServer)RunOnce() (TCPConnect, error) {
	self.lock.RLock()
	defer self.lock.RUnlock()

	if nil == self.listener {
		return NoListenerError
	}
	connect, err := self.listener.AcceptTCP()
	if nil != err {
		return err
	}
	tcpConnect := new(TCPConnect)
	tcpConnect.SetConnect(connect)
	return tcpConnect, nil
}


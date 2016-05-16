package NS

import (
	"errors"
	"github.com/SailorKGame/SimpleLog/SLog"
	"net"
	"sync"
)

var NoListenerError = errors.New("No listener has created!")

type TCPServer struct {
	lock     sync.RWMutex
	address  *net.TCPAddr
	listener *net.TCPListener
}

func (self *TCPServer) SetAddress(address string) (err error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.address, err = net.ResolveTCPAddr("tcp", address)
	return err
}

func (self *TCPServer) ClearListener() {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.listener = nil
}

func (self *TCPServer) Listen() (err error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.listener, err = net.ListenTCP("tcp", self.address)
	return err
}

func (self *TCPServer) RunOnce() (*TCPConnect, error) {
	self.lock.RLock()
	defer self.lock.RUnlock()

	if nil == self.listener {
		return nil, NoListenerError
	}
	connect, err := self.listener.AcceptTCP()
	if nil != err {
		return nil, err
	}
	tcpConnect := new(TCPConnect)
	tcpConnect.setConnect(connect)
	return tcpConnect, nil
}

func CreateTCPServer(address string) (*TCPServer, error) {
	server := new(TCPServer)
	err := server.SetAddress(address)
	return server, err
}

func RunTCPServer(server *TCPServer, onAccept func(*TCPConnect) error) {
	for {
		connect, err := server.RunOnce()
		if nil == err {
			err = onAccept(connect)
			if nil != err {
				SLog.W("TCPServer", err)
			}
			continue
		}
		if NoListenerError == err {
			err = server.Listen()
			if nil != err {
				SLog.W("TCPServer", err)
			}
			continue
		}
		SLog.W("TCPServer", err)
		server.ClearListener()
	}
}

package NS

import (
	"fmt"
	"github.com/sunny352/SimpleLog/SLog"
	"io"
	"net"
)

type TCPServer interface {
	io.Closer
	Listen(string, int) error
	RunLoop() error
	RegisterProcessor(MessageProcessor) error
}

func CreateTCPServer(ip string, port int) (TCPServer, error) {
	server := new(tcpServer_impl)
	err := server.Listen(ip, port)
	return server, err
}

type tcpServer_impl struct {
	tcpAddress *net.TCPAddr
	listener   *net.TCPListener
	processor  MessageProcessor
	isRunning  bool
}

func (self *tcpServer_impl) Listen(ip string, port int) (err error) {
	self.tcpAddress, err = net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ip, port))
	return err
}

func (self *tcpServer_impl) RegisterProcessor(processor MessageProcessor) error {
	self.processor = processor
	return nil
}

func (self *tcpServer_impl) Close() error {
	self.isRunning = false
	return nil
}

func (self *tcpServer_impl) RunLoop() (err error) {
	self.listener, err = net.ListenTCP("tcp", self.tcpAddress)
	if nil != err {
		return err
	}
	defer self.listener.Close()
	self.isRunning = true
	for self.isRunning {
		connect, err := self.listener.AcceptTCP()
		if nil != err {
			SLog.E("TCPServer", err)
			continue
		}
		go self.processConnect(connect)
	}
	return err
}

func (self tcpServer_impl) processConnect(connect net.Conn) (err error) {
	SLog.I("TCPServer", "Process", connect.RemoteAddr())
	defer connect.Close()
	for self.isRunning && nil == err {
		err = processConnect(connect, self.processor)
	}
	if nil != err {
		SLog.E("TCPServer", err)
	} else {
		SLog.I("TCPServer", "Closed")
	}
	return err
}

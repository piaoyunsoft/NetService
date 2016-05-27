package NS

import (
	"io"
	"net"
)

type TCPServer interface {
	Start() (err error)                                                                  //启动监听，开启线程
	SetOnAcceptCallback(callback func(server TCPServer, connect TCPConnect) (err error)) //设置收到连接后的处理
	io.Closer
}

func CreateTCPServer(name string, address string) (server TCPServer, err error) {
	serverImpl := new(tcpServer_impl)
	err = serverImpl.init(name, address)
	if nil != err {
		return nil, err
	}
	return serverImpl, nil
}

type tcpServer_impl struct {
	name      string
	address   *net.TCPAddr
	isRunning bool
	listener  *net.TCPListener
	onAccept  func(server TCPServer, connect TCPConnect) (err error)
}

func (self *tcpServer_impl) init(name string, address string) (err error) {
	self.name = name
	self.address, err = net.ResolveTCPAddr("tcp", address)
	return err
}
func (self *tcpServer_impl) Start() (err error) {
	self.isRunning = true
	go self.run()
	return nil
}
func (self *tcpServer_impl) SetOnAcceptCallback(callback func(server TCPServer, connect TCPConnect) (err error)) {
	self.onAccept = callback
}
func (self *tcpServer_impl) Close() error {
	self.isRunning = false
	return self.listener.Close()
}

func (self *tcpServer_impl) OnConnectClosed(connect TCPConnect) {

}

func (self *tcpServer_impl) run() {
	var err error
	self.listener, err = net.ListenTCP("tcp", self.address)
	if nil != err {
		return
	}
	for self.isRunning {
		connect, err := self.listener.AcceptTCP()
		if nil != err {
			continue
		}
		tcpConnect := new(tcpConnect_impl)
		tcpConnect.init(self, connect)
		self.onAccept(self, tcpConnect)
		go tcpConnect.run()
	}
}

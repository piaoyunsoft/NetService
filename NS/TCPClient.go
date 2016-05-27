package NS

import (
	"github.com/SailorKGame/SimpleLog/SLog"
	"io"
	"net"
)

type TCPClient interface {
	Start() //发起连接，创建线程
	SetOnConnectedCallback(callback func(client TCPClient, connect TCPConnect) (err error))
	GetConnect() TCPConnect
	io.Closer
}

func CreateTCPClient(name string, address string) (client TCPClient, err error) {
	clientImpl := new(tcpClient_impl)
	err = clientImpl.init(name, address)
	if nil != err {
		return nil, err
	}
	return clientImpl, nil
}

type tcpClient_impl struct {
	name            string
	address         *net.TCPAddr
	isAutoReconnect bool
	connect         TCPConnect
	onConnected     func(client TCPClient, connect TCPConnect) (err error)
}

func (self *tcpClient_impl) init(name string, address string) (err error) {
	self.name = name
	self.address, err = net.ResolveTCPAddr("tcp", address)
	self.isAutoReconnect = true
	return err
}
func (self *tcpClient_impl) Start() {
	go self.run()
}
func (self *tcpClient_impl) SetOnConnectedCallback(callback func(client TCPClient, connect TCPConnect) (err error)) {
	self.onConnected = callback
}
func (self *tcpClient_impl) GetConnect() TCPConnect {
	return self.connect
}
func (self *tcpClient_impl) Close() (err error) {
	self.isAutoReconnect = false
	return self.connect.Close()
}

func (self *tcpClient_impl) OnConnectClosed(connect TCPConnect) {

}

func (self *tcpClient_impl) run() {
	for self.isAutoReconnect {
		connect, err := net.DialTCP("tcp", nil, self.address)
		if nil != err {
			SLog.W("TCPClient", err)
			continue
		}
		tcpConnect := new(tcpConnect_impl)
		tcpConnect.init(self, connect)
		self.connect = tcpConnect
		self.onConnected(self, tcpConnect)
		tcpConnect.run()
		self.connect = nil
	}
}

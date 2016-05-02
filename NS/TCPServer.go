package NS

import (
	"io"
	"net"
)

type OnAccept func(net.Conn) error

type TCPServer interface {
	io.Closer
	RunLoop()
	RegisterOnAccept(OnAccept)
}

func CreateTCPServer(address string) TCPServer {
	server := new(tcpServer_impl)
	return server
}

type tcpServer_impl struct {
	address net.TCPAddr
}

func (self *tcpServer_impl) init(address string) (err error) {
	self.address, err = net.ResolveTCPAddr("tcp", address)
	return err
}

func (self *tcpServer_impl) RunLoop() {

}
func (self *tcpServer_impl) RegisterOnAccept(callback OnAccept) {

}
func (self *tcpServer_impl) Close() error {
	return nil
}

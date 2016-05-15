package NS

import (
	"errors"
	"net"
	"github.com/SailorKGame/SimpleLog/SLog"
)

var NoAddressError = errors.New("No address to connect!")

type ReceiveCallback func([]byte) error

type TCPClient struct {
	address net.TCPAddr
	isAutoReconnect bool
	TCPConnect
}

func (self *TCPClient) SetAddress(address string) (err error) {
	self.address, err = net.ResolveTCPAddr("tcp", address)
	return err
}

func (self TCPClient) GetAddress() net.TCPAddr {
	return self.address
}

func (self *TCPClient)SetAutoReconnect(isAuto bool) {
	self.isAutoReconnect = isAuto
}
func (self *TCPClient)IsAutoReconnect() bool {
	return self.isAutoReconnect
}

func (self *TCPClient) Connect() (err error) {
	if nil == self.address {
		return NoAddressError
	}
	self.connect, err = net.DialTCP("tcp", nil, self.address)
	return err
}

func CreateTCPClient(address string) *TCPClient {
	client := new(TCPClient)
	client.SetAddress(address)
	return client
}

func RunTCPClient(client *TCPClient) {
	for {
		err := client.RunOnce()
		if nil == err {
			continue
		}
		if NoConnectError == err {
			err == client.Connect()
			if nil == err {
				continue
			}
		}
		SLog.W("TCPClient", err)
		client.SetConnect(nil)
		if !client.IsAutoReconnect() {
			break
		}
	}
}

package NS

import (
	"errors"
	"github.com/SailorKGame/SimpleLog/SLog"
	"net"
)

var NoAddressError = errors.New("No address to connect!")

type TCPClient struct {
	address         *net.TCPAddr
	isAutoReconnect bool
	TCPConnect
}

func (self *TCPClient) setAddress(address string) (err error) {
	self.address, err = net.ResolveTCPAddr("tcp", address)
	return err
}

func (self TCPClient) GetAddress() *net.TCPAddr {
	return self.address
}

func (self *TCPClient) SetAutoReconnect(isAuto bool) {
	self.isAutoReconnect = isAuto
}
func (self *TCPClient) IsAutoReconnect() bool {
	return self.isAutoReconnect
}

func (self *TCPClient) Connect() (err error) {
	if nil == self.address {
		return NoAddressError
	}
	self.connect, err = net.DialTCP("tcp", nil, self.address)
	return err
}

func CreateTCPClient(address string) (*TCPClient, error) {
	client := &TCPClient{isAutoReconnect: true}
	err := client.setAddress(address)
	return client, err
}

func RunTCPClient(client *TCPClient) {
	for {
		err := client.RunOnce()
		if nil == err {
			continue
		}
		if NoConnectError == err {
			SLog.W("TCPClient", "No connected")
			err = client.Connect()
			if nil == err {
				err = client.onConnectedCallback(&client.TCPConnect)
				if nil != err {
					SLog.W("TCPClient", err)
				}
				continue
			}
		}
		SLog.W("TCPClient", err)
		client.setConnect(nil)
		if !client.IsAutoReconnect() {
			break
		}
	}
}

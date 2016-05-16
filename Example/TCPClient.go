package main

import (
	"github.com/SailorKGame/NetService/NS"
	"github.com/SailorKGame/SimpleLog/SLog"
)

func main() {
	SLog.D("Client", "Start")
	client, err := NS.CreateTCPClient("127.0.0.1:8123")
	if nil != err {
		SLog.E("Client", err)
		return
	}
	SLog.D("Client", "Created client")
	client.SetOnConnectedCallback(func(connect *NS.TCPConnect) error {
		n, err := connect.Write([]byte("test"))
		SLog.D("Client", "OnConnected", n)
		return err
	})
	client.SetReceiveCallback(func(connect *NS.TCPConnect, msgBytes []byte) error {
		SLog.I("Client", string(msgBytes), connect.GetConnect().RemoteAddr())
		return nil
	})
	NS.RunTCPClient(client)
}

package main

import (
	"fmt"
	"github.com/sunny352/NetService/NS"
	"github.com/sunny352/SimpleLog/SLog"
	"net"
)

func main() {
	SLog.I("NetService", "Start")
	client, err := NS.CreateTCPClient("127.0.0.1", 9527)
	if nil != err {
		SLog.E("NetService", err)
		return
	}
	SLog.I("NetService", "Created TCPClient")
	client.RegisterProcessor(func(connect net.Conn, message NS.Message) error {
		SLog.D("Client", connect.RemoteAddr(), fmt.Sprintf("length:%d type:%d body:%x \"%s\"", message.Length, message.TypeID, message.Body, string(message.Body)))
		return nil
	})
	client.RegisterOnConnected(func() error {
		client.Send(1, []byte("test"))
		SLog.I("NetService", "Sended")
		return nil
	})
	SLog.I("NetService", "TCPClient", "RegisterProcessor")
	client.RunLoop()
}

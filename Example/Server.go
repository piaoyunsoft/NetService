package main

import (
	"fmt"
	"github.com/sunny352/NetService/NS"
	"github.com/sunny352/SimpleLog/SLog"
	"net"
)

func main() {
	SLog.I("NetService", "Start")
	server, err := NS.CreateTCPServer("127.0.0.1", 9527)
	if nil != err {
		SLog.E("NetService", err)
		return
	}
	SLog.I("NetService", "Created TCPServer")
	server.RegisterProcessor(func(connect net.Conn, message NS.Message) error {
		SLog.D("Server", connect.RemoteAddr(), fmt.Sprintf("length:%d type:%d body:%x \"%s\"", message.Length, message.TypeID, message.Body[:], string(message.Body)))
		NS.SendMessage(connect, message)
		return nil
	})
	SLog.I("NetService", "TCPServer", "RegisterProcessor")
	server.RunLoop()
}

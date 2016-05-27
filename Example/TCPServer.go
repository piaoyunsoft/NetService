package main

import (
	"github.com/SailorKGame/NetService/NS"
	"github.com/SailorKGame/SimpleLog/SLog"
	"runtime"
)

func main() {
	SLog.I("Test", "Start Server")
	server, err := NS.CreateTCPServer("TestServer", ":9631")
	if nil != err {
		SLog.E("Test", err)
		return
	}
	server.SetOnAcceptCallback(func(server NS.TCPServer, connect NS.TCPConnect) (err error) {
		connect.SetOnConnectedCallback(func(connect NS.TCPConnect) (err error) {
			SLog.D("Test", "OnConnected")
			connect.Send([]byte("ServerToClient"))
			return nil
		})
		connect.SetOnReceiveMsgCallback(func(connect NS.TCPConnect, msgBytes []byte) (err error) {
			SLog.D("Test", "OnReceive", string(msgBytes))
			return nil
		})
		connect.SetOnClosedCallback(func(connect NS.TCPConnect) (err error) {
			SLog.D("Test", "OnClosed")
			return nil
		})
		return nil
	})
	server.Start()

	for {
		runtime.Gosched()
	}
}

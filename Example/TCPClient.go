package main

import (
	"github.com/SailorKGame/NetService/NS"
	"github.com/SailorKGame/SimpleLog/SLog"
	"runtime"
)

func main() {
	SLog.I("Test", "StartClient")
	client, err := NS.CreateTCPClient("TestClient", "127.0.0.1:9631")
	if nil != err {
		SLog.E("Test", err)
		return
	}
	client.SetOnConnectedCallback(func(client NS.TCPClient, connect NS.TCPConnect) (err error) {
		connect.SetOnConnectedCallback(func(connect NS.TCPConnect) (err error) {
			SLog.D("Test", "OnConnected")
			connect.Send([]byte("ClientToServer"))
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
	client.Start()

	for {
		runtime.Gosched()
	}
}

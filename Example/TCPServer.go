package main

import (
	"github.com/SailorKGame/NetService/NS"
	"github.com/SailorKGame/SimpleLog/SLog"
)

func main() {
	SLog.D("Server", "Start")
	server, err := NS.CreateTCPServer(":8123")
	if nil != err {
		SLog.E("Server", err)
		return
	}
	NS.RunTCPServer(server, func(connect *NS.TCPConnect) error {
		connect.SetOnConnectedCallback(func(connect *NS.TCPConnect) error {
			SLog.D("Server", connect.GetConnect().RemoteAddr())
			return nil
		})
		connect.SetReceiveCallback(func(connect *NS.TCPConnect, msgBytes []byte) error {
			SLog.I("Server", string(msgBytes))
			_, err := connect.Write([]byte("Server Echo"))
			return err
		})
		go NS.RunTCPConnect(connect)
		return nil
	})
}

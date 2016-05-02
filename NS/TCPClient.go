package NS

type TCPClient interface {
	TCPConnect
	Reconnect() error
}

func CreateTCPClient(address string) TCPConnect {
	client := new(tcpClient_impl)
	return client
}

type tcpClient_impl struct {
	connect TCPConnect
}

func (self *tcpClient_impl) RegisterMsgSender(MsgSender) {

}
func (self *tcpClient_impl) RegisterMsgProcessor(MsgProcessor) {

}
func (self *tcpClient_impl) RegisterOnStart(OnStart){

}
func (self *tcpClient_impl) RegisterOnEnd(OnEnd){

}
func (self *tcpClient_impl) Send([]byte) error {
	return nil
}
func (self *tcpClient_impl) RunLoop() {

}
func (self *tcpClient_impl) Reconnect() error {
	return nil
}
func (self *tcpClient_impl) Close() error {
	return nil
}

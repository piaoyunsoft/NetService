package NS

import (
	"net"
	"github.com/SailorKGame/SimpleLog/SLog"
	"encoding/binary"
	"io"
)

func processConnect(connect net.Conn, processor MessageProcessor) (err error) {
	SLog.D("NetService", "Process", connect.RemoteAddr())
	var length uint16
	err = binary.Read(connect, binary.BigEndian, &length)
	if nil != err {
		SLog.E("NetService", err)
		return err
	} else {
		SLog.D("NetService", "length is", length)
	}
	var msgType uint16
	err = binary.Read(connect, binary.BigEndian, &msgType)
	if nil != err {
		SLog.E("NetService", err)
		return err
	} else {
		SLog.D("NetService", "TypeID is", msgType)
	}
	buff := make([]byte, length-8)
	n, err := io.ReadFull(connect, buff[:])
	if nil != err {
		SLog.E("NetService", err)
		return err
	}
	if n != len(buff) {
		err = NotEnoughBytes
		SLog.E("NetService", err)
		return err
	}
	var msg Message
	msg.Length = length
	msg.TypeID = msgType
	msg.Body = buff
	if nil != processor {
		err = processor(connect, msg)
	} else {
		err = NoProcessorError
	}
	if nil != err {
		SLog.E("NetService", err)
	}
	return err
}

func SendMessage(connect net.Conn, msg Message) error {
	SLog.D("NetService", connect.RemoteAddr(), msg.Length, msg.TypeID)
	return binary.Write(connect, binary.BigEndian, &msg)
}
package main

import (
	"fmt"
	msgE "github.com/prestonTao/messageEngine"
	"time"
)

func main() {
	example1()
}

func example1() {
	engine := msgE.NewEngine("interClient")
	engine.RegisterMsg(111, RecvMsg)
	engine.AddClientConn("test", "127.0.0.1", 9090)

	//给服务器发送消息
	session, _ := engine.GetController().GetSession("test")
	hello := []byte("hello, I'm client")
	session.Send(111, &hello)
	time.Sleep(time.Second * 10)

}

func RecvMsg(c msgE.Controller, msg msgE.GetPacket) {
	fmt.Println(string(msg.Date))
}

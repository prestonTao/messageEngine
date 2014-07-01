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
	engine := msgE.NewEngine("interServer")
	engine.RegisterMsg(111, RecvMsg)
	engine.Listen("127.0.0.1", 9090)
	time.Sleep(time.Second * 10)

}

func RecvMsg(c msgE.Controller, msg msgE.GetPacket) {
	fmt.Println(string(msg.Date))
	session, ok := c.GetSession(msg.Name)
	if ok {
		hello := []byte("hello, I'm server")
		session.Send(111, &hello)
		session.Close()
	}
}

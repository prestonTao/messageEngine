package main

import (
	"fmt"
	msgE "github.com/prestonTao/messageEngine"
	"time"
)

func main() {
	simple1()
}

//启动服务器10秒钟
func simple1() {
	engine := msgE.NewEngine("interServer")
	engine.RegisterMsg(111, RecvMsg)
	engine.Listen("127.0.0.1", 9090)
	time.Sleep(time.Second * 10)
}

func RecvMsg(c msgE.Controller, msg msgE.GetPacket) {
	fmt.Println(string(msg.Date))
	session, ok := c.GetSession(msg.Name)
	if ok {
		session.Close()
	}

}
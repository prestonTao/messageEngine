package main

import (
	// "fmt"
	msgE "github.com/prestonTao/messageEngine"
	"time"
)

func main() {
	simple1()
}

//启动服务器10秒钟
func simple1() {

	engine := msgE.NewEngine("interClient")
	engine.AddClientConn("test", "127.0.0.1", 9090)
	session, _ := engine.GetController().GetSession("test")
	hello := []byte("hello")
	session.Send(1, &hello)

	time.Sleep(time.Second * 10)

}

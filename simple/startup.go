package main

import (
	msgE "github.com/prestonTao/messageEngine"
	"time"
)

func main() {

}

//启动服务器10秒钟
func simple1() {
	msgE.IP = "127.0.0.1"
	msgE.PORT = 9090
	server := new(msgE.ServerManager)
	server.Run()
	time.Sleep(time.Second * 10)
}

func example2() {
	msgE.IP = "127.0.0.1"
	msgE.PORT = 9090
	server := new(msgE.ServerManager)
	server.Run()
	time.Sleep(time.Second * 10)
}

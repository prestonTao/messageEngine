package main

import (
	msgE "github.com/prestonTao/messageEngine"
	"time"
)

func main() {
	msgE.IP = "127.0.0.1"
	msgE.PORT = 9090
	server := new(msgE.ServerManager)
	server.Run()
	time.Sleep(time.Second * 10)
}

package messageEngine

import (
	// "mandela/peerNode/messageEngine/net"
	"sync"
)

type MsgHandler func(c Controller, msg GetPacket)

var msgHandlerStore = make(map[int32]MsgHandler, 1000)

var routerLock sync.RWMutex

func init() {
	// connService := ConnService{}
	// AddRouter(1, connService.Connect)
	// AddRouter(2, connService.CloseConnect)
}

func AddRouter(msgId int32, handler MsgHandler) {
	routerLock.Lock()
	defer routerLock.Unlock()
	msgHandlerStore[msgId] = handler
}

func GetHandler(msgId int32) MsgHandler {
	routerLock.Lock()
	defer routerLock.Unlock()

	handler := msgHandlerStore[msgId]
	return handler
}

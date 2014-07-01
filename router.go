package messageEngine

import (
	// "mandela/peerNode/messageEngine/net"
	"sync"
)

type MsgHandler func(c Controller, msg GetPacket)

type RouterStore struct {
	lock     *sync.RWMutex
	handlers map[int32]MsgHandler
}

func (this *RouterStore) AddRouter(msgId int32, handler MsgHandler) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.handlers[msgId] = handler
}

func (this *RouterStore) GetHandler(msgId int32) MsgHandler {
	this.lock.Lock()
	defer this.lock.Unlock()

	handler := this.handlers[msgId]
	return handler
}

func NewRouter() *RouterStore {
	router := new(RouterStore)
	router.lock = new(sync.RWMutex)
	router.handlers = make(map[int32]MsgHandler)
	return router
}

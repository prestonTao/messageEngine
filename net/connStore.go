package net

import (
	"errors"
	"fmt"
	"sync"
)

type Conn interface {
	Send(msgID uint32, data *[]byte)
	Close()
}

//根据sessionId保存连接
var connStore = make(map[uint64]Conn, 10000)

var chs_lock sync.RWMutex

//注册Conn
func addConn(id uint64, client Conn) {
	chs_lock.Lock()
	defer chs_lock.Unlock()
	connStore[id] = client
	fmt.Println("添加一个conn", client, id, connStore)
}

//移除某个Conn
func removeConn(id uint64) {
	chs_lock.Lock()
	defer chs_lock.Unlock()
	delete(connStore, id)
}

//获取某个Conn
func getConn(id uint64) (client Conn, err error) {
	chs_lock.RLock()
	defer chs_lock.RUnlock()
	client, ok := connStore[id]
	if !ok {
		err = errors.New(fmt.Sprintf("uid %x 没有对应的 Session", id))
		return nil, err
	}
	return client, nil
}

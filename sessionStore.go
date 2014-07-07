package messageEngine

import (
	// "errors"
	// "fmt"
	"sync"
)

type sessionBase struct {
	name      string
	attrbutes map[string]interface{}
}

func (this *sessionBase) Set(name string, value interface{}) {
	this.attrbutes[name] = value
}
func (this *sessionBase) Get(name string) interface{} {
	return this.attrbutes[name]
}
func (this *sessionBase) GetName() string {
	return this.name
}
func (this *sessionBase) SetName(name string) {
	this.name = name
}
func (this *sessionBase) Send(msgID uint32, data *[]byte) (err error) { return }
func (this *sessionBase) Close()                                      {}

type Session interface {
	Send(msgID uint32, data *[]byte) error
	Close()
	Set(name string, value interface{})
	Get(name string) interface{}
	GetName() string
	SetName(name string)
}

type sessionStore struct {
	lock *sync.RWMutex
	// store     map[int64]Session
	nameStore map[string]Session
}

func (this *sessionStore) addSession(name string, session Session) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.nameStore[session.GetName()] = session
	// sessionStore.store[sessionId] = session
}

func (this *sessionStore) getSession(name string) (Session, bool) {
	this.lock.Lock()
	defer this.lock.Unlock()
	s, ok := this.nameStore[name]
	return s, ok
}

func (this *sessionStore) removeSession(name string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	delete(this.nameStore, name)
}

func NewSessionStore() *sessionStore {
	sessionStore := new(sessionStore)
	sessionStore.lock = new(sync.RWMutex)
	sessionStore.nameStore = make(map[string]Session)
	return sessionStore
}

// var sessionStore = new(sessionStore)

// func init() {
// 	sessionStore.lock = new(sync.RWMutex)
// 	sessionStore.nameStore = make(map[string]Session, 10000)
// }

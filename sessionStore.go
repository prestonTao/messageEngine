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
func (this *sessionBase) Send(msgID uint32, data *[]byte) {}
func (this *sessionBase) Close()                          {}

type Session interface {
	Send(msgID uint32, data *[]byte)
	Close()
	Set(name string, value interface{})
	Get(name string) interface{}
	GetName() string
	SetName(name string)
}

type sessionProvider struct {
	lock *sync.RWMutex
	// store     map[int64]Session
	nameStore map[string]Session
}

func addSession(name string, session Session) {
	sessionStore.lock.Lock()
	defer sessionStore.lock.Unlock()
	sessionStore.nameStore[session.GetName()] = session
	// sessionStore.store[sessionId] = session
}

func getSession(name string) (Session, bool) {
	sessionStore.lock.Lock()
	defer sessionStore.lock.Unlock()
	s, ok := sessionStore.nameStore[name]
	return s, ok
}

func removeSession(name string) {
	sessionStore.lock.Lock()
	defer sessionStore.lock.Unlock()
	delete(sessionStore.nameStore, name)
}

var sessionStore = new(sessionProvider)

func init() {
	sessionStore.lock = new(sync.RWMutex)
	sessionStore.nameStore = make(map[string]Session, 10000)
}

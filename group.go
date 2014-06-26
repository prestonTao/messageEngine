package messageEngine

import (
	"sync"
)

type groupOne struct {
	lock  *sync.RWMutex
	group map[string]Session
}

//
type msgGroupProvider struct {
	lock   *sync.RWMutex
	groups map[string]*groupOne
}

//添加一个小组
func addGroup(name string) {
	group.lock.Lock()
	defer group.lock.Unlock()
	group.groups[name] = make(map[string]Session)
}

func removeGroup(name string) {
	group.lock.Lock()
	defer group.lock.Unlock()
	delete(group.groups, name)
}
func addToGroup(groupName, name string, session Session) {

}

var group = new(msgGroupProvider)

func init() {
	group.lock = new(sync.RWMutex)
	group.groups = make(map[string]*groupOne)
}

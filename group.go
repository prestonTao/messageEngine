package messageEngine

import (
	"sync"
)

type groupOne struct {
	lock     *sync.RWMutex
	msgGroup map[string]Session
}

//
type msgGroupProvider struct {
	lock   *sync.RWMutex
	groups map[string]*groupOne
}

//创建一个小组
func createGroup(groupName string) {
	msgGroup.lock.Lock()
	defer msgGroup.lock.Unlock()
	msgGroup.groups[groupName] = new(groupOne)
}

//删除一个小组
func removeGroup(groupName string) {
	msgGroup.lock.Lock()
	defer msgGroup.lock.Unlock()
	delete(msgGroup.groups, groupName)
}

//将一个连接添加到组中
func addToGroup(groupName, name string) {
	groupTag, ok := msgGroup.groups[groupName]
	if !ok {
		createGroup(groupName)
		groupTag, _ = msgGroup.groups[groupName]
	}
	session, ok := getSession(name)
	if ok {
		groupTag.msgGroup[name] = session
	}
}

//检查一个name是否在某个组中
func checkNameInGroup(groupName, name string) bool {
	groupTag, ok := msgGroup.groups[groupName]
	if ok {
		if _, ok = groupTag.msgGroup[name]; ok {
			return true
		}
	}
	return false
}

var msgGroup = new(msgGroupProvider)

func init() {
	msgGroup.lock = new(sync.RWMutex)
	msgGroup.groups = make(map[string]*groupOne)
}

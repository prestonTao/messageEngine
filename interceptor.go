package messageEngine

import (
	"sync"
)

// type interceptor struct{}

// func (this *interceptor) in() {

// }

type Interceptor interface {
	In(c Controller, msg GetPacket) bool
	Out(c Controller, msg GetPacket)
}

type InterceptorProvider struct {
	lock         *sync.RWMutex
	interceptors []Interceptor
}

func addInterceptor(itpr Interceptor) {
	interceptors.lock.Lock()
	defer interceptors.lock.Unlock()
	interceptors.interceptors = append(interceptors.interceptors, itpr)
}
func getInterceptors() []Interceptor {
	interceptors.lock.Lock()
	defer interceptors.lock.Unlock()
	return interceptors.interceptors

}

var interceptors *InterceptorProvider

func init() {
	interceptors = new(InterceptorProvider)
	interceptors.lock = new(sync.RWMutex)
	interceptors.interceptors = make([]Interceptor, 0)
}

// var chanS =

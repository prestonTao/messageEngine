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

func (this *InterceptorProvider) addInterceptor(itpr Interceptor) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.interceptors = append(this.interceptors, itpr)
}
func (this *InterceptorProvider) getInterceptors() []Interceptor {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.interceptors

}

func NewInterceptor() *InterceptorProvider {
	interceptor := new(InterceptorProvider)
	interceptor.lock = new(sync.RWMutex)
	interceptor.interceptors = make([]Interceptor, 0)
	return interceptor
}

// var interceptors *InterceptorProvider

// func init() {
// 	interceptors = new(InterceptorProvider)
// 	interceptors.lock = new(sync.RWMutex)
// 	interceptors.interceptors = make([]Interceptor, 0)
// }

// var chanS =

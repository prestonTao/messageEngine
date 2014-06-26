package messageEngine

// type interceptor struct{}

// func (this *interceptor) in() {

// }

type Interceptor interface {
	In(c Controller, msg GetPacket) bool
	Out(c Controller, msg GetPacket)
}

func addInterceptor(itpr Interceptor) {
	interceptors = append(interceptors, itpr)
}
func getInterceptor() {

}

var interceptors = make([]Interceptor, 0)
var chanS = 
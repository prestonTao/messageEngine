package messageEngine

import (
	// "mandela/peerNode/messageEngine/net"
	// "bufio"
	"fmt"
	// "os"
	// "strconv"
	"sync"
)

// var (
// 	Name       = "myserver"
// 	IP         = "127.0.0.1"
// 	PORT int32 = 9090
// )

type Engine struct {
	name        string
	status      int //服务器状态
	net         *Net
	receive     <-chan *GetPacket
	controller  Controller
	auth        Auth
	interceptor *InterceptorProvider
	onceRead    *sync.Once
}

func (this *Engine) Listen(ip string, port int32) {
	this.run()
	this.net.Listen(ip, port)
}

//添加一个连接，给这个连接取一个名字，连接名字可以在自定义权限验证方法里面修改
func (this *Engine) AddClientConn(name, ip string, port int32) {
	this.run()
	this.net.AddClientConn(name, ip, port)
}

//添加一个拦截器，所有消息到达业务方法之前都要经过拦截器处理
func (this *Engine) AddInterceptor(itpr Interceptor) {
	this.interceptor.addInterceptor(itpr)
}

//得到控制器
func (this *Engine) GetController() Controller {
	return this.controller
}

func (this *Engine) SetAuth(auth Auth) {
	if auth == nil {
		return
	}
	defaultAuth = auth
}

func (this *Engine) run() {
	//保证方法只执行一次
	go this.onceRead.Do(func() {
		this.receive = this.net.Recv
		//构建控制器
		this.buildController()
		go this.read()
	})
}

func (this *Engine) buildController() {
	c := new(ControllerImpl)
	c.net = this.net
	c.engine = this
	c.attributes = make(map[string]interface{})
	c.msgGroup = NewMsgGroupManager()
	c.msgGroup.controller = c
	this.controller = c
}

//读取网络模块发送来的消息
func (this *Engine) read() {
	//保证将消息处理完才关闭服务器
	for msg := range this.receive {
		fmt.Println("haha")
		this.handler(msg)
	}
}

//负责将接收到的消息转换为结构体
func (this *Engine) handler(msg *GetPacket) {
	handler := GetHandler(msg.MsgID)
	if handler == nil {
		fmt.Println("该消息未注册，消息编号：", msg.MsgID)
		return
	}
	//这里决定了消息是否异步处理
	go this.handlerProcess(handler, msg)
}

func (this *Engine) handlerProcess(handler MsgHandler, msg *GetPacket) {
	defer func() {
		if err := recover(); err != nil {
			e, ok := err.(error)
			if ok {
				fmt.Println(e.Error())
			}
		}
	}()
	//消息处理前先通过拦截器
	itps := this.interceptor.getInterceptors()
	itpsLen := len(itps)
	for i := 0; i < itpsLen; i++ {
		itps[i].In(this.controller, *msg)
	}
	handler(this.controller, *msg)
	//消息处理后也要通过拦截器
	for i := itpsLen - 1; i >= 0; i-- {
		itps[i].Out(this.controller, *msg)
	}
}

func NewEngine(name string) *Engine {
	engine := new(Engine)
	engine.name = name
	engine.interceptor = NewInterceptor()
	engine.onceRead = new(sync.Once)
	engine.net = NewNet()
	return engine
}

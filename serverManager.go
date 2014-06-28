package messageEngine

import (
	// "mandela/peerNode/messageEngine/net"
	// "bufio"
	"fmt"
	// "os"
	// "strconv"
)

var (
	Name       = "myserver"
	IP         = "127.0.0.1"
	PORT int32 = 9090
)

type ServerManager struct {
	isRun      bool //服务器是否开启
	net        *Net
	receive    <-chan *GetPacket
	controller Controller
	auth       Auth
}

func (this *ServerManager) Run() {
	if !this.isRun {
		this.isRun = true
		//启动网络
		this.net = NewNet(IP, int32(PORT), this.auth)
		this.receive = this.net.Recv
		go this.read()
		//构建控制器
		this.buildController()
		//运行启动钩子
		// runHookFunc(this.controller)
	}
}

func (this *ServerManager) buildController() {
	c := new(ControllerImpl)
	c.net = this.net
	c.serverManager = this
	c.attributes = make(map[string]interface{})
	c.msgGroup = NewMsgGroupManager()
	this.controller = c
}

//添加一个连接，给这个连接取一个名字
func (this *ServerManager) AddClientConn(name, ip string, port int32) {
	this.Run()
	this.net.AddClientConn(name, ip, port)
	// addAcc(name, client.Session)
}

//添加一个拦截器，所有消息到达业务方法之前都要经过拦截器处理
func (this *ServerManager) AddInterceptor(itpr Interceptor) {
	addInterceptor(itpr)
}

//得到控制器
func (this *ServerManager) GetController() Controller {
	return this.controller
}

//读取网络模块发送来的消息
func (this *ServerManager) read() {
	//保证将消息处理完才关闭服务器
	for msg := range this.receive {
		this.handler(msg)
	}

}

//负责将接收到的消息转换为结构体
func (this *ServerManager) handler(msg *GetPacket) {
	handler := GetHandler(msg.MsgID)
	if handler == nil {
		fmt.Println("该消息未注册，消息编号：", msg.MsgID)
		return
	}
	//这里决定了消息是否异步处理
	this.handlerProcess(handler, msg)
}

func (this *ServerManager) handlerProcess(handler MsgHandler, msg *GetPacket) {
	defer func() {
		if err := recover(); err != nil {
			e, ok := err.(error)
			if ok {
				fmt.Println(e.Error())
			}
		}
	}()
	itps := getInterceptors()
	itpsLen := len(itps)
	for i := 0; i < itpsLen; i++ {
		itps[i].In(this.controller, *msg)
	}
	handler(this.controller, *msg)
	for i := itpsLen - 1; i >= 0; i-- {
		itps[i].Out(this.controller, *msg)
	}
}

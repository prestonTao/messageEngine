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
	isRun       bool //服务器是否开启
	net         *Net
	receive     <-chan *GetPacket
	controller  Controller
	auth        Auth
	interceptor *interceptor
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
	this.controller = c
}

func (this *ServerManager) AddClientConn(name, ip string, port int32) {
	this.Run()
	this.net.AddClientConn(name, ip, port)
	// addAcc(name, client.Session)
}

func (this *ServerManager) AddInterceptor(itpr Interceptor) {

}

//得到控制器
func (this *ServerManager) GetController() Controller {
	return this.controller
}

//读取网络模块发送来的消息
func (this *ServerManager) read() {
	// for this.isRun {
	// 	msg := <-this.receive
	// 	this.handler(msg)
	// }
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
	handler(this.controller, *msg)
}

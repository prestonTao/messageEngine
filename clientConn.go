package messageEngine

import (
	// "encoding/binary"
	"fmt"
	// "hash/crc32"
	// "io"
	"net"
	"strconv"
	"time"
)

type CloseCallback func(name string)

//本机向其他服务器的连接
type Client struct {
	sessionBase
	serverName string
	ip         string
	port       int32
	conn       net.Conn
	// outData    chan *[]byte    //发送队列
	inPack     chan *GetPacket //接收队列
	isClose    bool            //该连接是否被关闭
	isPowerful bool            //是否是强连接，强连接有短线重连功能
	net        *Net
	// call       CloseCallback
}

func (this *Client) Connect(ip string, port int32) (remoteName string, err error) {

	this.ip = ip
	this.port = port

	this.conn, err = net.Dial("tcp", ip+":"+strconv.Itoa(int(port)))
	if err != nil {
		return
	}

	//权限验证
	remoteName, err = defaultAuth.SendKey(this.conn, this, this.serverName)
	if err != nil {
		return
	}

	fmt.Println("Connecting to", ip, ":", strconv.Itoa(int(port)))

	go this.recv()
	// go this.send()
	// go this.hold()
	return
}
func (this *Client) reConnect() {
	for {
		//十秒钟后重新连接
		time.Sleep(time.Second * 10)
		var err error
		this.conn, err = net.Dial("tcp", this.ip+":"+strconv.Itoa(int(this.port)))
		if err != nil {
			continue
		}

		fmt.Println("Connecting to", this.ip, ":", strconv.Itoa(int(this.port)))

		go this.recv()
		// go this.send()
		// go this.hold()
		return
	}
}

func (this *Client) recv() {
	for !this.isClose {
		packet, err, isClose := RecvPackage(this.conn)

		if isClose {
			this.isClose = true
			break
		}
		if err == nil {
			packet.Name = this.GetName()
			this.inPack <- packet
			continue
		}
		fmt.Println("接收数据出错  ", err.Error())
	}
	// fmt.Println(this.call, this.isPowerful)
	// if this.call != nil {
	// 	this.call(this.GetName())
	// }

	this.net.CloseClient(this.GetName())
	if this.isPowerful {
		go this.reConnect()
	}
	//最后一个包接收了之后关闭chan
	//如果有超时包需要等超时了才关闭，目前未做处理
	// close(this.outData)
	// fmt.Println("recv 协成走完")
}

// func (this *Client) send() {
// 	defer func() {
// 		// close(this.outData)
// 		this.isClose = true
// 		// fmt.Println("send 协成走完")
// 	}()
// 	// //处理客户端主动断开连接的情况
// 	for msg := range this.outData {
// 		if _, err := this.conn.Write(*msg); err != nil {
// 			log.Println("发送数据出错", err)
// 			return
// 		}
// 	}

// }

//心跳连接
// func (this *Client) hold() {
// 	for !this.isClose {
// 		// fmt.Println("hold")
// 		time.Sleep(time.Second * 2)
// 		bs := []byte("")
// 		this.Send(0, &bs)
// 	}
// 	// close(this.outData)
// 	this.net.CloseClient(this.GetName())
// 	fmt.Println("hold 协成走完")
// }

//发送序列化后的数据
func (this *Client) Send(msgID uint32, data *[]byte) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err, _ = e.(error)
			fmt.Println("发送序列化的数据出错  ", err.Error())
		}
	}()
	buff := MarshalPacket(msgID, data)
	// this.outData <- buff
	_, err = this.conn.Write(*buff)
	// if _, err = this.conn.Write(*msg); err != nil {
	// 	log.Println("发送数据出错", err)
	// 	return
	// }
	return
}

// func (this *Client) GetOneMsg() {

// }

// //发送
// func (this *Client) SendBytes(msgID uint32, data []byte) {
// 	buff := MarshalPacket(msgID, &data)
// 	this.outData <- buff
// }

//客户端关闭时,退出recv,send
func (this *Client) Close() {
	this.isClose = true
}
func NewClient(name, ip string, port int32) *Client {
	client := new(Client)
	client.name = name
	client.inPack = make(chan *GetPacket, 1000)
	// client.outData = make(chan *[]byte, 1000)
	client.Connect(ip, port)
	return client
}

package messageEngine

import (
	"fmt"
	"net"
	// "time"
)

//其他计算机对本机的连接
type ServerConn struct {
	sessionBase
	conn           net.Conn
	Ip             string
	Connected_time string
	CloseTime      string
	// outData        chan *[]byte //序列化后的GetPacket
	inPack  chan *GetPacket
	isClose bool //该连接是否已经关闭
	net     *Net
}

func (this *ServerConn) run() {
	go this.recv()
	// go this.send()
	// go this.hold()
}

//接收客户端消息协程
func (this *ServerConn) recv() {
	//处理客户端主动断开连接的情况
	for !this.isClose {
		packet, err, isClose := RecvPackage(this.conn)
		if isClose {
			this.isClose = true
			break
		}
		if err == nil {
			if packet.MsgID == 0 {
				//hold 心跳包
				continue
			}
			packet.Name = this.GetName()
			this.inPack <- packet
			continue
		}
		fmt.Println("接收数据出错  ", err.Error())
	}

	this.net.CloseClient(this.GetName())
	//最后一个包接收了之后关闭chan
	//如果有超时包需要等超时了才关闭，目前未做处理
	// close(this.outData)
	// fmt.Println("关闭连接")
}

//发送给客户端消息协程
// func (this *ServerConn) send() {
// 	//处理客户端主动断开连接的情况
// 	//确保消息发送完后再关闭连接
// 	for msg := range this.outData {
// 		if _, err := this.conn.Write(*msg); err != nil {
// 			log.Println("发送数据出错", err)
// 			return
// 		}
// 	}
// }

// //心跳连接
// func (this *ServerConn) hold() {
// 	for !this.isClose {
// 		time.Sleep(time.Second * 5)
// 		bs := []byte("")
// 		this.Send(0, &bs)
// 	}
// }

//给客户端发送数据
func (this *ServerConn) Send(msgID uint32, data *[]byte) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err, _ = e.(error)
		}
	}()
	buff := MarshalPacket(msgID, data)
	// this.outData <- buff
	_, err = this.conn.Write(*buff)
	return
}

//关闭这个连接
func (this *ServerConn) Close() {
	// fmt.Println("调用关闭连接方法")
	this.isClose = true
}

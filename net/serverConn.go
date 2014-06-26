package net

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"net"
)

//其他计算机对本机的连接
type ServerConn struct {
	conn           net.Conn
	session        uint64
	Ip             string
	Connected_time string
	CloseTime      string
	outData        chan *[]byte //序列化后的GetPacket
	inPack         chan *GetPacket
	isClose        bool //该连接是否已经关闭
}

func (this *ServerConn) run() {
	go this.recv()
	go this.send()
}

//接收客户端消息协程
func (this *ServerConn) recv() {
	//处理客户端主动断开连接的情况
	defer func(clientConn *ServerConn) {
		log.Println("客户端主动关闭连接")
	}(this)

	for !this.isClose {
		var data []byte
		header := make([]byte, 12)

		n, err := io.ReadFull(this.conn, header)
		if n == 0 && err == io.EOF {
			temp := new(GetPacket)
			temp.ConnId = this.session
			temp.MsgID = int32(Close)
			this.inPack <- temp
			fmt.Println("客户端断开连接")
			return
		} else if err != nil {
			fmt.Println("接收数据出错:", err)
			return
		}
		//数据包长度

		size := binary.BigEndian.Uint32(header)
		//crc值
		crc1 := binary.BigEndian.Uint32(header[4:8])
		msgID := binary.BigEndian.Uint32(header[8:12])

		data = make([]byte, size)
		n, err = io.ReadFull(this.conn, data)
		if uint32(n) != size {
			log.Println("数据包长度不正确", n, "!=", size)
			return
		}
		if err != nil {
			log.Println("读取数据出错:", err)
			return
		}

		crc2 := crc32.Checksum(data, crc32.IEEETable)
		if crc1 != crc2 {
			log.Println("crc 数据验证不正确: ", crc1, " != ", crc2)
			return
		}

		temp := new(GetPacket)
		temp.ConnId = this.session
		temp.Date = data
		temp.MsgID = int32(msgID)
		temp.Size = uint32(len(data))
		this.inPack <- temp
	}
	//最后一个包接收了之后关闭chan
	//如果有超时包需要等超时了才关闭，目前未做处理
	close(this.outData)
}

//发送给客户端消息协程
func (this *ServerConn) send() {
	// //处理客户端主动断开连接的情况
	// defer func(clientConn *ServerConn) {
	// 	log.Println("关闭此连接发送协程")
	// 	this.conn.Close()
	// 	close(this.outData)
	// }(this)
	// for !this.isClose {
	// 	msg := <-this.outData
	// 	if _, err := this.conn.Write(*msg); err != nil {
	// 		log.Println("发送数据出错", err)
	// 		return
	// 	}
	// }
	//确保消息发送完后再关闭连接
	for msg := range this.outData {
		if _, err := this.conn.Write(*msg); err != nil {
			log.Println("发送数据出错", err)
			return
		}
	}
}

//给客户端发送数据
func (this *ServerConn) Send(msgID uint32, data *[]byte) {
	buff := PacketData(msgID, data)
	this.outData <- buff
}

//关闭这个连接
func (this *ServerConn) Close() {
	this.isClose = true

}

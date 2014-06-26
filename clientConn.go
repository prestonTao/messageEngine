package messageEngine

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

//本机向其他服务器的连接
type Client struct {
	sessionBase
	// Session uint64
	ip      string
	port    int32
	conn    net.Conn
	outData chan *[]byte    //发送队列
	inPack  chan *GetPacket //接收队列
	isClose bool            //该连接是否被关闭
}

func (this *Client) Connect(ip string, port int32) error {

	this.ip = ip
	this.port = port

	var err error
	this.conn, err = net.Dial("tcp", ip+":"+strconv.Itoa(int(port)))
	if err != nil {
		fmt.Println(err)
		return err
	}

	//权限验证
	err = defaultAuth.SendKey(this.conn, this)
	if err != nil {
		return err

	}
	addSession(this.name, this)

	fmt.Println("Connecting to", ip)

	// temp := new(GetPacket)
	// temp.ConnId = this.Session
	// temp.MsgID = int32(Connect)
	// this.inPack <- temp

	go this.recv()
	go this.send()
	return err
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

		fmt.Println("Connecting to", this.ip)

		// temp := new(GetPacket)
		// temp.ConnId = this.Session
		// temp.MsgID = int32(Connect)
		// this.inPack <- temp

		go this.recv()
		go this.send()
		return
	}
}

func (this *Client) recv() {
	for !this.isClose {

		header := make([]byte, 12)

		n, err := io.ReadFull(this.conn, header)

		if n == 0 && err == io.EOF {
			fmt.Println("客户端断开连接")
			go this.reConnect()
			return
		} else if err != nil {
			fmt.Println("接受数据出错:", err)
			go this.reConnect()
			return
		}

		//数据包长度
		size := binary.BigEndian.Uint32(header)
		//crc值
		crc1 := binary.BigEndian.Uint32(header[4:8])

		// msgID := binary.BigEndian.Uint32(header[8:12])

		data := make([]byte, size)

		n, err = io.ReadFull(this.conn, data)
		if err != nil {
			log.Println("读取数据出错:", err)
			go this.reConnect()
			return
		}
		if uint32(n) != size {
			log.Println("数据包长度不正确", n, "!=", size)
			continue
		}

		crc2 := crc32.Checksum(data, crc32.IEEETable)

		if crc1 != crc2 {
			log.Println("crc 数据验证不正确: ", crc1, " != ", crc2)
			continue
		}

		// temp := new(GetPacket)
		// temp.ConnId = this.Session
		// temp.Date = data
		// temp.MsgID = int32(msgID)
		// temp.Size = uint32(len(data))
		// this.inPack <- temp
	}
	//最后一个包接收了之后关闭chan
	//如果有超时包需要等超时了才关闭，目前未做处理
	close(this.outData)
}

func (this *Client) send() {
	// //处理客户端主动断开连接的情况
	// defer func(clientConn *Client) {
	// 	log.Println("关闭此连接发送协程")
	// 	this.conn.Close()
	// 	close(this.outData)
	// }(this)
	// for {
	// 	select {
	// 	case msg := <-this.outData:
	// 		if _, err := this.conn.Write(*msg); err != nil {
	// 			log.Println("发送数据出错", err)
	// 			return
	// 		}
	// 	default:
	// 		if this.isClose {
	// 			return
	// 		}
	// 	}
	// }

	for msg := range this.outData {
		if _, err := this.conn.Write(*msg); err != nil {
			log.Println("发送数据出错", err)
			return
		}
	}
}

//发送序列化后的数据
func (this *Client) Send(msgID uint32, data *[]byte) {
	buff := PacketData(msgID, data)
	this.outData <- buff
}

func (this *Client) GetOneMsg() {

}

//发送
func (this *Client) SendBytes(msgID uint32, data []byte) {
	buff := PacketData(msgID, &data)
	this.outData <- buff
}

//客户端关闭时,退出recv,send
func (this *Client) Close() {
	this.isClose = true
}
func NewClient(name, ip string, port int32) *Client {
	client := new(Client)
	client.name = name
	client.inPack = make(chan *GetPacket, 1000)
	client.outData = make(chan *[]byte, 1000)
	client.Connect(ip, port)
	return client
}

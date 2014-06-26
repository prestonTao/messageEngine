package net

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

type Net struct {
	lock    *sync.RWMutex
	session uint64
	Recv    chan *GetPacket //获得数据
}

func (this *Net) start(ip string, port int32) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ip+":"+strconv.Itoa(int(port)))
	if err != nil {
		fmt.Println(err)
		return
	}

	listener, err := net.ListenTCP("tcp4", tcpAddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(ip + ":" + strconv.Itoa(int(port)) + "成功启动服务器")

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go this.newConnect(conn)
	}
}

//创建一个新的连接
func (this *Net) newConnect(conn net.Conn) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.session++

	client := new(ServerConn)
	client.session = this.session
	client.conn = conn
	client.Ip = conn.RemoteAddr().String()
	client.Connected_time = time.Now().String()
	client.outData = make(chan *[]byte, 1000)
	client.inPack = this.Recv
	client.run()

	addConn(this.session, client)

	fmt.Println(time.Now().String(), "建立连接", conn.RemoteAddr().String())

	temp := new(GetPacket)
	temp.ConnId = this.session
	temp.MsgID = int32(Connect)

	this.Recv <- temp

}

//关闭连接
func (this *Net) CloseClient(id uint64) error {
	client, err := getConn(id)
	if err != nil {
		return err
	}
	client.Close()
	removeConn(id)
	return nil
}

func (this *Net) AddClientConn(ip string, port int32) *Client {
	this.lock.Lock()
	defer this.lock.Unlock()
	//-------------------
	//保证把原有的队列里的数据取出才能替换
	//-------------------
	this.session++

	client := new(Client)
	client.Session = this.session
	client.inPack = this.Recv
	client.outData = make(chan *[]byte, 2000)
	client.Connect(ip, port)
	addConn(this.session, client)

	return client
}

func (this *Net) GetConn(session uint64) Conn {
	client, err := getConn(session)
	if err != nil {
		return nil
	}
	return client
}

//发送数据
func (this *Net) Send(connId uint64, msgID uint32, data []byte) error {

	client, err := getConn(connId)
	if err != nil {
		return err
	}
	client.Send(msgID, &data)

	return nil
}

func NewNet(ip string, port int32) *Net {

	net := new(Net)
	net.Recv = make(chan *GetPacket, 5000)
	net.lock = new(sync.RWMutex)
	go net.start(ip, port)
	time.Sleep(time.Millisecond * 3)
	return net
}

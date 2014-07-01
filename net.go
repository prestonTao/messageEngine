package messageEngine

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

type Net struct {
	Recv         chan *GetPacket //获得数据
	sessionStore *sessionStore
}

func (this *Net) Listen(ip string, port int32) {
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
	fmt.Println("监听一个地址：", ip+":"+strconv.Itoa(int(port)))
	// fmt.Println(ip + ":" + strconv.Itoa(int(port)) + "成功启动服务器")
	go this.listener(listener)
}

func (this *Net) listener(listener net.Listener) {
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
	name, err := defaultAuth.RecvKey(conn)
	if err != nil {
		return
	}

	// this.lock.Lock()
	// defer this.lock.Unlock()

	// this.session++

	serverConn := &ServerConn{
		// attrbutes:      make(map[string]interface{}),
		conn: conn,
		// session:        this.session,
		Ip:             conn.RemoteAddr().String(),
		Connected_time: time.Now().String(),
		outData:        make(chan *[]byte, 1000),
		inPack:         this.Recv,
	}
	serverConn.name = name
	serverConn.attrbutes = make(map[string]interface{})
	serverConn.run()

	// client := new(ServerConn)
	// client.session = this.session
	// client.conn = conn
	// client.Ip = conn.RemoteAddr().String()
	// client.Connected_time = time.Now().String()
	// client.outData = make(chan *[]byte, 1000)
	// client.inPack = this.Recv
	// client.run()
	this.sessionStore.addSession(name, serverConn)
	// addConn(this.session, serverConn)

	fmt.Println(time.Now().String(), "建立连接", conn.RemoteAddr().String())

	// temp := new(GetPacket)
	// temp.ConnId = this.session
	// temp.MsgID = int32(Connect)

	// this.Recv <- temp

}

//关闭连接
func (this *Net) CloseClient(name string) bool {
	session, ok := this.sessionStore.getSession(name)
	if ok {
		this.sessionStore.removeSession(name)
		session.Close()
		return true
	}
	return false
}

func (this *Net) AddClientConn(name, ip, serverName string, port int32) (*Client, error) {
	// this.lock.Lock()
	// defer this.lock.Unlock()
	//-------------------
	//保证把原有的队列里的数据取出才能替换
	//-------------------
	// this.session++

	clientConn := &Client{
		serverName: serverName,
		outData:    make(chan *[]byte, 2000),
		inPack:     this.Recv,
	}
	clientConn.name = name
	clientConn.attrbutes = make(map[string]interface{})
	err := clientConn.Connect(ip, port)
	if err == nil {
		this.sessionStore.addSession(name, clientConn)
		return clientConn, nil
	}
	return nil, err
}

func (this *Net) GetSession(name string) (Session, bool) {
	return this.sessionStore.getSession(name)
}

//发送数据
func (this *Net) Send(name string, msgID uint32, data []byte) bool {
	session, ok := this.sessionStore.getSession(name)
	if ok {
		session.Send(msgID, &data)
		return true
	} else {
		return false
	}
}

func NewNet() *Net {

	net := new(Net)
	net.sessionStore = NewSessionStore()
	net.Recv = make(chan *GetPacket, 5000)
	// net.start(ip, port)
	// time.Sleep(time.Millisecond * 3)
	return net
}

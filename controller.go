package messageEngine

import (
// "fmt"
// "mandela/peerNode/messageEngine/net"
// "common/message"
// "github.com/ziutek/mymysql/mysql"
// "sync"
)

type Controller interface {
	// Register(name string, sessionId uint64) //通过账号注册一个客户端连接
	GetSession(name string) (Session, bool) //通过accId得到客户端的连接Id
	// GetConfig() *ConfigFile                      //得到配置文件
	// GetDBConn(dbNumber int32) mysql.Conn         //得到数据库连接
	GetNet() *Net                                //获得连接到本地的计算机连接
	SetAttribute(name string, value interface{}) //设置共享数据，实现业务模块之间通信
	GetAttribute(name string) interface{}        //得到共享数据，实现业务模块之间通信
	// GetClientById(session uint64) Conn           //得到本机到其他计算机的连接
	// GetDatabaseId(accName string) (int, int)     //获得数据库分库分表的id
	// Forward()                                    //服务器跳转

}

type ControllerImpl struct {
	// lock       sync.RWMutex
	// Config     *ConfigFile
	// ConnStore  map[int32]mysql.Conn
	net           *Net
	serverManager *ServerManager
	attributes    map[string]interface{}
	// Client     *Client
}

// func (this *ControllerImpl) Register(name string, sessionId uint64) {
// 	addAcc(name, sessionId)
// }

// func (this *ControllerImpl) GetClientByName(name string) Conn {
// 	sessionId, err := getSessions(name)
// 	if err != nil {
// 		return nil
// 	}
// 	return this.GetClientById(sessionId)
// }

//连接数据库
// func (this *ControllerImpl) connDB() {
// 	conn := DBConn{}
// 	conn.Connect(this.Config)
// 	this.ConnStore = conn.ConnStore
// }

//作为客户端连接其他服务器
// func (this *ControllerImpl) connClient() {
// 	c := &Client{
// 		net: this.net,
// 	}
// 	c.invoke(this.Config)
// 	this.Client = c
// }

//添加csv文件配置
// func (this *ControllerImpl) deploy() {

// }

//构建控制器
// func (this *ControllerImpl) build(config *ConfigFile) {
// 	this.Config = config
// 	//初始化数据库连接
// 	// this.connDB()
// 	// this.connClient()
// 	// this.deploy()
// }

//得到一个数据库连接
// func (this *ControllerImpl) GetDBConn(dbNumber int32) mysql.Conn {
// 	conn, b := this.ConnStore[dbNumber]
// 	if b {
// 		return conn
// 	}
// 	return nil
// }

//得到net模块，用于给用户发送消息
func (this *ControllerImpl) GetNet() *Net {
	return this.net
}

func (this *ControllerImpl) SetAttribute(name string, value interface{}) {
	// this.lock.Lock()
	// defer this.lock.Unlock()
	this.attributes[name] = value
}
func (this *ControllerImpl) GetAttribute(name string) interface{} {
	return this.attributes[name]
}

//
func (this *ControllerImpl) GetSession(name string) (Session, bool) {
	return this.net.GetSession(name)
}

// func (this *ControllerImpl) GetDatabaseId(accName string) (int, int) {
// 	Dbid := (int(accName[0]) + int(accName[len(accName)-1])) % message.PLAYER_DATABASE_COUNT
// 	Tbid := int(accName[len(accName)-1]) % message.PLAYER_DATABASE_COUNT
// 	return Dbid, Tbid
// }

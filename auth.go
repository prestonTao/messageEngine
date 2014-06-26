package messageEngine

import (
	"net"
	"strconv"
)

var defaultAuth Auth = new(NoneAuth)

type Auth interface {
	SendKey(conn net.Conn, session Session) (err error)
	RecvKey(conn net.Conn) (name string, err error)
}

type NoneAuth struct {
	session int64
}

//发送
func (this *NoneAuth) SendKey(conn net.Conn, session Session) (err error) {
	return
}

//接收
func (this *NoneAuth) RecvKey(conn net.Conn) (name string, err error) {
	this.session++
	// name = strconv.ParseInt(this.session, 10, )
	// name = strconv.Itoa(this.session)
	name = strconv.FormatInt(this.session, 10)
	return
}

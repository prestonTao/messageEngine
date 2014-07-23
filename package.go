package messageEngine

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"net"
	"time"
)

type GetPacket struct {
	Name  string
	MsgID int32
	Size  uint32
	Date  []byte
}

//一个packet包括包头和包体，保证在接收到包头后两秒钟内接收到包体，否则线程会一直阻塞
//因此，引入了超时机制
func RecvPackage(conn net.Conn) (packet *GetPacket, e error, isClose bool) {
	var data []byte
	header := make([]byte, 12)
	if conn == nil {
		// e = errors.New("")
		fmt.Println("连接已经关闭")
		isClose = true
		return
	}
	n, err := io.ReadFull(conn, header)
	if n == 0 && err == io.EOF {
		// fmt.Println("客户端断开连接")
		isClose = true
		e = err
		return
	} else if err != nil {
		// fmt.Println("接收数据出错:", err)
		isClose = true
		e = err
		return
	}
	//数据包长度
	size := binary.BigEndian.Uint32(header)
	//crc值
	crc1 := binary.BigEndian.Uint32(header[4:8])
	msgID := binary.BigEndian.Uint32(header[8:12])

	data = make([]byte, size)

	timeout := NewTimeOut(func() {

		n, err = io.ReadFull(conn, data)
		if uint32(n) != size {
			log.Println("数据包长度不正确", n, "!=", size)
			e = errors.New(fmt.Sprint("数据包长度不正确:%d!=%d", n, size))
			return
		}
		if err != nil {
			log.Println("读取数据出错:", err)
			e = err
			return
		}
		crc2 := crc32.Checksum(data, crc32.IEEETable)
		if crc1 != crc2 {
			log.Println("crc 数据验证不正确: ", crc1, " != ", crc2)
			e = errors.New(fmt.Sprint("crc 数据验证不正确:%d!=%d", crc1, crc2))
			return
		}
		packet = new(GetPacket)
		packet.Date = data
		packet.MsgID = int32(msgID)
		packet.Size = uint32(len(data))
	})
	isTimeOut := timeout.Do(time.Second * 5)
	if isTimeOut {
		e = errors.New("数据包头和数据包体不完整")
		return
	}
	return
}

func MarshalPacket(msgID uint32, data *[]byte) *[]byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, uint32(uint32(len(*data))))
	crc32 := crc32.Checksum(*data, crc32.IEEETable)
	binary.Write(buf, binary.BigEndian, crc32)
	binary.Write(buf, binary.BigEndian, msgID)
	buf.Write(*data)
	bs := buf.Bytes()

	return &bs
}

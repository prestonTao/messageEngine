package net

import (
	"common/packet"
	"hash/crc32"
)

var Connect int8 = 1
var Close int8 = 2

type GetPacket struct {
	ConnId uint64
	Size   uint32
	MsgID  int32
	Date   []byte
}

func PacketData(msgID uint32, data *[]byte) *[]byte {
	writer := packet.Writer()
	//size uint32
	writer.WriteU32(uint32(len(*data)))
	//crc32 uint32
	crc32 := crc32.Checksum(*data, crc32.IEEETable)
	writer.WriteU32(crc32)

	//msgID
	writer.WriteU32(msgID)
	//Data
	writer.WriteRawBytes(*data)
	bs := writer.Data()

	return &bs
}

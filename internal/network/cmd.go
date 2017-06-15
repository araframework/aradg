package network

import (
	"bytes"
	"encoding/binary"
	"github.com/araframework/aradg/internal/consts"
	"github.com/araframework/aradg/internal/consts/code"
	"github.com/araframework/aradg/internal/consts/status"
)

// command struct
type CmdHeader struct {
	// uint16
	Magic []byte
	Code  byte
}

type Member struct {
	Status    byte
	StartTime uint64
	Ip        []byte
	Port      uint16
}

type Cluster struct {
	CmdHeader
	Members []Member
}

// --- CmdJoin begin-------------
type CmdJoin struct {
	CmdHeader
	Member
}

func NewCmdJoin(me Member) []byte {
	buff := bytes.NewBuffer(make([]byte, 0))

	// header
	bMagic := make([]byte, 2)
	binary.LittleEndian.PutUint16(bMagic, consts.Magic)

	buff.Write(bMagic)        // 2
	buff.WriteByte(code.Join) //1

	// body
	bodyBuf := bytes.NewBuffer(make([]byte, 0))
	bodyBuf.WriteByte(status.New) //1

	bStartTime := make([]byte, 8)
	binary.LittleEndian.PutUint64(bStartTime, me.StartTime)
	bodyBuf.Write(bStartTime) //8

	bodyBuf.Write(me.Ip) //16

	bPort := make([]byte, 2)
	binary.LittleEndian.PutUint16(bPort, me.Port)
	bodyBuf.Write(bPort) // 2

	bodyLen := make([]byte, 4)
	binary.LittleEndian.PutUint32(bodyLen, uint32(bodyBuf.Len()))
	buff.Write(bodyLen)

	buff.Write(bodyBuf.Bytes())

	return buff.Bytes()
}

func DecodeCmdJoin(msgs []byte) CmdJoin {
	header := CmdHeader{msgs[:2], msgs[2]}

	pos := 3 + 4 // body len is 4 bytes
	status := msgs[pos]
	pos += 1
	startTime := binary.LittleEndian.Uint64(msgs[pos : pos+8])
	pos += 8
	ip := msgs[pos : pos+16]
	pos += 16
	port := binary.LittleEndian.Uint16(msgs[pos : pos+2])
	member := Member{status, startTime, ip, port}

	return CmdJoin{header, member}
}

// --- CmdJoin end-------------

func Code(code byte) byte {
	return 0
}

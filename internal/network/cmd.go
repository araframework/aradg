package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/araframework/aradg/internal/consts"
	"github.com/araframework/aradg/internal/consts/code"
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

func newCmdJoin(me Member) *CmdJoin {
	bMagic := make([]byte, 2)
	binary.LittleEndian.PutUint16(bMagic, consts.Magic)
	header := CmdHeader{bMagic, code.Join}
	return &CmdJoin{header, me}
}

// TODO
func (cmd *CmdJoin) encode() []byte {
	bodyBuf := bytes.NewBuffer(make([]byte, 30))
	bodyBuf.Write(cmd.Magic)      // 2
	bodyBuf.WriteByte(cmd.Code)   //1
	bodyBuf.WriteByte(cmd.Status) //1

	bStartTime := make([]byte, 8)
	binary.LittleEndian.PutUint64(bStartTime, cmd.StartTime)
	bodyBuf.Write(bStartTime) //8
	bodyBuf.Write(cmd.Ip)     //16

	bPort := make([]byte, 2)
	binary.LittleEndian.PutUint16(bPort, cmd.Port)
	bodyBuf.Write(bPort) // 2
	fmt.Println(bodyBuf.Len())
	fmt.Println(bodyBuf.Cap())
	return bodyBuf.Bytes()
}

// --- CmdJoin end-------------

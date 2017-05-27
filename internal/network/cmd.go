package network

import (
	"github.com/araframework/aradg/internal/consts"
	"github.com/araframework/aradg/internal/consts/code"
)

// command struct
type CmdHeader struct {
	Magic uint16
	Code  consts.Code
}

type Member struct {
	Status    consts.Status
	StartTime int64
	Ip        []byte
	Port      uint16
}

type Cluster struct {
	CmdHeader
	Members []Member
}

type CmdJoin struct {
	CmdHeader
	Member Member
}

func newCmdJoin(me Member) *CmdJoin {
	header := CmdHeader{consts.Magic, code.Join}
	return &CmdJoin{header, me}
}

// TODO
func (cmd *CmdJoin) encode() ([]byte, error) {
	var buff []byte
	return buff, nil
}

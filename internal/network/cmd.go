package network

import (
	"github.com/araframework/aradg/internal/consts"
	"github.com/araframework/aradg/internal/consts/code"
	"bytes"
)

// command struct
type Cmd struct {
	Magic   uint16
	Code    consts.Code
}

type CmdWrap struct {
	Magic   uint16
	Code    consts.Code
	Buff *bytes.Buffer
}

type Member struct {
	Status    consts.Status
	StartTime int64
	Interface string
}

type Cluster struct {
	Cmd
	Members []Member
}

type CmdJoin struct {
	Magic   uint16
	Code    consts.Code
	Member Member
}

func newCmdJoin(me Member) *CmdJoin {
	return &CmdJoin{consts.Magic, code.Join,me}
}

package network

import (
	"github.com/araframework/aradg/internal/consts"
	"github.com/araframework/aradg/internal/consts/code"
)

// request for join and cluster members
type Cluster struct {
	Leader  Member
	Members []Member
}

type Member struct {
	Status    consts.Status
	StartTime int64
	Interface string
}

// command: Join a cluster
type CmdJoin struct {
	Magic uint16
	Code  consts.Code
	Me    *Member
}

func newJoin(me *Member) *CmdJoin {
	return &CmdJoin{consts.Magic, code.Join, me}
}

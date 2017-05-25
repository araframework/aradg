package network

import (
	"github.com/araframework/aradg/internal/consts"
	"github.com/araframework/aradg/internal/consts/code"
)

type Member struct {
	Status    consts.Status
	StartTime int64
	Interface string
}

// command struct
type Cluster struct {
	Magic   uint16
	Code    consts.Code
	Members []Member
}

func newJoin(me Member) *Cluster {
	return &Cluster{consts.Magic, code.Join, []Member{me}}
}

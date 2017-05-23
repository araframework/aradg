package status

import "github.com/araframework/aradg/internal/consts"

// define the node status
const (
	New consts.Status = iota
	Member
	Leader
)

package network

import (
	"testing"
	"net"
	"github.com/araframework/aradg/internal/consts/status"
	"encoding/hex"
	"bytes"
)

func TestNewCmdJoin(t *testing.T) {
	var startTime uint64 = 1495911752212
	ip := net.ParseIP("127.0.0.1")
	port := uint16(2800)
	me := Member{status.New, startTime, ip, port}

	src := []byte("cefa001b0000000014e2494b5c01000000000000000000000000ffff7f000001f00a")

	dst := make([]byte, hex.DecodedLen(len(src)))
	_, err := hex.Decode(dst, src)
	if err != nil {
		t.Fatal(err)
	}

	join := NewCmdJoin(me)
	len := len(join)
	if join == nil || len != 34 || !bytes.Equal(join, dst) {
		t.Fatalf("NewCmdJoin fail: %x\n", join)
	}
}

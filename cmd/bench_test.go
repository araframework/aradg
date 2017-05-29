package main

import (
	"testing"
	"github.com/araframework/aradg/internal/network"
	"github.com/araframework/aradg/internal/consts/status"
	"net"
)

var res bool

func BenchmarkStrCast(b *testing.B) {
	var t uint64 = 1495911752212
	ip := net.ParseIP("127.0.0.1")
	port := uint16(2800)
	me := network.Member{Status: status.New, StartTime:t, Ip:ip, Port:port}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		network.NewCmdJoin(me)
	}
}

//func BenchmarkBytesCmp(b *testing.B) {
//
//	for n := 0; n < b.N; n++ {
//		//res = bytes.Equal(b1, b2)
//		b1 := make([]byte, 8)
//		binary.LittleEndian.PutUint64(b1, 1495911752212)
//		b2 := make([]byte, 8)
//		binary.LittleEndian.PutUint64(b2, 2495911752212)
//	}
//}

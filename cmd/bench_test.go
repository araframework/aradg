package main

import (
	"encoding/binary"
	"testing"
)

var res bool

func BenchmarkStrCast(b *testing.B) {
	var t uint64 = 1495911752212
	var t2 uint64 = 2495911752212
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		res = t == t2
	}
}

func BenchmarkBytesCmp(b *testing.B) {

	for n := 0; n < b.N; n++ {
		//res = bytes.Equal(b1, b2)
		b1 := make([]byte, 8)
		binary.LittleEndian.PutUint64(b1, 1495911752212)
		b2 := make([]byte, 8)
		binary.LittleEndian.PutUint64(b2, 2495911752212)
	}
}

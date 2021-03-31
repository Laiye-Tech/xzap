package tfmt

import (
	"testing"
	"time"
)

func Benchmark_formatTimeHeader(b *testing.B) {
	b.ReportAllocs()
	t := time.Now()
	for i := 0; i < b.N; i++ {
		var v [23]byte
		FormatTimeHeader(t, v[:])
	}
}

func Benchmark_GetDay(b *testing.B) {
	var s [23]byte
	v, _, _ := FormatTimeHeader(time.Now(), s[:])
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		GetDay(v)
	}
}

func Test_GetDay(t *testing.T) {
	var s [23]byte
	v, _, _ := FormatTimeHeader(time.Now(), s[:])
	r := GetDay(v)
	t.Log(string(r[:]))
}

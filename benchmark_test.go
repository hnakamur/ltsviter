package ltsviter

import "testing"

func BenchmarkFields(b *testing.B) {
	input := []byte("time:2025-12-17T03:46:56.123456+09:00\tua:value\\twith\\\\escapes\\n\n")
	nop := func([]byte, []byte) {}
	for b.Loop() {
		var buf [64]byte
		for field, err := range Fields(input) {
			if err != nil {
				return
			}
			nop(field.Label(), field.Value(buf[:]))
		}
	}
}

package ltsviter

import "testing"

func BenchmarkIter(b *testing.B) {
	input := []byte("time:2025-12-17T03:46:56.123456+09:00\tua:value\\twith\\\\escapes\\n\n")
	nop := func(*Entry) {}
	for b.Loop() {
		var buf [64]byte
		for entry, err := range EntryIter(input, buf[:]) {
			if err != nil {
				return
			}
			nop(&entry)
		}
	}
}

func BenchmarkRawIter(b *testing.B) {
	input := []byte("time:2025-12-17T03:46:56.123456+09:00\tua:value\\twith\\\\escapes\\n\n")
	nop := func([]byte, []byte) {}
	for b.Loop() {
		var buf [64]byte
		for entry, err := range RawEntryIter(input) {
			if err != nil {
				return
			}
			value := entry.RawValue
			if IsEscapedValue(value) {
				value = AppendUnescapedValue(buf[:0], value)
			}
			nop(entry.Label, value)
		}
	}
}

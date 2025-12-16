package ltsviter

import "fmt"

func ExampleFields() {
	input := []byte("time:2025-12-17T03:46:56.123456+09:00\tua:value\\twith\\\\escapes\\n\n")
	var buf [256]byte
	for entry, err := range Fields(input, buf[:]) {
		if err != nil {
			fmt.Printf("failed to iterate LTSV fields: %s\n", err)
			return
		}
		fmt.Printf("label=%s, value=%q\n", entry.Label, entry.Value)
	}

	// Output:
	// label=time, value="2025-12-17T03:46:56.123456+09:00"
	// label=ua, value="value\twith\\escapes\n"
}

func ExampleRawFields() {
	input := []byte("time:2025-12-17T03:46:56.123456+09:00\tua:value\\twith\\\\escapes\\n\n")
	var buf [256]byte
	for entry, err := range RawFields(input) {
		if err != nil {
			fmt.Printf("failed to iterate LTSV fields: %s\n", err)
			return
		}
		value := entry.RawValue
		if IsEscapedValue(value) {
			value = AppendUnescapedValue(buf[:0], value)
		}
		fmt.Printf("label=%s, value=%q\n", entry.Label, value)
	}

	// Output:
	// label=time, value="2025-12-17T03:46:56.123456+09:00"
	// label=ua, value="value\twith\\escapes\n"
}

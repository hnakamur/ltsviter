package ltsviter

import "fmt"

func ExampleFields() {
	input := []byte("time:2025-12-17T03:46:56.123456+09:00\tua:value\\twith\\\\escapes\\n\n")
	var buf [256]byte
	for field, err := range Fields(input, buf[:]) {
		if err != nil {
			fmt.Printf("failed to iterate LTSV fields: %s\n", err)
			return
		}
		fmt.Printf("label=%s, value=%q\n", field.Label, field.Value)
	}

	// Output:
	// label=time, value="2025-12-17T03:46:56.123456+09:00"
	// label=ua, value="value\twith\\escapes\n"
}

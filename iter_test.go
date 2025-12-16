package ltsviter

import (
	"bytes"
	"testing"
)

func TestRawFields(t *testing.T) {
	t.Run("normalCases", func(t *testing.T) {
		tsstCases := []struct {
			name  string
			input []byte
		}{
			{name: "withLF", input: []byte("time:2025-12-17T03:46:56.123456+09:00\tua:value\\twith\\\\escapes\\n\n")},
			{name: "withoutLF", input: []byte("time:2025-12-17T03:46:56.123456+09:00\tua:value\\twith\\\\escapes\\n")},
		}
		for _, tc := range tsstCases {
			t.Run(tc.name, func(t *testing.T) {
				wants := []rawField{
					{Label: []byte("time"), RawValue: []byte("2025-12-17T03:46:56.123456+09:00")},
					{Label: []byte("ua"), RawValue: []byte("value\\twith\\\\escapes\\n")},
				}
				i := 0
				for field, err := range rawFields(tc.input) {
					if err != nil {
						t.Fatalf("unexpected error: %v", err)
					}

					if got, want := field.Label, wants[i].Label; !bytes.Equal(got, want) {
						t.Errorf("label mismatch, i=%d, got=%s, want=%s", i, string(got), string(want))
					}
					if got, want := field.RawValue, wants[i].RawValue; !bytes.Equal(got, want) {
						t.Errorf("raw value mismatch, i=%d, got=%s, want=%s", i, string(got), string(want))
					}
					i++
				}
			})
		}
	})
	t.Run("errorCases", func(t *testing.T) {
		testCases := []struct {
			name  string
			input []byte
			want  string
		}{
			{
				name:  "lineEndsWithFieldSeparator",
				input: []byte("a:1\t"),
				want:  "invalid LTSV: line ends with a field separator",
			},
			{
				name:  "noLabelSeparatorInField",
				input: []byte("a"),
				want:  "invalid LTSV: no label separator in field",
			},
		}
		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				var firstError error
				for _, err := range rawFields(tc.input) {
					if err != nil {
						firstError = err
						break
					}
				}
				if firstError != nil {
					if got, want := firstError.Error(), tc.want; got != want {
						t.Errorf("error mismatched, input=%s, got=%s, want=%s", string(tc.input), got, want)
					}
				} else {
					t.Errorf("expected an error: %s for input %s, got no error", tc.want, string(tc.input))
				}
			})
		}
	})
}

func TestIsEscapedValue(t *testing.T) {
	testCases := []struct {
		name  string
		input []byte
		want  bool
	}{
		{name: "notEscaped", input: []byte("2025-12-17T03:46:56.123456+09:00"), want: false},
		{name: "escaped", input: []byte("value\\twith\\\\escapes\\n"), want: true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got, want := isEscapedValue(tc.input), tc.want; got != want {
				t.Errorf("result mismatch, input=%s, got=%v, want=%v", string(tc.input), got, want)
			}
		})
	}
}

func TestIsAppendUnescapedValue(t *testing.T) {
	testCases := []struct {
		name  string
		input []byte
		want  []byte
	}{
		{
			name:  "notEscaped",
			input: []byte("2025-12-17T03:46:56.123456+09:00"),
			want:  []byte("2025-12-17T03:46:56.123456+09:00"),
		},
		{
			name:  "escaped",
			input: []byte("value\\twith\\\\escapes\\n"),
			want:  []byte("value\twith\\escapes\n"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got, want := appendUnescapedValue([]byte{}, tc.input), tc.want; !bytes.Equal(got, want) {
				t.Errorf("result mismatch, input=%s, got=%s, want=%s", string(tc.input), string(got), string(want))
			}
		})
	}
}

// Package ltsviter provides a function to iterate over fields in a LTSV (Labeled Tab-Separated Values; http://ltsv.org/) line.
//
// This package supports an extended specification of LTSV, including value escaping.
//
// # The original LTSV specification
//
//   - Fields are separated with a tab character.
//   - A label and its corresponding value within a field are separated by a colon character.
//
// # Value Escape Extension
//
// The following special characters within values are escaped using a backslash character (\):
//   - Tab
//   - Newline
//   - Backslash
package ltsviter

import (
	"bytes"
	"fmt"
	"iter"
)

const (
	fieldSepartor = '\t'
	labelSepartor = ':'
	escapeChar    = '\\'
	newline       = '\n'
)

// Field is a field in a LTSV line.
type Field struct {
	label    []byte
	rawValue []byte
}

// Label returns the unescaped value in the Field.
func (f *Field) Label() []byte { return f.label }

// RawValue returns the raw (escaped) value in the Field.
func (f *Field) RawValue() []byte { return f.rawValue }

// Value returns the unescaped value in the Field.
func (f *Field) Value(buf []byte) []byte {
	if isEscapedValue(f.rawValue) {
		return appendUnescapedValue(buf[:0], f.rawValue)
	}
	return f.rawValue
}

// Fields returns an interator for fields in a LTSV line.
func Fields(line []byte) iter.Seq2[Field, error] {
	return func(yield func(Field, error) bool) {
		// Cut newline at end
		if len(line) > 0 && line[len(line)-1] == newline {
			line = line[:len(line)-1]
		}
		for {
			if len(line) == 0 {
				return
			}

			var field []byte
			i := bytes.IndexByte(line, fieldSepartor)
			if i == -1 {
				field, line = line, nil
			} else {
				field, line = line[:i], line[i+1:]
			}

			if j := bytes.IndexByte(field, labelSepartor); j == -1 {
				yield(Field{}, &invalidLTSVError{
					detail: noLabelSeparatorInField,
				})
				return
			} else if !yield(Field{
				label:    field[:j],
				rawValue: field[j+1:],
			}, nil) {
				return
			}

			if i != -1 && len(line) == 0 {
				yield(Field{}, &invalidLTSVError{
					detail: lineEndsWithFieldSeparator,
				})
				return
			}
		}
	}
}

type invalidLTSVError struct {
	detail string
}

func (e *invalidLTSVError) Error() string {
	return fmt.Sprintf("invalid LTSV: %s", e.detail)
}

const (
	noLabelSeparatorInField    = "no label separator in field"
	lineEndsWithFieldSeparator = "line ends with a field separator"
)

func isEscapedValue(rawValue []byte) bool {
	return bytes.IndexByte(rawValue, escapeChar) != -1
}

func appendUnescapedValue(dest, rawValue []byte) []byte {
	seenEscape := false
	for _, b := range rawValue {
		if seenEscape {
			switch b {
			case 't':
				b = '\t'
			case 'n':
				b = '\n'
			}
			dest = append(dest, b)
			seenEscape = false
		} else if b == escapeChar {
			seenEscape = true
		} else {
			dest = append(dest, b)
		}
	}
	return dest
}

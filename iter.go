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

type Entry struct {
	Label []byte
	Value []byte
}

func EntryIter(line, unescapeBuf []byte) iter.Seq2[Entry, error] {
	return func(yield func(Entry, error) bool) {
		for entry, err := range RawEntryIter(line) {
			if err != nil {
				yield(Entry{}, err)
				return
			}
			value := entry.RawValue
			if IsEscapedValue(value) {
				value = AppendUnescapedValue(unescapeBuf[:0], value)
			}
			if !yield(Entry{
				Label: entry.Label,
				Value: value,
			}, nil) {
				return
			}
		}
	}
}

type RawEntry struct {
	Label    []byte
	RawValue []byte
}

func RawEntryIter(line []byte) iter.Seq2[RawEntry, error] {
	return func(yield func(RawEntry, error) bool) {
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
				yield(RawEntry{}, &invalidLTSVError{
					detail: noLabelSeparatorInField,
				})
				return
			} else if !yield(RawEntry{
				Label:    field[:j],
				RawValue: field[j+1:],
			}, nil) {
				return
			}

			if i != -1 && len(line) == 0 {
				yield(RawEntry{}, &invalidLTSVError{
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

func IsEscapedValue(rawValue []byte) bool {
	return bytes.IndexByte(rawValue, escapeChar) != -1
}

func AppendUnescapedValue(dest, rawValue []byte) []byte {
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

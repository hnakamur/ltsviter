# ltsviter [![PkgGoDev](https://pkg.go.dev/badge/github.com/hnakamur/ltsviter)](https://pkg.go.dev/github.com/hnakamur/ltsviter)

Package ltsviter provides a function to iterate over fields in a LTSV (Labeled Tab-Separated Values; http://ltsv.org/) line.

This package supports an extended specification of LTSV, including value escaping.

## The original LTSV specification

  - Fields are separated with a tab character.
  - A label and its corresponding value within a field are separated by a colon character.

## Value Escape Extension

The following special characters within values are escaped using a backslash character (\):
  - Tab
  - Newline
  - Backslash

## License
MIT

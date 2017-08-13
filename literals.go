package goparsify

import (
	"bytes"
	"strconv"
	"unicode/utf8"
)

// StringLit matches a quoted string and returns it in .Result. It may contain:
//  - unicode
//  - escaped characters, eg \" or \n
//  - unicode sequences, eg \uBEEF
func StringLit(allowedQuotes string) Parser {
	return NewParser("string literal", func(ps *State, node *Result) {
		ps.WS(ps)

		if !stringContainsByte(allowedQuotes, ps.Input[ps.Pos]) {
			ps.ErrorHere(allowedQuotes)
			return
		}
		quote := ps.Input[ps.Pos]

		var end = ps.Pos + 1

		inputLen := len(ps.Input)
		var buf *bytes.Buffer

		for end < inputLen {
			switch ps.Input[end] {
			case '\\':
				if end+1 >= inputLen {
					ps.ErrorHere(string(quote))
					return
				}

				if buf == nil {
					buf = bytes.NewBufferString(ps.Input[ps.Pos+1 : end])
				}

				c := ps.Input[end+1]
				if c == 'u' {
					if end+6 >= inputLen {
						ps.Error.expected = "[a-f0-9]{4}"
						ps.Error.pos = end + 2
						return
					}

					r, ok := unhex(ps.Input[end+2 : end+6])
					if !ok {
						ps.Error.expected = "[a-f0-9]"
						ps.Error.pos = end + 2
						return
					}
					buf.WriteRune(r)
					end += 6
				} else {
					buf.WriteByte(c)
					end += 2
				}
			case quote:
				if buf == nil {
					node.Result = ps.Input[ps.Pos+1 : end]
					ps.Pos = end + 1
					return
				}
				ps.Pos = end + 1
				node.Result = buf.String()
				return
			default:
				if buf == nil {
					if ps.Input[end] < 127 {
						end++
					} else {
						_, w := utf8.DecodeRuneInString(ps.Input[end:])
						end += w
					}
				} else {
					r, w := utf8.DecodeRuneInString(ps.Input[end:])
					end += w
					buf.WriteRune(r)
				}
			}
		}

		ps.ErrorHere(string(quote))
	})
}

// NumberLit matches a floating point or integer number and returns it as a int64 or float64 in .Result
func NumberLit() Parser {
	return NewParser("number literal", func(ps *State, node *Result) {
		ps.WS(ps)
		end := ps.Pos
		float := false
		inputLen := len(ps.Input)

		if end < inputLen && (ps.Input[end] == '-' || ps.Input[end] == '+') {
			end++
		}

		for end < inputLen && ps.Input[end] >= '0' && ps.Input[end] <= '9' {
			end++
		}

		if end < inputLen && ps.Input[end] == '.' {
			float = true
			end++
		}

		for end < inputLen && ps.Input[end] >= '0' && ps.Input[end] <= '9' {
			end++
		}

		if end < inputLen && (ps.Input[end] == 'e' || ps.Input[end] == 'E') {
			end++
			float = true

			if end < inputLen && (ps.Input[end] == '-' || ps.Input[end] == '+') {
				end++
			}

			for end < inputLen && ps.Input[end] >= '0' && ps.Input[end] <= '9' {
				end++
			}
		}

		if end == ps.Pos {
			ps.ErrorHere("number")
			return
		}

		var err error
		if float {
			node.Result, err = strconv.ParseFloat(ps.Input[ps.Pos:end], 10)
		} else {
			node.Result, err = strconv.ParseInt(ps.Input[ps.Pos:end], 10, 64)
		}
		if err != nil {
			ps.ErrorHere("number")
			return
		}
		ps.Pos = end
	})
}

func stringContainsByte(s string, b byte) bool {
	for i := 0; i < len(s); i++ {
		if b == s[i] {
			return true
		}
	}
	return false
}

func unhex(b string) (v rune, ok bool) {
	for _, c := range b {
		v <<= 4
		switch {
		case '0' <= c && c <= '9':
			v |= c - '0'
		case 'a' <= c && c <= 'f':
			v |= c - 'a' + 10
		case 'A' <= c && c <= 'F':
			v |= c - 'A' + 10
		default:
			return 0, false
		}
	}

	return v, true
}

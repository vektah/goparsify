package goparsify

import (
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"
)

type Node struct {
	Token    string
	Children []Node
	Result   interface{}
}

type Parser func(*State) Node

// Parserish types are any type that can be turned into a Parser by Parsify
// These currently include *Parser and string literals.
//
// This makes recursive grammars cleaner and allows string literals to be used directly in most contexts.
// eg, matching balanced paren:
// ```go
// var group Parser
// group = And("(", Maybe(&group), ")")
// ```
// vs
// ```go
// var group ParserPtr{}
// group.P = And(Exact("("), Maybe(group.Parse), Exact(")"))
// ```
type Parserish interface{}

func Parsify(p Parserish) Parser {
	switch p := p.(type) {
	case nil:
		return nil
	case func(*State) Node:
		return NewParser("anonymous func", p)
	case Parser:
		return p
	case *Parser:
		// Todo: Maybe capture this stack and on nil show it? Is there a good error library to do this?
		return func(ptr *State) Node {
			return (*p)(ptr)
		}
	case string:
		return Exact(p)
	default:
		panic(fmt.Errorf("cant turn a `%T` into a parser", p))
	}
}

func ParsifyAll(parsers ...Parserish) []Parser {
	ret := make([]Parser, len(parsers))
	for i, parser := range parsers {
		ret[i] = Parsify(parser)
	}
	return ret
}

func WS() Parser {
	return NewParser("AutoWS", func(ps *State) Node {
		ps.WS()
		return Node{}
	})
}

func ParseString(parser Parserish, input string) (result interface{}, remaining string, err error) {
	p := Parsify(parser)
	ps := InputString(input)

	ret := p(ps)
	ps.AutoWS()

	if ps.Error.Expected != "" {
		return nil, ps.Get(), ps.Error
	}

	return ret.Result, ps.Get(), nil
}

func Exact(match string) Parser {
	if len(match) == 1 {
		matchByte := match[0]
		return NewParser(match, func(ps *State) Node {
			ps.AutoWS()
			if ps.Pos >= len(ps.Input) || ps.Input[ps.Pos] != matchByte {
				ps.ErrorHere(match)
				return Node{}
			}

			ps.Advance(1)

			return Node{Token: match}
		})
	}

	return NewParser(match, func(ps *State) Node {
		ps.AutoWS()
		if !strings.HasPrefix(ps.Get(), match) {
			ps.ErrorHere(match)
			return Node{}
		}

		ps.Advance(len(match))

		return Node{Token: match}
	})
}

func parseRepetition(defaultMin, defaultMax int, repetition ...int) (min int, max int) {
	min = defaultMin
	max = defaultMax
	switch len(repetition) {
	case 0:
	case 1:
		min = repetition[0]
	case 2:
		min = repetition[0]
		max = repetition[1]
	default:
		panic(fmt.Errorf("Dont know what %d repetion args mean", len(repetition)))
	}
	return min, max
}

// parseMatcher turns a string in the format a-f01234A-F into:
//   - a set string of matches string(01234)
//   - a set of ranges [][]rune{{'a', 'f'}, {'A', 'F'}}
func parseMatcher(matcher string) (matches string, ranges [][]rune) {
	runes := []rune(matcher)

	for i := 0; i < len(runes); i++ {

		if i+2 < len(runes) && runes[i+1] == '-' {
			start := runes[i]
			end := runes[i+2]
			if start <= end {
				ranges = append(ranges, []rune{start, end})
			} else {
				ranges = append(ranges, []rune{end, start})
			}
		} else if i+1 < len(runes) && runes[i] == '\\' {
			matches += string(runes[i+1])
		} else {
			matches += string(runes[i])
		}

	}

	return matches, ranges
}

func Chars(matcher string, repetition ...int) Parser {
	return NewParser("["+matcher+"]", charsImpl(matcher, false, repetition...))
}

func NotChars(matcher string, repetition ...int) Parser {
	return NewParser("!["+matcher+"]", charsImpl(matcher, true, repetition...))
}

func charsImpl(matcher string, stopOn bool, repetition ...int) Parser {
	min, max := parseRepetition(1, -1, repetition...)
	matches, ranges := parseMatcher(matcher)

	return func(ps *State) Node {
		ps.AutoWS()
		matched := 0
		for ps.Pos+matched < len(ps.Input) {
			if max != -1 && matched >= max {
				break
			}

			r, w := utf8.DecodeRuneInString(ps.Input[ps.Pos+matched:])

			anyMatched := strings.ContainsRune(matches, r)
			if !anyMatched {
				for _, rng := range ranges {
					if r >= rng[0] && r <= rng[1] {
						anyMatched = true
					}
				}
			}

			if anyMatched == stopOn {
				break
			}

			matched += w
		}

		if matched < min {
			ps.ErrorHere(matcher)
			return Node{}
		}

		result := ps.Input[ps.Pos : ps.Pos+matched]
		ps.Advance(matched)
		return Node{Token: result}
	}
}

func String(quote rune) Parser {
	return NewParser("string", func(ps *State) Node {
		ps.AutoWS()
		var r rune
		var w int
		var matched int
		r, matched = utf8.DecodeRuneInString(ps.Input[ps.Pos:])
		if r != quote {
			ps.ErrorHere("\"")
			return Node{}
		}

		result := &bytes.Buffer{}

		for ps.Pos+matched < len(ps.Input) {
			r, w = utf8.DecodeRuneInString(ps.Input[ps.Pos+matched:])
			matched += w

			if r == '\\' {
				r, w = utf8.DecodeRuneInString(ps.Input[ps.Pos+matched:])
				result.WriteRune(r)
				matched += w
				continue
			}

			if r == quote {
				ps.Advance(matched)
				return Node{Token: result.String()}
			}
			result.WriteRune(r)
		}

		ps.ErrorHere("\"")
		return Node{}
	})
}

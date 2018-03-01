package goparsify

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Parser is the workhorse of parsify. A parser takes a State and returns a result, consuming some
// of the State in the process.
// Given state is shared there are a few rules that should be followed:
//  - A parser that errors must set state.Error
//  - A parser that errors must not change state.Pos
//  - A parser that consumed some input should advance state.Pos
type Parser func(*State, *Result)

// Map shorthand for Map(p, func())
func (p Parser) Map(f func(n *Result)) Parser {
	return Map(p, f)
}

// VoidParser is a special type of parser that never returns anything but can still consume input
type VoidParser func(*State)

// Parserish types are any type that can be turned into a Parser by Parsify
// These currently include *Parser and string literals.
//
// This makes recursive grammars cleaner and allows string literals to be used directly in most contexts.
// eg, matching balanced paren:
//  var group Parser
//  group = Seq("(", Maybe(&group), ")")
// vs
//  var group ParserPtr{}
//  group.P = Seq(Exact("("), Maybe(group.Parse), Exact(")"))
type Parserish interface{}

// Parsify takes a Parserish and makes a Parser out of it. It should be called by
// any Parser that accepts a Parser as an argument. It should never be called during
// instead call it during parser creation so there is no runtime cost.
//
// See Parserish for details.
func Parsify(p Parserish) Parser {
	switch p := p.(type) {
	case func(*State, *Result):
		return p
	case Parser:
		return p
	case *Parser:
		// Todo: Maybe capture this stack and on nil show it? Is there a good error library to do this?
		return func(ptr *State, node *Result) {
			(*p)(ptr, node)
		}
	case string:
		return Exact(p)
	case func(*State):
		return func(ptr *State, node *Result) {
			p(ptr)
		}
	default:
		panic(fmt.Errorf("cant turn a `%T` into a parser", p))
	}
}

// ParsifyAll calls Parsify on all parsers
func ParsifyAll(parsers ...Parserish) []Parser {
	ret := make([]Parser, len(parsers))
	for i, parser := range parsers {
		ret[i] = Parsify(parser)
	}
	return ret
}

// Run applies some input to a parser and returns the result, failing if the input isnt fully consumed.
// It is a convenience method for the most common way to invoke a parser.
func Run(parser Parserish, input string, ws ...VoidParser) (result interface{}, err error) {
	p := Parsify(parser)
	ps := NewState(input)
	if len(ws) > 0 {
		ps.WS = ws[0]
	}

	ret := Result{}
	p(ps, &ret)
	ps.WS(ps)

	if ps.Error.expected != "" {
		return ret.Result, &ps.Error
	}

	if ps.Get() != "" {
		return ret.Result, UnparsedInputError{ps.Get()}
	}

	return ret.Result, nil
}

// Cut prevents backtracking beyond this point. Usually used after keywords when you
// are sure this is the correct path. Improves performance and error reporting.
func Cut() Parser {
	return func(ps *State, node *Result) {
		ps.Cut = ps.Pos
	}
}

// Regex returns a match if the regex successfully matches
func Regex(pattern string) Parser {
	re := regexp.MustCompile("^" + pattern)
	return NewParser(pattern, func(ps *State, node *Result) {
		ps.WS(ps)
		if match := re.FindString(ps.Get()); match != "" {
			ps.Advance(len(match))
			node.Token = match
			return
		}
		ps.ErrorHere(pattern)
	})
}

// Exact will fully match the exact string supplied, or error. The match will be stored in .Token
func Exact(match string) Parser {
	if len(match) == 1 {
		matchByte := match[0]
		return NewParser(match, func(ps *State, node *Result) {
			ps.WS(ps)
			if ps.Pos >= len(ps.Input) || ps.Input[ps.Pos] != matchByte {
				ps.ErrorHere(match)
				return
			}

			ps.Advance(1)

			node.Token = match
		})
	}

	return NewParser(match, func(ps *State, node *Result) {
		ps.WS(ps)
		if !strings.HasPrefix(ps.Get(), match) {
			ps.ErrorHere(match)
			return
		}

		ps.Advance(len(match))

		node.Token = match
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
		panic(fmt.Errorf("Dont know what %d repetition args mean", len(repetition)))
	}
	return min, max
}

// parseMatcher turns a string in the format a-f01234A-F into:
//   - an alphabet of matches string(01234)
//   - a set of ranges [][]rune{{'a', 'f'}, {'A', 'F'}}
func parseMatcher(matcher string) (alphabet string, ranges [][]rune) {
	runes := []rune(matcher)

	i := 0
	for {
		if i >= len(runes) {
			break
		}

		if i+2 < len(runes) && runes[i+1] == '-' && runes[i] != '\\' {
			start := runes[i]
			end := runes[i+2]
			if start <= end {
				ranges = append(ranges, []rune{start, end})
			} else {
				ranges = append(ranges, []rune{end, start})
			}
			i += 3 // we just consumed 3 bytes: range start, hyphen, and range end
			continue
		}

		if i+1 < len(runes) && runes[i] == '\\' {
			alphabet += string(runes[i+1])
		} else {
			alphabet += string(runes[i])
		}

		i++
	}

	return alphabet, ranges
}

// Chars is the swiss army knife of character matches. It can match:
//  - ranges: Chars("a-z") will match one or more lowercase letter
//  - alphabets: Chars("abcd") will match one or more of the letters abcd in any order
//  - min and max: Chars("a-z0-9", 4, 6) will match 4-6 lowercase alphanumeric characters
// the above can be combined in any order
func Chars(matcher string, repetition ...int) Parser {
	return NewParser("["+matcher+"]", charsImpl(matcher, false, repetition...))
}

// NotChars accepts the full range of input from Chars, but it will stop when any
// character matches. If you need to match until you see a sequence use Until instead
func NotChars(matcher string, repetition ...int) Parser {
	return NewParser("!["+matcher+"]", charsImpl(matcher, true, repetition...))
}

func charsImpl(matcher string, stopOn bool, repetition ...int) Parser {
	min, max := parseRepetition(1, -1, repetition...)
	alphabet, ranges := parseMatcher(matcher)

	return func(ps *State, node *Result) {
		ps.WS(ps)
		matched := 0
		for ps.Pos+matched < len(ps.Input) {
			if max != -1 && matched >= max {
				break
			}

			r, w := utf8.DecodeRuneInString(ps.Input[ps.Pos+matched:])

			anyMatched := strings.ContainsRune(alphabet, r)
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
			return
		}

		node.Token = ps.Input[ps.Pos : ps.Pos+matched]
		ps.Advance(matched)
	}
}

// Until will consume all input until one of the given terminator sequences is found. If you want to stop when seeing
// single characters see NotChars instead
func Until(terminators ...string) Parser {

	return NewParser("Until", func(ps *State, node *Result) {
		startPos := ps.Pos
	loop:
		for ps.Pos < len(ps.Input) {
			for _, terminator := range terminators {
				if ps.Pos+len(terminator) <= len(ps.Input) && ps.Input[ps.Pos:ps.Pos+len(terminator)] == terminator {
					break loop
				}
			}
			ps.Pos++
		}

		if ps.Pos == startPos {
			ps.ErrorHere("something")
		}
		node.Token = ps.Input[startPos:ps.Pos]
	})
}

package parsec

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type Parser func(Pointer) (Node, Pointer)

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
	case func(Pointer) (Node, Pointer):
		return Parser(p)
	case Parser:
		return p
	case *Parser:
		// Todo: Maybe capture this stack and on nil show it? Is there a good error library to do this?
		return func(ptr Pointer) (Node, Pointer) {
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

func ParseString(parser Parserish, input string) (result Node, remaining string, err error) {
	p := Parsify(parser)
	result, pointer := p(Pointer{input, 0})

	if err, isErr := result.(Error); isErr {
		return nil, pointer.Get(), err
	}

	return result, pointer.Get(), nil
}

func Exact(match string) Parser {
	return func(p Pointer) (Node, Pointer) {
		if !strings.HasPrefix(p.Get(), match) {
			return NewError(p.pos, "Expected "+match), p
		}

		return NewToken(p.pos, match), p.Advance(len(match))
	}
}

func Char(match string) Parser {
	return func(p Pointer) (Node, Pointer) {
		r, w := utf8.DecodeRuneInString(p.Get())

		if !strings.ContainsRune(match, r) {
			return NewError(p.pos, "Expected one of "+string(match)), p

		}
		return NewToken(p.pos, string(r)), p.Advance(w)
	}
}

func CharRun(match string) Parser {
	return func(p Pointer) (Node, Pointer) {
		matched := 0
		for p.pos+matched < len(p.input) {
			r, w := utf8.DecodeRuneInString(p.input[p.pos+matched:])
			if !strings.ContainsRune(match, r) {
				break
			}
			matched += w
		}

		if matched == 0 {
			return NewError(p.pos, "Expected some of "+match), p
		}

		return NewToken(p.pos, p.input[p.pos:p.pos+matched]), p.Advance(matched)
	}
}

func CharRunUntil(match string) Parser {
	return func(p Pointer) (Node, Pointer) {
		matched := 0
		for p.pos+matched < len(p.input) {
			r, w := utf8.DecodeRuneInString(p.input[p.pos+matched:])
			if strings.ContainsRune(match, r) {
				break
			}
			matched += w
		}

		if matched == 0 {
			return NewError(p.pos, "Expected some of "+match), p
		}

		return NewToken(p.pos, p.input[p.pos:p.pos+matched]), p.Advance(matched)
	}
}

func Range(r string, repetition ...int) Parser {
	min := int(1)
	max := int(-1)
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

	runes := []rune(r)
	if len(runes)%3 != 0 {
		panic("ranges should be in the form a-z0-9")
	}

	var ranges [][]rune
	for i := 0; i < len(runes); i += 3 {
		start := runes[i]
		end := runes[i+2]
		if start <= end {
			ranges = append(ranges, []rune{start, end})
		} else {
			ranges = append(ranges, []rune{end, start})
		}
	}

	return func(p Pointer) (Node, Pointer) {
		matched := 0
		for p.pos+matched < len(p.input) {
			if max != -1 && matched >= max {
				break
			}

			r, w := utf8.DecodeRuneInString(p.input[p.pos+matched:])

			anyMatched := false
			for _, rng := range ranges {
				if r >= rng[0] && r <= rng[1] {
					anyMatched = true
				}
			}
			if !anyMatched {
				break
			}

			matched += w
		}

		if matched < min {
			return NewError(p.pos+matched, fmt.Sprintf("Expected at least %d more of %s", min-matched, r)), p
		}

		return NewToken(p.pos, p.input[p.pos:p.pos+matched]), p.Advance(matched)
	}
}

func WS(p Pointer) (Node, Pointer) {
	_, p2 := CharRun("\t\n\v\f\r \x85\xA0")(p)
	return nil, p2
}

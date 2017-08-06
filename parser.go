package parsec

import (
	"fmt"
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

func Exact(match string) Parser {
	return func(p Pointer) (Node, Pointer) {
		if !p.HasPrefix(match) {
			return NewError(p.pos, "Expected "+match), p
		}

		return NewToken(p.pos, match), p.Advance(len(match))
	}
}

func Char(match string) Parser {
	return func(p Pointer) (Node, Pointer) {
		r, p2 := p.Accept(match)
		if r == "" {
			return NewError(p.pos, "Expected one of "+string(match)), p
		}

		return NewToken(p.pos, string(r)), p2
	}
}

func CharRun(match string) Parser {
	return func(p Pointer) (Node, Pointer) {
		s, p2 := p.AcceptRun(match)
		if s == "" {
			return NewError(p.pos, "Expected some of "+match), p
		}

		return NewToken(p.pos, s), p2
	}
}

func CharRunUntil(match string) Parser {
	return func(p Pointer) (Node, Pointer) {
		s, p2 := p.AcceptUntil(match)
		if s == "" {
			return NewError(p.pos, "Expected some of "+match), p
		}

		return NewToken(p.pos, s), p2
	}
}

func Range(r string) string {
	runes := []rune(r)
	if len(runes)%3 != 0 {
		panic("ranges should be in the form a-z0-9")
	}

	match := ""

	for i := 0; i < len(runes); i += 3 {
		start := runes[i]
		end := runes[i+2]
		if start > end {
			tmp := start
			start = end
			end = tmp
		}
		for c := start; c <= end; c++ {
			match += string(c)
		}
	}

	return match
}

func WS(p Pointer) (Node, Pointer) {
	_, p2 := p.AcceptRun("\t\n\v\f\r \x85\xA0")

	return nil, p2
}

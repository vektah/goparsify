package goparsify

import (
	"bytes"
)

// Seq matches all of the given parsers in order and returns their result as .Child[n]
func Seq(parsers ...Parserish) Parser {
	parserfied := ParsifyAll(parsers...)

	return NewParser("Seq()", func(ps *State) Result {
		result := Result{Child: make([]Result, len(parserfied))}
		startpos := ps.Pos
		for i, parser := range parserfied {
			result.Child[i] = parser(ps)
			if ps.Errored() {
				ps.Pos = startpos
				return result
			}
		}
		return result
	})
}

// NoAutoWS disables automatically ignoring whitespace between tokens for all parsers underneath
func NoAutoWS(parser Parserish) Parser {
	parserfied := Parsify(parser)
	return func(ps *State) Result {
		ps.NoAutoWS = true

		ret := parserfied(ps)

		ps.NoAutoWS = false
		return ret
	}
}

// Any matches the first successful parser and returns its result
func Any(parsers ...Parserish) Parser {
	parserfied := ParsifyAll(parsers...)

	return NewParser("Any()", func(ps *State) Result {
		longestError := Error{}
		startpos := ps.Pos
		for _, parser := range parserfied {
			node := parser(ps)
			if ps.Errored() {
				if ps.Error.pos > longestError.pos {
					longestError = ps.Error
				}
				if ps.Cut > startpos {
					break
				}
				ps.Recover()
				continue
			}
			return node
		}

		ps.Error = longestError
		ps.Pos = startpos
		return Result{}
	})
}

// Some matches one or more parsers and returns the value as .Child[n]
// an optional separator can be provided and that value will be consumed
// but not returned. Only one separator can be provided.
func Some(parser Parserish, separator ...Parserish) Parser {
	return NewParser("Some()", manyImpl(0, parser, separator...))
}

// Many matches zero or more parsers and returns the value as .Child[n]
// an optional separator can be provided and that value will be consumed
// but not returned. Only one separator can be provided.
func Many(parser Parserish, separator ...Parserish) Parser {
	return NewParser("Many()", manyImpl(1, parser, separator...))
}

func manyImpl(min int, op Parserish, sep ...Parserish) Parser {
	var opParser = Parsify(op)
	var sepParser Parser
	if len(sep) > 0 {
		sepParser = Parsify(sep[0])
	}

	return func(ps *State) Result {
		var result Result
		startpos := ps.Pos
		for {
			node := opParser(ps)
			if ps.Errored() {
				if len(result.Child) < min || ps.Cut > ps.Pos {
					ps.Pos = startpos
					return result
				}
				ps.Recover()
				return result
			}
			result.Child = append(result.Child, node)

			if sepParser != nil {
				sepParser(ps)
				if ps.Errored() {
					ps.Recover()
					return result
				}
			}
		}
	}
}

// Maybe will 0 or 1 of the parser
func Maybe(parser Parserish) Parser {
	parserfied := Parsify(parser)

	return NewParser("Maybe()", func(ps *State) Result {
		startpos := ps.Pos
		node := parserfied(ps)
		if ps.Errored() && ps.Cut <= startpos {
			ps.Recover()
		}

		return node
	})
}

// Bind will set the node .Result when the given parser matches
// This is useful for giving a value to keywords and constant literals
// like true and false. See the json parser for an example.
func Bind(parser Parserish, val interface{}) Parser {
	p := Parsify(parser)

	return func(ps *State) Result {
		node := p(ps)
		if ps.Errored() {
			return node
		}
		node.Result = val
		return node
	}
}

// Map applies the callback if the parser matches. This is used to set the Result
// based on the matched result.
func Map(parser Parserish, f func(n Result) Result) Parser {
	p := Parsify(parser)

	return func(ps *State) Result {
		node := p(ps)
		if ps.Errored() {
			return node
		}
		return f(node)
	}
}

func flatten(n Result) string {
	if n.Token != "" {
		return n.Token
	}

	if len(n.Child) > 0 {
		sbuf := &bytes.Buffer{}
		for _, node := range n.Child {
			sbuf.WriteString(flatten(node))
		}
		return sbuf.String()
	}

	return ""
}

// Merge all child Tokens together recursively
func Merge(parser Parserish) Parser {
	return Map(parser, func(n Result) Result {
		return Result{Token: flatten(n)}
	})
}

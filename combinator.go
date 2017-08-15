package goparsify

import (
	"bytes"
)

// Seq matches all of the given parsers in order and returns their result as .Child[n]
func Seq(parsers ...Parserish) Parser {
	parserfied := ParsifyAll(parsers...)

	return NewParser("Seq()", func(ps *State, node *Result) {
		node.Child = make([]Result, len(parserfied))
		startpos := ps.Pos
		for i, parser := range parserfied {
			parser(ps, &node.Child[i])
			if ps.Errored() {
				ps.Pos = startpos
				return
			}
		}
	})
}

// NoAutoWS disables automatically ignoring whitespace between tokens for all parsers underneath
func NoAutoWS(parser Parserish) Parser {
	parserfied := Parsify(parser)
	return func(ps *State, node *Result) {
		oldWS := ps.WS
		ps.WS = NoWhitespace
		parserfied(ps, node)
		ps.WS = oldWS
	}
}

// Any matches the first successful parser and returns its result
func Any(parsers ...Parserish) Parser {
	parserfied := ParsifyAll(parsers...)
	// Records which parser was successful for each byte, and will use it first next time.
	predictor := [255]int{}

	return NewParser("Any()", func(ps *State, node *Result) {
		ps.WS(ps)
		if ps.Pos >= len(ps.Input) {
			ps.ErrorHere("!EOF")
			return
		}
		startpos := ps.Pos
		predictorChar := ps.Input[startpos]
		predicted := predictor[predictorChar]

		parserfied[predicted](ps, node)
		if !ps.Errored() {
			return
		}

		longestError := ps.Error
		if ps.Cut <= startpos {
			ps.Recover()
		} else {
			return
		}

		for i, parser := range parserfied {
			if i == predicted {
				continue
			}
			parser(ps, node)
			if ps.Errored() {
				if ps.Error.pos >= longestError.pos {
					longestError = ps.Error
				}
				if ps.Cut > startpos {
					break
				}
				ps.Recover()
				continue
			}
			predictor[predictorChar] = i
			return
		}

		ps.Error = longestError
		ps.Pos = startpos
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

	return func(ps *State, node *Result) {
		node.Child = make([]Result, 0, 5)
		startpos := ps.Pos
		for {
			node.Child = append(node.Child, Result{})
			opParser(ps, &node.Child[len(node.Child)-1])
			if ps.Errored() {
				if len(node.Child)-1 < min || ps.Cut > ps.Pos {
					ps.Pos = startpos
					return
				}
				ps.Recover()
				node.Child = node.Child[0 : len(node.Child)-1]
				return
			}

			if sepParser != nil {
				sepParser(ps, TrashResult)
				if ps.Errored() {
					ps.Recover()
					return
				}
			}
		}
	}
}

// Maybe will 0 or 1 of the parser
func Maybe(parser Parserish) Parser {
	parserfied := Parsify(parser)

	return NewParser("Maybe()", func(ps *State, node *Result) {
		startpos := ps.Pos
		parserfied(ps, node)
		if ps.Errored() && ps.Cut <= startpos {
			ps.Recover()
		}
	})
}

// Bind will set the node .Result when the given parser matches
// This is useful for giving a value to keywords and constant literals
// like true and false. See the json parser for an example.
func Bind(parser Parserish, val interface{}) Parser {
	p := Parsify(parser)

	return func(ps *State, node *Result) {
		p(ps, node)
		if ps.Errored() {
			return
		}
		node.Result = val
	}
}

// Map applies the callback if the parser matches. This is used to set the Result
// based on the matched result.
func Map(parser Parserish, f func(n *Result)) Parser {
	p := Parsify(parser)

	return func(ps *State, node *Result) {
		p(ps, node)
		if ps.Errored() {
			return
		}
		f(node)
	}
}

func flatten(n *Result) {
	if len(n.Child) > 0 {
		sbuf := &bytes.Buffer{}
		for _, child := range n.Child {
			flatten(&child)
			sbuf.WriteString(child.Token)
		}
		n.Token = sbuf.String()
	}
}

// Merge all child Tokens together recursively
func Merge(parser Parserish) Parser {
	return Map(parser, flatten)
}

package goparsify

import (
	"bytes"
)

var Nil = NewParser("Nil", func(ps *State) Node {
	return Node{}
})

func And(parsers ...Parserish) Parser {
	if len(parsers) == 0 {
		return Nil
	}

	parserfied := ParsifyAll(parsers...)

	return NewParser("And()", func(ps *State) Node {
		result := Node{Children: make([]Node, len(parserfied))}
		startpos := ps.Pos
		for i, parser := range parserfied {
			result.Children[i] = parser(ps)
			if ps.Errored() {
				ps.Pos = startpos
				return result
			}
		}
		return result
	})
}

func NoAutoWS(parser Parserish) Parser {
	parserfied := Parsify(parser)
	return func(ps *State) Node {
		ps.NoAutoWS = true

		ret := parserfied(ps)

		ps.NoAutoWS = false
		return ret
	}
}

func Any(parsers ...Parserish) Parser {
	if len(parsers) == 0 {
		return Nil
	}

	parserfied := ParsifyAll(parsers...)

	return NewParser("Any()", func(ps *State) Node {
		longestError := Error{}
		startpos := ps.Pos
		for _, parser := range parserfied {
			node := parser(ps)
			if ps.Errored() {
				if ps.Error.pos > longestError.pos {
					longestError = ps.Error
				}
				ps.ClearError()
				continue
			}
			return node
		}

		ps.Error = longestError
		ps.Pos = startpos
		return Node{}
	})
}

func Kleene(opScan Parserish, sepScan ...Parserish) Parser {
	return NewParser("Kleene()", manyImpl(0, opScan, sepScan...))
}

func Many(opScan Parserish, sepScan ...Parserish) Parser {
	return NewParser("Many()", manyImpl(1, opScan, sepScan...))
}

func manyImpl(min int, op Parserish, sep ...Parserish) Parser {
	var opParser = Parsify(op)
	var sepParser Parser
	if len(sep) > 0 {
		sepParser = Parsify(sep[0])
	}

	return func(ps *State) Node {
		var result Node
		startpos := ps.Pos
		for {
			node := opParser(ps)
			if ps.Errored() {
				if len(result.Children) < min {
					ps.Pos = startpos
					return result
				}
				ps.ClearError()
				return result
			}
			result.Children = append(result.Children, node)

			if sepParser != nil {
				sepParser(ps)
				if ps.Errored() {
					ps.ClearError()
					return result
				}
			}
		}
	}
}

func Maybe(parser Parserish) Parser {
	parserfied := Parsify(parser)

	return NewParser("Maybe()", func(ps *State) Node {
		node := parserfied(ps)
		if ps.Errored() {
			ps.ClearError()
		}

		return node
	})
}

func Map(parser Parserish, f func(n Node) Node) Parser {
	p := Parsify(parser)

	return NewParser("Map()", func(ps *State) Node {
		node := p(ps)
		if ps.Errored() {
			return node
		}
		return f(node)
	})
}

func flatten(n Node) string {
	if n.Token != "" {
		return n.Token
	}

	if len(n.Children) > 0 {
		sbuf := &bytes.Buffer{}
		for _, node := range n.Children {
			sbuf.WriteString(flatten(node))
		}
		return sbuf.String()
	}

	return ""
}

func Merge(parser Parserish) Parser {
	return NewParser("Merge()", Map(parser, func(n Node) Node {
		return Node{Token: flatten(n)}
	}))
}

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
		var nodes = make([]Node, 0, len(parserfied))
		startpos := ps.Pos
		for _, parser := range parserfied {
			node := parser(ps)
			if ps.Errored() {
				ps.Pos = startpos
				return Node{}
			}
			nodes = append(nodes, node)
		}
		return Node{Children: nodes}
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
	return NewParser("Kleene()", manyImpl(0, opScan, nil, sepScan...))
}

func KleeneUntil(opScan Parserish, untilScan Parserish, sepScan ...Parserish) Parser {
	return NewParser("KleeneUntil()", manyImpl(0, opScan, untilScan, sepScan...))
}

func Many(opScan Parserish, sepScan ...Parserish) Parser {
	return NewParser("Many()", manyImpl(1, opScan, nil, sepScan...))
}

func ManyUntil(opScan Parserish, untilScan Parserish, sepScan ...Parserish) Parser {
	return NewParser("ManyUntil()", manyImpl(1, opScan, untilScan, sepScan...))
}

func manyImpl(min int, op Parserish, until Parserish, sep ...Parserish) Parser {
	opParser := Parsify(op)
	untilParser := Parsify(until)
	sepParser := Nil
	if len(sep) > 0 {
		sepParser = Parsify(sep[0])
	}

	return func(ps *State) Node {
		var node Node
		nodes := make([]Node, 0, 20)
		startpos := ps.Pos
		for {
			tempPos := ps.Pos
			if untilParser != nil {
				node = untilParser(ps)
				if !ps.Errored() {
					ps.Pos = tempPos
					if len(nodes) < min {
						ps.Pos = startpos
						ps.ErrorHere("something else")
						return Node{}
					}
					break
				}
				ps.ClearError()
			}

			node = opParser(ps)
			if ps.Errored() {
				if len(nodes) < min {
					ps.Pos = startpos
					return Node{}
				}
				ps.ClearError()
				break
			}

			nodes = append(nodes, node)

			node = sepParser(ps)
			if ps.Errored() {
				ps.ClearError()
				break
			}
		}
		return Node{Children: nodes}
	}
}

func Maybe(parser Parserish) Parser {
	parserfied := Parsify(parser)

	return NewParser("Maybe()", func(ps *State) Node {
		node := parserfied(ps)
		if ps.Errored() {
			ps.ClearError()
			return Node{}
		}

		return node
	})
}

func Map(parser Parserish, f func(n Node) Node) Parser {
	p := Parsify(parser)

	return NewParser("Map()", func(ps *State) Node {
		node := p(ps)
		if ps.Errored() {
			return Node{}
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

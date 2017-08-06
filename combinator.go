package parsec

import (
	"bytes"
	"fmt"
)

func Nil(ps *State) interface{} {
	return nil
}

func Never(ps *State) interface{} {
	ps.ErrorHere("not anything")
	return nil
}

func And(parsers ...Parserish) Parser {
	if len(parsers) == 0 {
		return Nil
	}

	parserfied := ParsifyAll(parsers...)

	return func(ps *State) interface{} {
		var nodes = make([]interface{}, 0, len(parserfied))
		startpos := ps.Pos
		for _, parser := range parserfied {
			node := parser(ps)
			if ps.Errored() {
				ps.Pos = startpos
				return nil
			}
			if node != nil {
				nodes = append(nodes, node)
			}
		}
		return nodes
	}
}

func Any(parsers ...Parserish) Parser {
	if len(parsers) == 0 {
		return Nil
	}

	parserfied := ParsifyAll(parsers...)

	return func(ps *State) interface{} {
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
		return nil
	}
}

func Kleene(opScan Parserish, sepScan ...Parserish) Parser {
	return manyImpl(0, opScan, Never, sepScan...)
}

func KleeneUntil(opScan Parserish, untilScan Parserish, sepScan ...Parserish) Parser {
	return manyImpl(0, opScan, untilScan, sepScan...)
}

func Many(opScan Parserish, sepScan ...Parserish) Parser {
	return manyImpl(1, opScan, Never, sepScan...)
}

func ManyUntil(opScan Parserish, untilScan Parserish, sepScan ...Parserish) Parser {
	return manyImpl(1, opScan, untilScan, sepScan...)
}

func manyImpl(min int, op Parserish, until Parserish, sep ...Parserish) Parser {
	opParser := Parsify(op)
	untilParser := Parsify(until)
	sepParser := Nil
	if len(sep) > 0 {
		sepParser = Parsify(sep[0])
	}

	return func(ps *State) interface{} {
		var node interface{}
		nodes := make([]interface{}, 0, 20)
		startpos := ps.Pos
		for {
			tempPos := ps.Pos
			node = untilParser(ps)
			if !ps.Errored() {
				ps.Pos = tempPos
				if len(nodes) < min {
					ps.Pos = startpos
					ps.ErrorHere("something else")
					return nil
				}
				break
			}
			ps.ClearError()

			node = opParser(ps)
			if ps.Errored() {
				if len(nodes) < min {
					ps.Pos = startpos
					return nil
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
		return nodes
	}
}

func Maybe(parser Parserish) Parser {
	parserfied := Parsify(parser)

	return func(ps *State) interface{} {
		node := parserfied(ps)
		if ps.Errored() {
			ps.ClearError()
			return nil
		}

		return node
	}
}

func Map(parser Parserish, f func(n interface{}) interface{}) Parser {
	p := Parsify(parser)

	return func(ps *State) interface{} {
		node := p(ps)
		if ps.Errored() {
			return nil
		}
		return f(node)
	}
}

func flatten(n interface{}) interface{} {
	if s, ok := n.(string); ok {
		return s
	}

	if nodes, ok := n.([]interface{}); ok {
		sbuf := &bytes.Buffer{}
		for _, node := range nodes {
			sbuf.WriteString(flatten(node).(string))
		}
		return sbuf.String()
	}

	panic(fmt.Errorf("Dont know how to flatten %t", n))
}

func Merge(parser Parserish) Parser {
	return Map(parser, flatten)
}

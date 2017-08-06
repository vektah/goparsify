package parsec

func Nil(p Pointer) (Node, Pointer) {
	return nil, p
}

func Never(p Pointer) (Node, Pointer) {
	return Error{p.pos, "Never matches"}, p
}

func And(parsers ...Parserish) Parser {
	if len(parsers) == 0 {
		return Nil
	}

	ps := ParsifyAll(parsers...)

	return func(p Pointer) (Node, Pointer) {
		var nodes = make([]Node, 0, len(ps))
		var node Node
		newP := p
		for _, parser := range ps {
			node, newP = parser(newP)
			if node == nil {
				continue
			}
			if IsError(node) {
				return node, p
			}
			nodes = append(nodes, node)
		}
		return NewSequence(p.pos, nodes...), newP
	}
}

func Any(parsers ...Parserish) Parser {
	if len(parsers) == 0 {
		return Nil
	}

	ps := ParsifyAll(parsers...)

	return func(p Pointer) (Node, Pointer) {
		errors := []Error{}
		for _, parser := range ps {
			node, newP := parser(p)
			if err, isErr := node.(Error); isErr {
				errors = append(errors, err)
				continue
			}
			return node, newP
		}

		longestError := errors[0]
		for _, e := range errors[1:] {
			if e.pos > longestError.pos {
				longestError = e
			}
		}

		return longestError, p
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

	return func(p Pointer) (Node, Pointer) {
		var node Node
		nodes := make([]Node, 0)
		newP := p
		for {
			if node, _ := untilParser(newP); !IsError(node) {
				if len(nodes) < min {
					return NewError(newP.pos, "Unexpected input"), p
				}
				break
			}

			if node, newP = opParser(newP); IsError(node) {
				if len(nodes) < min {
					return node, p
				}
				break
			}
			nodes = append(nodes, node)
			if node, newP = sepParser(newP); IsError(node) {
				break
			}
		}
		return NewSequence(p.pos, nodes...), newP
	}
}

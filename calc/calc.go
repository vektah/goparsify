package calc

import (
	"errors"
	"fmt"

	. "github.com/vektah/goparsify"
)

var (
	value Parser

	sumOp  = Chars("+-", 1, 1)
	prodOp = Chars("/*", 1, 1)

	groupExpr = Map(Seq("(", sum, ")"), func(n Node) Node {
		return Node{Result: n.Child[1].Result}
	})

	number = Map(NumberLit(), func(n Node) Node {
		switch i := n.Result.(type) {
		case int64:
			return Node{Result: float64(i)}
		case float64:
			return Node{Result: i}
		default:
			panic(fmt.Errorf("unknown value %#v", i))
		}
	})

	sum = Map(Seq(prod, Some(Seq(sumOp, prod))), func(n Node) Node {
		i := n.Child[0].Result.(float64)

		for _, op := range n.Child[1].Child {
			switch op.Child[0].Token {
			case "+":
				i += op.Child[1].Result.(float64)
			case "-":
				i -= op.Child[1].Result.(float64)
			}
		}

		return Node{Result: i}
	})

	prod = Map(Seq(&value, Some(Seq(prodOp, &value))), func(n Node) Node {
		i := n.Child[0].Result.(float64)

		for _, op := range n.Child[1].Child {
			switch op.Child[0].Token {
			case "/":
				i /= op.Child[1].Result.(float64)
			case "*":
				i *= op.Child[1].Result.(float64)
			}
		}

		return Node{Result: i}
	})

	Y = Maybe(sum)
)

func init() {
	value = Any(number, groupExpr)
}

func Calc(input string) (float64, error) {
	result, remaining, err := ParseString(Y, input)

	if err != nil {
		return 0, err
	}

	if remaining != "" {
		return result.(float64), errors.New("left unparsed: " + remaining)
	}

	return result.(float64), nil
}

package calc

import (
	"fmt"

	. "github.com/ijt/goparsify"
)

var (
	value Parser

	sumOp  = Chars("+-", 1, 1)
	prodOp = Chars("/*", 1, 1)

	groupExpr = Seq("(", sum, ")").Map(func(n *Result) {
		n.Result = n.Child[1].Result
	})

	number = NumberLit().Map(func(n *Result) {
		switch i := n.Result.(type) {
		case int64:
			n.Result = float64(i)
		case float64:
			n.Result = i
		default:
			panic(fmt.Errorf("unknown value %#v", i))
		}
	})

	sum = Seq(prod, Some(Seq(sumOp, prod))).Map(func(n *Result) {
		i := n.Child[0].Result.(float64)

		for _, op := range n.Child[1].Child {
			switch op.Child[0].Token {
			case "+":
				i += op.Child[1].Result.(float64)
			case "-":
				i -= op.Child[1].Result.(float64)
			}
		}

		n.Result = i
	})

	prod = Seq(&value, Some(Seq(prodOp, &value))).Map(func(n *Result) {
		i := n.Child[0].Result.(float64)

		for _, op := range n.Child[1].Child {
			switch op.Child[0].Token {
			case "/":
				i /= op.Child[1].Result.(float64)
			case "*":
				i *= op.Child[1].Result.(float64)
			}
		}

		n.Result = i
	})

	y = Maybe(sum)
)

func init() {
	value = Any(number, groupExpr)
}

func calc(input string) (float64, error) {
	result, err := Run(y, input)
	if err != nil {
		return 0, err
	}

	return result.(float64), nil
}

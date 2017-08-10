package calc

import (
	"fmt"

	. "github.com/vektah/goparsify"
)

var (
	value Parser

	sumOp  = Chars("+-", 1, 1)
	prodOp = Chars("/*", 1, 1)

	groupExpr = Map(Seq("(", sum, ")"), func(n Result) Result {
		return Result{Result: n.Child[1].Result}
	})

	number = Map(NumberLit(), func(n Result) Result {
		switch i := n.Result.(type) {
		case int64:
			return Result{Result: float64(i)}
		case float64:
			return Result{Result: i}
		default:
			panic(fmt.Errorf("unknown value %#v", i))
		}
	})

	sum = Map(Seq(prod, Some(Seq(sumOp, prod))), func(n Result) Result {
		i := n.Child[0].Result.(float64)

		for _, op := range n.Child[1].Child {
			switch op.Child[0].Token {
			case "+":
				i += op.Child[1].Result.(float64)
			case "-":
				i -= op.Child[1].Result.(float64)
			}
		}

		return Result{Result: i}
	})

	prod = Map(Seq(&value, Some(Seq(prodOp, &value))), func(n Result) Result {
		i := n.Child[0].Result.(float64)

		for _, op := range n.Child[1].Child {
			switch op.Child[0].Token {
			case "/":
				i /= op.Child[1].Result.(float64)
			case "*":
				i *= op.Child[1].Result.(float64)
			}
		}

		return Result{Result: i}
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

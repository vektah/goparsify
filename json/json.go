package json

import (
	"errors"

	. "github.com/vektah/goparsify"
)

var (
	value Parser

	array = Map(And(WS, "[", Kleene(&value, And(WS, ",")), "]"), func(n Node) Node {
		return n.([]Node)[1].([]Node)
	})
	properties = Kleene(And(WS, String('"'), WS, ":", WS, &value), ",")
	object     = Map(And(WS, "{", WS, properties, WS, "}"), func(n Node) Node {
		ret := map[string]interface{}{}

		for _, prop := range n.([]Node)[1].([]Node) {
			vals := prop.([]Node)
			if len(vals) == 3 {
				ret[vals[0].(string)] = vals[2]
			} else {
				ret[vals[0].(string)] = nil
			}
		}

		return ret
	})

	_null = Map(And(WS, "null"), func(n Node) Node {
		return nil
	})

	_true = Map(And(WS, "true"), func(n Node) Node {
		return true
	})

	_false = Map(And(WS, "false"), func(n Node) Node {
		return false
	})

	Y = Map(And(&value, WS), func(n Node) Node {
		nodes := n.([]Node)
		if len(nodes) > 0 {
			return nodes[0]
		}
		return nil
	})
)

func init() {
	value = Any(_null, _true, _false, String('"'), array, object)
}

func Unmarshal(input string) (interface{}, error) {
	result, remaining, err := ParseString(Y, input)

	if err != nil {
		return result, err
	}

	if remaining != "" {
		return result, errors.New("left unparsed: " + remaining)
	}

	return result, err
}

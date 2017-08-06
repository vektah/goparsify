package json

import (
	"errors"

	. "github.com/vektah/goparsify"
)

var (
	value Parser

	array = Map(And(WS, "[", Kleene(&value, And(WS, ",")), "]"), func(n interface{}) interface{} {
		return n.([]interface{})[1].([]interface{})
	})
	properties = Kleene(And(WS, String('"'), WS, ":", WS, &value), ",")
	object     = Map(And(WS, "{", WS, properties, WS, "}"), func(n interface{}) interface{} {
		ret := map[string]interface{}{}

		for _, prop := range n.([]interface{})[1].([]interface{}) {
			vals := prop.([]interface{})
			if len(vals) == 3 {
				ret[vals[0].(string)] = vals[2]
			} else {
				ret[vals[0].(string)] = nil
			}
		}

		return ret
	})

	_null = Map(And(WS, "null"), func(n interface{}) interface{} {
		return nil
	})

	_true = Map(And(WS, "true"), func(n interface{}) interface{} {
		return true
	})

	_false = Map(And(WS, "false"), func(n interface{}) interface{} {
		return false
	})

	Y = Map(And(&value, WS), func(n interface{}) interface{} {
		nodes := n.([]interface{})
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

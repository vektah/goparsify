package json

import (
	"errors"

	. "github.com/vektah/goparsify"
)

var (
	value Parser

	_array = Map(And("[", Kleene(&value, ","), "]"), func(n Node) Node {
		ret := []interface{}{}
		for _, child := range n.Children[1].Children {
			ret = append(ret, child.Result)
		}
		return Node{Result: ret}
	})
	properties = Kleene(And(StringLit(`"`), ":", &value), ",")
	_object    = Map(And("{", properties, "}"), func(n Node) Node {
		ret := map[string]interface{}{}

		for _, prop := range n.Children[1].Children {
			ret[prop.Children[0].Token] = prop.Children[2].Result
		}

		return Node{Result: ret}
	})

	_null  = Bind("null", nil)
	_true  = Bind("true", true)
	_false = Bind("false", false)

	_string = Map(StringLit(`"`), func(n Node) Node {
		return Node{Result: n.Token}
	})
)

func init() {
	value = Any(_null, _true, _false, _string, _array, _object)
}

func Unmarshal(input string) (interface{}, error) {
	result, remaining, err := ParseString(value, input)

	if err != nil {
		return result, err
	}

	if remaining != "" {
		return result, errors.New("left unparsed: " + remaining)
	}

	return result, err
}

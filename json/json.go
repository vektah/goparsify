package json

import "errors"
import . "github.com/vektah/goparsify"

var (
	_value      Parser
	_null       = Bind("null", nil)
	_true       = Bind("true", true)
	_false      = Bind("false", false)
	_string     = StringLit(`"`)
	_number     = NumberLit()
	_properties = Kleene(And(StringLit(`"`), ":", &_value), ",")

	_array = Map(And("[", Kleene(&_value, ","), "]"), func(n Node) Node {
		ret := []interface{}{}
		for _, child := range n.Child[1].Child {
			ret = append(ret, child.Result)
		}
		return Node{Result: ret}
	})

	_object = Map(And("{", _properties, "}"), func(n Node) Node {
		ret := map[string]interface{}{}

		for _, prop := range n.Child[1].Child {
			ret[prop.Child[0].Result.(string)] = prop.Child[2].Result
		}

		return Node{Result: ret}
	})
)

func init() {
	_value = Any(_null, _true, _false, _string, _number, _array, _object)
}

func Unmarshal(input string) (interface{}, error) {
	result, remaining, err := ParseString(_value, input)

	if err != nil {
		return result, err
	}

	if remaining != "" {
		return result, errors.New("left unparsed: " + remaining)
	}

	return result, err
}

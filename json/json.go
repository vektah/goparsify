package json

import (
	. "github.com/vektah/goparsify"
)

var (
	_value      Parser
	_null       = Bind("null", nil)
	_true       = Bind("true", true)
	_false      = Bind("false", false)
	_string     = StringLit(`"`)
	_number     = NumberLit()
	_properties = Some(Seq(StringLit(`"`), ":", &_value), ",")

	_array = Seq("[", Cut(), Some(&_value, ","), "]").Map(func(n Result) Result {
		ret := []interface{}{}
		for _, child := range n.Child[2].Child {
			ret = append(ret, child.Result)
		}
		return Result{Result: ret}
	})

	_object = Seq("{", Cut(), _properties, "}").Map(func(n Result) Result {
		ret := map[string]interface{}{}

		for _, prop := range n.Child[2].Child {
			ret[prop.Child[0].Result.(string)] = prop.Child[2].Result
		}

		return Result{Result: ret}
	})
)

func init() {
	_value = Any(_null, _true, _false, _string, _number, _array, _object)
}

// Unmarshall json string into map[string]interface{} or []interface{}
func Unmarshal(input string) (interface{}, error) {
	return Run(_value, input, ASCIIWhitespace)
}

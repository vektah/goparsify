package json

import . "github.com/vektah/goparsify"

var (
	_value      Parser
	_null       = Bind("null", nil)
	_true       = Bind("true", true)
	_false      = Bind("false", false)
	_string     = StringLit(`"`)
	_number     = NumberLit()
	_properties = Some(Seq(StringLit(`"`), ":", &_value), ",")

	_array = Map(Seq("[", Some(&_value, ","), "]"), func(n Result) Result {
		ret := []interface{}{}
		for _, child := range n.Child[1].Child {
			ret = append(ret, child.Result)
		}
		return Result{Result: ret}
	})

	_object = Map(Seq("{", _properties, "}"), func(n Result) Result {
		ret := map[string]interface{}{}

		for _, prop := range n.Child[1].Child {
			ret[prop.Child[0].Result.(string)] = prop.Child[2].Result
		}

		return Result{Result: ret}
	})
)

func init() {
	_value = Any(_null, _true, _false, _string, _number, _array, _object)
}

func Unmarshal(input string) (interface{}, error) {
	return Run(_value, input)
}

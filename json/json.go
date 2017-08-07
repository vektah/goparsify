package json

import (
	"errors"

	. "github.com/vektah/goparsify"
)

var (
	value Parser

	_array = Map(And(WS, "[", Kleene(&value, And(WS, ",")), "]"), func(n *Node) *Node {
		ret := []interface{}{}
		for _, child := range n.Children[1].Children {
			ret = append(ret, child.Result)
		}
		return &Node{Result: ret}
	})
	properties = Kleene(And(WS, String('"'), WS, ":", WS, &value), ",")
	_object    = Map(And(WS, "{", WS, properties, WS, "}"), func(n *Node) *Node {
		ret := map[string]interface{}{}

		for _, prop := range n.Children[1].Children {
			ret[prop.Children[0].Token] = prop.Children[2].Result
		}

		return &Node{Result: ret}
	})

	_null = Map(And(WS, "null"), func(n *Node) *Node {
		return &Node{Result: nil}
	})

	_true = Map(And(WS, "true"), func(n *Node) *Node {
		return &Node{Result: true}
	})

	_false = Map(And(WS, "false"), func(n *Node) *Node {
		return &Node{Result: false}
	})

	_string = Map(String('"'), func(n *Node) *Node {
		return &Node{Result: n.Token}
	})

	Y = Map(And(&value, WS), func(n *Node) *Node {
		return &Node{Result: n.Children[0].Result}
	})
)

func init() {
	value = Any(_null, _true, _false, _string, _array, _object)
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

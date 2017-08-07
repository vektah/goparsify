package html

import (
	. "github.com/vektah/goparsify"
)

func Parse(input string) (result interface{}, remaining string, err error) {
	return ParseString(tag, input)
}

type Tag struct {
	Name       string
	Attributes map[string]string
	Body       []interface{}
}

var (
	tag Parser

	identifier = NoAutoWS(Merge(And(WS(), Chars("a-zA-Z", 1), Chars("a-zA-Z0-9", 0))))
	text       = Map(NotChars("<>"), func(n Node) Node {
		return Node{Result: n.Token}
	})

	element  = Any(text, &tag)
	elements = Map(Kleene(element), func(n Node) Node {
		ret := []interface{}{}
		for _, child := range n.Children {
			ret = append(ret, child.Result)
		}
		return Node{Result: ret}
	})

	attr  = And(identifier, "=", Any(String('"'), String('\'')))
	attrs = Map(Kleene(attr), func(node Node) Node {
		attr := map[string]string{}

		for _, attrNode := range node.Children {
			attr[attrNode.Children[0].Token] = attrNode.Children[2].Token
		}

		return Node{Result: attr}
	})

	tstart = And("<", identifier, attrs, ">")
	tend   = And("</", identifier, ">")
)

func init() {
	tag = Map(And(tstart, elements, tend), func(node Node) Node {
		openTag := node.Children[0]
		return Node{Result: Tag{
			Name:       openTag.Children[1].Token,
			Attributes: openTag.Children[2].Result.(map[string]string),
			Body:       node.Children[1].Result.([]interface{}),
		}}

	})
}

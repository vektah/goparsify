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

	identifier = NoAutoWS(Merge(Seq(WS(), Chars("a-zA-Z", 1), Chars("a-zA-Z0-9", 0))))
	text       = Map(NotChars("<>"), func(n Node) Node {
		return Node{Result: n.Token}
	})

	element  = Any(text, &tag)
	elements = Map(Some(element), func(n Node) Node {
		ret := []interface{}{}
		for _, child := range n.Child {
			ret = append(ret, child.Result)
		}
		return Node{Result: ret}
	})

	attr  = Seq(identifier, "=", StringLit(`"'`))
	attrs = Map(Some(attr), func(node Node) Node {
		attr := map[string]string{}

		for _, attrNode := range node.Child {
			attr[attrNode.Child[0].Token] = attrNode.Child[2].Result.(string)
		}

		return Node{Result: attr}
	})

	tstart = Seq("<", identifier, attrs, ">")
	tend   = Seq("</", identifier, ">")
)

func init() {
	tag = Map(Seq(tstart, elements, tend), func(node Node) Node {
		openTag := node.Child[0]
		return Node{Result: Tag{
			Name:       openTag.Child[1].Token,
			Attributes: openTag.Child[2].Result.(map[string]string),
			Body:       node.Child[1].Result.([]interface{}),
		}}

	})
}

package html

import (
	. "github.com/ijt/goparsify"
)

func parse(input string) (result interface{}, err error) {
	return Run(tag, input)
}

type htmlTag struct {
	Name       string
	Attributes map[string]string
	Body       []interface{}
}

var (
	tag Parser

	identifier = Regex("[a-zA-Z][a-zA-Z0-9]*")
	text       = NotChars("<>").Map(func(n *Result) { n.Result = n.Token })

	element  = Any(text, &tag)
	elements = Some(element).Map(func(n *Result) {
		ret := []interface{}{}
		for _, child := range n.Child {
			ret = append(ret, child.Result)
		}
		n.Result = ret
	})

	attr  = Seq(identifier, "=", StringLit(`"'`))
	attrs = Some(attr).Map(func(node *Result) {
		attr := map[string]string{}

		for _, attrNode := range node.Child {
			attr[attrNode.Child[0].Token] = attrNode.Child[2].Token
		}

		node.Result = attr
	})

	tstart = Seq("<", identifier, Cut(), attrs, ">")
	tend   = Seq("</", Cut(), identifier, ">")
)

func init() {
	tag = Seq(tstart, Cut(), elements, tend).Map(func(node *Result) {
		openTag := node.Child[0]
		node.Result = htmlTag{
			Name:       openTag.Child[1].Token,
			Attributes: openTag.Child[3].Result.(map[string]string),
			Body:       node.Child[2].Result.([]interface{}),
		}
	})
}

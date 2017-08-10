package html

import (
	. "github.com/vektah/goparsify"
)

func Parse(input string) (result interface{}, err error) {
	return Run(tag, input)
}

type Tag struct {
	Name       string
	Attributes map[string]string
	Body       []interface{}
}

var (
	tag Parser

	identifier = Regex("[a-zA-Z][a-zA-Z0-9]*")
	text       = Map(NotChars("<>"), func(n Result) Result {
		return Result{Result: n.Token}
	})

	element  = Any(text, &tag)
	elements = Map(Some(element), func(n Result) Result {
		ret := []interface{}{}
		for _, child := range n.Child {
			ret = append(ret, child.Result)
		}
		return Result{Result: ret}
	})

	attr  = Seq(identifier, "=", StringLit(`"'`))
	attrs = Map(Some(attr), func(node Result) Result {
		attr := map[string]string{}

		for _, attrNode := range node.Child {
			attr[attrNode.Child[0].Token] = attrNode.Child[2].Result.(string)
		}

		return Result{Result: attr}
	})

	tstart = Seq("<", identifier, Cut(), attrs, ">")
	tend   = Seq("</", Cut(), identifier, ">")
)

func init() {
	tag = Map(Seq(tstart, Cut(), elements, tend), func(node Result) Result {
		openTag := node.Child[0]
		return Result{Result: Tag{
			Name:       openTag.Child[1].Token,
			Attributes: openTag.Child[3].Result.(map[string]string),
			Body:       node.Child[2].Result.([]interface{}),
		}}

	})
}

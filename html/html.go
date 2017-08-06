package html

import . "github.com/vektah/goparsify"

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

	identifier = Merge(And(Chars("a-z", 1, 1), Chars("a-zA-Z0-9", 0)))
	text       = NotChars("<>")

	element  = Any(text, &tag)
	elements = Kleene(element)
	//attr := And(identifier, equal, String())
	attr  = And(WS, identifier, WS, "=", WS, Any(String('"'), String('\'')))
	attrs = Map(Kleene(attr, WS), func(node interface{}) interface{} {
		nodes := node.([]interface{})
		attr := map[string]string{}

		for _, attrNode := range nodes {
			attrNodes := attrNode.([]interface{})
			attr[attrNodes[0].(string)] = attrNodes[2].(string)
		}

		return attr
	})

	tstart = And("<", identifier, attrs, ">")
	tend   = And("</", identifier, ">")
)

func init() {
	tag = Map(And(tstart, elements, tend), func(node interface{}) interface{} {
		nodes := node.([]interface{})
		openTag := nodes[0].([]interface{})
		return Tag{
			Name:       openTag[1].(string),
			Attributes: openTag[2].(map[string]string),
			Body:       nodes[1].([]interface{}),
		}

	})
}

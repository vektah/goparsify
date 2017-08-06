package main

import (
	"fmt"

	. "github.com/vektah/goparsify"
)

func html(p Pointer) (Node, Pointer) {
	opentag := Exact("<")
	closetag := Exact(">")
	equal := Exact("=")
	slash := Exact("/")
	identifier := And(Char(Range("a-z")), CharRun(Range("a-zA-Z0-9")))
	text := CharRunUntil("<>")

	var tag Parser

	element := Any(text, &tag)
	elements := Kleene(element)
	//attr := And(identifier, equal, String())
	attr := And(identifier, equal, Exact(`"test"`))
	attrws := And(attr, WS)
	attrs := Kleene(attrws)
	tstart := And(opentag, identifier, attrs, closetag)
	tend := And(opentag, slash, identifier, closetag)
	tag = And(tstart, elements, tend)

	return element(p)
}

func main() {
	node, _ := html(Input("<h1>hello world</h1>"))
	fmt.Printf("%#v\n", node)
}

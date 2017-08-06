package main

import (
	"fmt"

	. "github.com/vektah/goparsify"
)

func html(p Pointer) (Node, Pointer) {
	identifier := And(Range("a-z", 1, 1), Range("a-zA-Z0-9"))
	text := CharRunUntil("<>")

	var tag Parser

	element := Any(text, &tag)
	elements := Kleene(element)
	//attr := And(identifier, equal, String())
	attr := And(identifier, "=", `"test"`)
	attrws := And(attr, WS)
	attrs := Kleene(attrws)
	tstart := And("<", identifier, attrs, ">")
	tend := And("</", identifier, ">")
	tag = And(tstart, elements, tend)

	return element(p)
}

func main() {
	result, _, err := ParseString(html, "<h1>hello world</h1>")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", result)
}

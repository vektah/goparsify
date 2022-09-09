goparsify [![CircleCI](https://circleci.com/gh/ijt/goparsify/tree/master.svg?style=shield)](https://circleci.com/gh/ijt/goparsify/tree/master) [![godoc](http://b.repl.ca/v1/godoc-reference-blue.png)](https://godoc.org/github.com/ijt/goparsify) [![Go Report Card](https://goreportcard.com/badge/github.com/ijt/goparsify)](https://goreportcard.com/report/github.com/ijt/goparsify)
=========

This is a fork of github.com/vektah/goparsify. The fork brings it up to date
with some current Go practices (modules) and works to improve error messages.
The original readme follows:

A parser-combinator library for building easy to test, read and maintain parsers using functional composition.

Everything should be unicode safe by default, but you can opt out of unicode whitespace for a decent ~20% performance boost.
```go
Run(parser, input, ASCIIWhitespace)
```

### benchmarks
I dont have many benchmarks set up yet, its pretty quick:
```
$ go test -benchmem -bench=. ./json
BenchmarkUnmarshalParsec-8         20000             74880 ns/op           50846 B/op       1318 allocs/op
BenchmarkUnmarshalParsify-8        30000             50631 ns/op           45055 B/op        233 allocs/op
BenchmarkUnmarshalStdlib-8         30000             46989 ns/op           14210 B/op        260 allocs/op
PASS
ok      github.com/ijt/goparsify/json        6.124s
```

Most of the remaining small allocs are from putting things in `interface{}` and are pretty unavoidable. https://www.darkcoding.net/software/go-the-price-of-interface/ is a good read.

### debugging parsers

When a parser isnt working as you intended you can build with debugging and enable logging to get a detailed log of exactly what the parser is doing.

1. First build with debug using `-tags debug`
2. enable logging by calling `EnableLogging(os.Stdout)` in your code

This works great with tests, eg in the goparsify source tree
```
adam:goparsify(master)$ go test -tags debug ./html -v
=== RUN   TestParse
html.go:48 | <body>hello <p  | tag {
html.go:43 | <body>hello <p  |   tstart {
html.go:43 | body>hello <p c |     < found <
html.go:20 | >hello <p color |     identifier found body
html.go:33 | >hello <p color |     attrs {
html.go:32 | >hello <p color |       attr {
html.go:20 | >hello <p color |         identifier did not find [a-zA-Z][a-zA-Z0-9]*
html.go:32 | >hello <p color |       } did not find [a-zA-Z][a-zA-Z0-9]*
html.go:33 | >hello <p color |     } found
html.go:43 | hello <p color= |     > found >
html.go:43 | hello <p color= |   } found [<,body,,map[string]string{},>]
html.go:24 | hello <p color= |   elements {
html.go:23 | hello <p color= |     element {
html.go:21 | <p color="blue" |       text found hello
html.go:23 | <p color="blue" |     } found "hello "
html.go:23 | <p color="blue" |     element {
html.go:21 | <p color="blue" |       text did not find <>
html.go:48 | <p color="blue" |       tag {
html.go:43 | <p color="blue" |         tstart {
html.go:43 | p color="blue"> |           < found <
html.go:20 |  color="blue">w |           identifier found p
html.go:33 |  color="blue">w |           attrs {
html.go:32 |  color="blue">w |             attr {
html.go:20 | ="blue">world</ |               identifier found color
html.go:32 | "blue">world</p |               = found =
html.go:32 | >world</p></bod |               string literal found "blue"
html.go:32 | >world</p></bod |             } found [color,=,"blue"]
html.go:32 | >world</p></bod |             attr {
html.go:20 | >world</p></bod |               identifier did not find [a-zA-Z][a-zA-Z0-9]*
html.go:32 | >world</p></bod |             } did not find [a-zA-Z][a-zA-Z0-9]*
html.go:33 | >world</p></bod |           } found [[color,=,"blue"]]
html.go:43 | world</p></body |           > found >
html.go:43 | world</p></body |         } found [<,p,,map[string]string{"color":"blue"},>]
html.go:24 | world</p></body |         elements {
html.go:23 | world</p></body |           element {
html.go:21 | </p></body>     |             text found world
html.go:23 | </p></body>     |           } found "world"
html.go:23 | </p></body>     |           element {
html.go:21 | </p></body>     |             text did not find <>
html.go:48 | </p></body>     |             tag {
html.go:43 | </p></body>     |               tstart {
html.go:43 | /p></body>      |                 < found <
html.go:20 | /p></body>      |                 identifier did not find [a-zA-Z][a-zA-Z0-9]*
html.go:43 | </p></body>     |               } did not find [a-zA-Z][a-zA-Z0-9]*
html.go:48 | </p></body>     |             } did not find [a-zA-Z][a-zA-Z0-9]*
html.go:23 | </p></body>     |           } did not find [a-zA-Z][a-zA-Z0-9]*
html.go:24 | </p></body>     |         } found ["world"]
html.go:44 | </p></body>     |         tend {
html.go:44 | p></body>       |           </ found </
html.go:20 | ></body>        |           identifier found p
html.go:44 | </body>         |           > found >
html.go:44 | </body>         |         } found [</,,p,>]
html.go:48 | </body>         |       } found "hello "
html.go:23 | </body>         |     } found html.htmlTag{Name:"p", Attributes:map[string]string{"color":"blue"}, Body:[]interface {}{"world"}}
html.go:23 | </body>         |     element {
html.go:48 | </body>         |       tag {
html.go:43 | </body>         |         tstart {
html.go:43 | /body>          |           < found <
html.go:20 | /body>          |           identifier did not find [a-zA-Z][a-zA-Z0-9]*
html.go:43 | </body>         |         } did not find [a-zA-Z][a-zA-Z0-9]*
html.go:48 | </body>         |       } did not find [a-zA-Z][a-zA-Z0-9]*
html.go:21 | </body>         |       text did not find <>
html.go:23 | </body>         |     } did not find [a-zA-Z][a-zA-Z0-9]*
html.go:24 | </body>         |   } found ["hello ",html.htmlTag{Name:"p", Attributes:map[string]string{"color":"blue"}, Body:[]interface {}{"world"}}]
html.go:44 | </body>         |   tend {
html.go:44 | body>           |     </ found </
html.go:20 | >               |     identifier found body
html.go:44 |                 |     > found >
html.go:44 |                 |   } found [</,,body,>]
html.go:48 |                 | } found [[<,body,,map[string]string{},>],,[]interface {}{"hello ", html.htmlTag{Name:"p", Attributes:map[string]string{"color":"blue"}, Body:[]interface {}{"world"}}},[</,,body,>]]
--- PASS: TestParse (0.00s)
PASS
ok      github.com/ijt/goparsify/html        0.117s
```

### debugging performance
If you build the parser with -tags debug it will instrument each parser and a call to DumpDebugStats() will show stats:

|             var name |              matches |      total time |       self time |      calls |     errors | location
| -------------------- | -------------------- | --------------- | --------------- | ---------- | ---------- | ----------
|               _value |                Any() |      5.0685431s |       34.0131ms |     878801 |          0 | json.go:36
|              _object |                Seq() |      3.7513821s |       10.5038ms |     161616 |      40403 | json.go:24
|          _properties |               Some() |      3.6863512s |        5.5028ms |     121213 |          0 | json.go:14
|          _properties |                Seq() |      3.4912614s |       46.0229ms |     818185 |          0 | json.go:14
|               _array |                Seq() |      931.4679ms |        3.5014ms |      65660 |      55558 | json.go:16
|               _array |               Some() |      911.4597ms |              0s |      10102 |          0 | json.go:16
|          _properties |       string literal |      126.0662ms |       44.5201ms |     818185 |          0 | json.go:14
|              _string |       string literal |        67.033ms |       26.0126ms |     671723 |     136369 | json.go:12
|          _properties |                    : |       50.0238ms |       45.0205ms |     818185 |          0 | json.go:14
|          _properties |                    , |       48.5189ms |       36.0146ms |     818185 |     121213 | json.go:14
|              _number |       number literal |       28.5159ms |       10.5062ms |     287886 |     106066 | json.go:13
|                _true |                 true |       17.5086ms |       12.5069ms |     252537 |     232332 | json.go:10
|                _null |                 null |       14.5082ms |        11.007ms |     252538 |     252535 | json.go:9
|              _object |                    } |       10.5051ms |       10.5033ms |     121213 |          0 | json.go:24
|               _false |                false |       10.5049ms |        5.0019ms |     232333 |     222229 | json.go:11
|              _object |                    { |       10.0046ms |        5.0052ms |     161616 |      40403 | json.go:24
|               _array |                    , |        4.5024ms |        4.0018ms |      50509 |      10102 | json.go:16
|               _array |                    [ |        4.5014ms |        2.0006ms |      65660 |      55558 | json.go:16
|               _array |                    ] |              0s |              0s |      10102 |          0 | json.go:16

All times are cumulative, it would be nice to break this down into a parse tree with relative times. This is a nice addition to pprof as it will break down the parsers based on where they are used instead of grouping them all by type.

This is **free** when the debug tag isnt used.

### example calculator
Lets say we wanted to build a calculator that could take an expression and calculate the result.

Lets start with test:
```go
func TestNumbers(t *testing.T) {
	result, err := Calc(`1`)
	require.NoError(t, err)
	require.EqualValues(t, 1, result)
}
```

Then define a parser for numbers
```go
var number = NumberLit().Map(func(n Result) Result {
    switch i := n.Result.(type) {
    case int64:
        return Result{Result: float64(i)}
    case float64:
        return Result{Result: i}
    default:
        panic(fmt.Errorf("unknown value %#v", i))
    }
})

func Calc(input string) (float64, error) {
	result, err := Run(y, input)
	if err != nil {
		return 0, err
	}

	return result.(float64), nil
}

```

This parser will return numbers either as float64 or int depending on the literal, for this calculator we only want floats so we Map the results and type cast.

Run the tests and make sure everything is ok.

Time to add addition

```go
func TestAddition(t *testing.T) {
	result, err := Calc(`1+1`)
	require.NoError(t, err)
	require.EqualValues(t, 2, result)
}


var sumOp  = Chars("+-", 1, 1)

sum = Seq(number, Some(And(sumOp, number))).Map(func(n Result) Result {
    i := n.Child[0].Result.(float64)

    for _, op := range n.Child[1].Child {
        switch op.Child[0].Token {
        case "+":
            i += op.Child[1].Result.(float64)
        case "-":
            i -= op.Child[1].Result.(float64)
        }
    }

    return Result{Result: i}
})

// and update Calc to point to the new root parser -> `result, err := ParseString(sum, input)`
```

This parser will match number ([+-] number)+, then map its to be the sum. See how the Child map directly to the positions in the parsers? n is the result of the and, `n.Child[0]` is its first argument, `n.Child[1]` is the result of the Some parser, `n.Child[1].Child[0]` is the result of the first And and so fourth. Given how closely tied the parser and the Map are it is good to keep the two together.

You can continue like this and add multiplication and parenthesis fairly easily. Eventually if you keep adding parsers you will end up with a loop, and go will give you a handy error message like:
```
typechecking loop involving value = goparsify.Any(number, groupExpr)
```

we need to break the loop using a pointer, then set its value in init
```go
var (
    value Parser
    prod = Seq(&value, Some(And(prodOp, &value)))
)

func init() {
	value = Any(number, groupExpr)
}
```

Take a look at [calc](calc/calc.go) for a full example.

### preventing backtracking with cuts
A cut is a marker that prevents backtracking past the point it was set. This greatly improves error messages when used correctly:
```go
alpha := Chars("a-z")

// without a cut if the close tag is left out the parser will backtrack and ignore the rest of the string
nocut := Many(Any(Seq("<", alpha, ">"), alpha))
_, err := Run(nocut, "asdf <foo")
fmt.Println(err.Error())
// Outputs: left unparsed: <foo

// with a cut, once we see the open tag we know there must be a close tag that matches it, so the parser will error
cut := Many(Any(Seq("<", Cut(), alpha, ">"), alpha))
_, err = Run(cut, "asdf <foo")
fmt.Println(err.Error())
// Outputs: offset 9: expected >
```

### prior art

Inspired by https://github.com/prataprc/goparsec

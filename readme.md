goparsify [![CircleCI](https://circleci.com/gh/Vektah/goparsify/tree/master.svg?style=shield)](https://circleci.com/gh/Vektah/goparsify/tree/master) [![godoc](http://b.repl.ca/v1/godoc-reference-blue.png)](https://godoc.org/github.com/Vektah/goparsify) [![Go Report Card](https://goreportcard.com/badge/github.com/vektah/goparsify)](https://goreportcard.com/report/github.com/vektah/goparsify)
=========

A parser-combinator library for building easy to test, read and maintain parsers using functional composition.

Everything should be unicode safe by default, but you can opt out of unicode whitespace for a decent ~20% performance boost.
```go
Run(parser, input, ASCIIWhitespace)
```

### benchmarks
I dont have many benchmarks set up yet, but the json parser is very promising. Nearly keeping up with the stdlib for raw speed:
```
$ go test -bench=. -benchtime=2s -benchmem ./json
BenchmarkUnmarshalParsec-8         20000             65682 ns/op           50460 B/op       1318 allocs/op
BenchmarkUnmarshalParsify-8        30000             51292 ns/op           45104 B/op        334 allocs/op
BenchmarkUnmarshalStdlib-8         30000             46522 ns/op           13953 B/op        262 allocs/op
PASS
ok      github.com/vektah/goparsify/json        10.840s
```

### debugging parsers

When a parser isnt working as you intended you can build with debugging and enable logging to get a detailed log of exactly what the parser is doing.

1. First build with debug using `-tags debug`
2. enable logging by calling `EnableLogging(os.Stdout)` in your code

This works great with tests, eg in the goparsify source tree
```
$ cd html
$ go test -tags debug -parselogs
html.go:50 | <body>hello <p  |            | tag
html.go:45 | <body>hello <p  |            |   tstart
html.go:45 | body>hello <p c | <          |     <
html.go:20 | >hello <p color | body       |     identifier
html.go:35 | >hello <p color |            |     attrs
html.go:34 | >hello <p color |            |       attr
html.go:20 | >hello <p color | fail       |         identifier
html.go:45 | hello <p color= | >          |     >
html.go:26 | hello <p color= |            |   elements
html.go:25 | hello <p color= |            |     element
html.go:21 | <p color="blue" | hello      |       text
html.go:25 | <p color="blue" |            |     element
html.go:21 | <p color="blue" | fail       |       text
html.go:50 | <p color="blue" |            |       tag
html.go:45 | <p color="blue" |            |         tstart
html.go:45 | p color="blue"> | <          |           <
html.go:20 |  color="blue">w | p          |           identifier
html.go:35 |  color="blue">w |            |           attrs
html.go:34 |  color="blue">w |            |             attr
html.go:20 | ="blue">world</ | color      |               identifier
html.go:34 | "blue">world</p | =          |               =
html.go:34 | >world</p></bod |            |               string literal
html.go:34 | >world</p></bod |            |             attr
html.go:20 | >world</p></bod | fail       |               identifier
html.go:45 | world</p></body | >          |           >
html.go:26 | world</p></body |            |         elements
html.go:25 | world</p></body |            |           element
html.go:21 | </p></body>     | world      |             text
html.go:25 | </p></body>     |            |           element
html.go:21 | </p></body>     | fail       |             text
html.go:50 | </p></body>     |            |             tag
html.go:45 | </p></body>     |            |               tstart
html.go:45 | /p></body>      | <          |                 <
html.go:20 | /p></body>      | fail       |                 identifier
html.go:46 | </p></body>     |            |         tend
html.go:46 | p></body>       | </         |           </
html.go:20 | ></body>        | p          |           identifier
html.go:46 | </body>         | >          |           >
html.go:25 | </body>         |            |     element
html.go:21 | </body>         | fail       |       text
html.go:50 | </body>         |            |       tag
html.go:45 | </body>         |            |         tstart
html.go:45 | /body>          | <          |           <
html.go:20 | /body>          | fail       |           identifier
html.go:46 | </body>         |            |   tend
html.go:46 | body>           | </         |     </
html.go:20 | >               | body       |     identifier
html.go:46 |                 | >          |     >
PASS
ok      github.com/vektah/goparsify/html        0.118s
```

### debugging performance
If you build the parser with -tags debug it will instrument each parser and a call to DumpDebugStats() will show stats:
```
              _value    12.6186996s        2618801      calls   json.go:36
             _object    9.0349494s          361213      calls   json.go:24
         _properties    8.9393995s          361213      calls   json.go:14
         _properties    8.5702176s         2438185      calls   json.go:14
              _array    2.3471541s          391315      calls   json.go:16
              _array     2.263117s           30102      calls   json.go:16
         _properties    257.1277ms         2438185      calls   json.go:14
             _string    165.0818ms         2528489      calls   json.go:12
         _properties     94.0519ms         2438185      calls   json.go:14
               _true     81.5423ms         2618798      calls   json.go:10
              _false      74.032ms         2558593      calls   json.go:11
               _null     70.0318ms         2618801      calls   json.go:9
         _properties     56.5289ms         2438185      calls   json.go:14
             _number     52.0277ms          933135      calls   json.go:13
              _array      20.008ms          391315      calls   json.go:16
             _object     17.5049ms          361213      calls   json.go:24
             _object      9.0073ms          361213      calls   json.go:24
              _array      3.0025ms          150509      calls   json.go:16
              _array      3.0019ms           30102      calls   json.go:16
```
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
var number = Map(NumberLit(), func(n Result) Result {
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

sum = Map(Seq(number, Some(And(sumOp, number))), func(n Result) Result {
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

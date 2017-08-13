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
      var name           matches               total time        self time           calls           errors      location
              _value                 Any()       6.3303551s        46.0214ms      878801 calls           0 errors json.go:36
             _string        string literal       100.0559ms         44.019ms      848489 calls      313135 errors json.go:12
              _false                 false        52.0288ms        43.0197ms      858593 calls      848489 errors json.go:11
               _null                  null        58.0309ms        42.0222ms      878801 calls      878798 errors json.go:9
         _properties        string literal       119.3651ms        42.0151ms      818185 calls           0 errors json.go:14
         _properties                     :        54.5277ms         41.018ms      818185 calls           0 errors json.go:14
               _true                  true        56.5292ms        37.0166ms      878798 calls      858593 errors json.go:10
         _properties                 Seq()       4.2989722s        35.5217ms      818185 calls           0 errors json.go:14
         _properties                     ,        45.0263ms         35.519ms      818185 calls      121213 errors json.go:14
             _number        number literal        30.0208ms        11.5093ms      313135 calls      131315 errors json.go:13
              _array                     [        12.0045ms         10.504ms      131315 calls      121213 errors json.go:16
         _properties                Some()       4.4800665s         9.0051ms      121213 calls           0 errors json.go:14
             _object                     {        11.0053ms         8.5041ms      121213 calls           0 errors json.go:24
             _object                     }         9.0022ms         8.0031ms      121213 calls           0 errors json.go:24
             _object                 Seq()       4.5375994s         6.5055ms      121213 calls           0 errors json.go:24
              _array                 Seq()       1.1524115s         5.5023ms      131315 calls      121213 errors json.go:16
              _array                     ,         3.0008ms         4.0012ms       50509 calls       10102 errors json.go:16
              _array                     ]         1.5013ms         1.5011ms       10102 calls           0 errors json.go:16
              _array                Some()        1.116393s               0s       10102 calls           0 errors json.go:16
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

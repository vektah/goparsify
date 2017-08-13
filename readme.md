goparsify [![CircleCI](https://circleci.com/gh/Vektah/goparsify/tree/master.svg?style=shield)](https://circleci.com/gh/Vektah/goparsify/tree/master) [![godoc](http://b.repl.ca/v1/godoc-reference-blue.png)](https://godoc.org/github.com/Vektah/goparsify) [![Go Report Card](https://goreportcard.com/badge/github.com/vektah/goparsify)](https://goreportcard.com/report/github.com/vektah/goparsify)
=========

A parser-combinator library for building easy to test, read and maintain parsers using functional composition.

Everything should be unicode safe by default, but you can opt out of unicode whitespace for a decent ~20% performance boost.
```go
Run(parser, input, ASCIIWhitespace)
```

### benchmarks
I dont have many benchmarks set up yet, but the json parser keeps up with the stdlib for raw speed:
```
$ go test -bench=. -benchtime=2s -benchmem ./json
BenchmarkUnmarshalParsec-8         50000             66012 ns/op           50462 B/op       1318 allocs/op
BenchmarkUnmarshalParsify-8       100000             46713 ns/op           44543 B/op        332 allocs/op
BenchmarkUnmarshalStdlib-8        100000             46967 ns/op           13952 B/op        262 allocs/op
PASS
ok      github.com/vektah/goparsify/json        14.424s
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
              _value                 Any()       5.1725682s        57.0243ms      878801 calls           0 errors json.go:36
         _properties        string literal       131.5662ms        45.0273ms      818185 calls           0 errors json.go:14
         _properties                 Seq()        3.579274s         42.016ms      818185 calls           0 errors json.go:14
         _properties                     ,        50.5254ms        35.5182ms      818185 calls      121213 errors json.go:14
         _properties                     :        51.5256ms        35.0183ms      818185 calls           0 errors json.go:14
             _string        string literal        78.0462ms        28.0172ms      671723 calls      136369 errors json.go:12
             _number        number literal        34.5187ms        15.5065ms      287886 calls      106066 errors json.go:13
               _null                  null         17.011ms         8.5058ms      252538 calls      252535 errors json.go:9
         _properties                Some()       3.7588588s         7.5023ms      121213 calls           0 errors json.go:14
             _object                     {        10.5049ms         7.0029ms      161616 calls       40403 errors json.go:24
               _true                  true        10.0072ms          6.505ms      252537 calls      232332 errors json.go:10
              _false                 false         9.0039ms         4.5032ms      232333 calls      222229 errors json.go:11
             _object                 Seq()         3.81739s         4.5016ms      161616 calls       40403 errors json.go:24
              _array                     [         5.0013ms         4.0011ms       65660 calls       55558 errors json.go:16
             _object                     }         5.5023ms         2.5021ms      121213 calls           0 errors json.go:24
              _array                     ,         2.0018ms         1.5026ms       50509 calls       10102 errors json.go:16
              _array                Some()       933.4591ms          500.8µs       10102 calls           0 errors json.go:16
              _array                 Seq()       952.9664ms               0s       65660 calls       55558 errors json.go:16
              _array                     ]               0s               0s       10102 calls           0 errors json.go:16
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

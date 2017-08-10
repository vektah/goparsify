goparsify [![CircleCI](https://circleci.com/gh/Vektah/goparsify/tree/master.svg?style=shield)](https://circleci.com/gh/Vektah/goparsify/tree/master) [![godoc](http://b.repl.ca/v1/godoc-reference-blue.png)](https://godoc.org/github.com/Vektah/goparsify)
=========

A parser-combinator library for building easy to test, read and maintain parsers using functional composition.

### benchmarks
I dont have many benchmarks set up yet, but the json parser is very promising. Nearly keeping up with the stdlib for raw speed:
```
$ go test -bench=. -benchtime=2s -benchmem ./json
BenchmarkUnmarshalParsec-8         50000             71447 ns/op           50464 B/op       1318 allocs/op
BenchmarkUnmarshalParsify-8        50000             56414 ns/op           43887 B/op        334 allocs/op
BenchmarkUnmarshalStdlib-8         50000             50187 ns/op           13949 B/op        262 allocs/op
PASS
ok      github.com/vektah/goparsify/json        10.840s
```

### debugging parsers

When a parser isnt working as you intended you can build with debugging and enable logging to get a detailed log of exactly what the parser is doing.

1. First build with debug using `-tags debug`
2. enable logging by passing a runtime flag -parselogs or calling `EnableLogging(os.Stdout)` in your code.

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
               Any()    415.7136ms           87000      calls  json.go:35
               Map()    309.6569ms           12000      calls  json.go:31
               Seq()    298.6519ms           12000      calls  json.go:23
              Some()    290.6462ms           12000      calls  json.go:13
               Seq()    272.6392ms           81000      calls  json.go:13
               Seq()     78.0404ms           13000      calls  json.go:15
               Map()     78.0404ms           13000      calls  json.go:21
              Some()     77.0401ms            1000      calls  json.go:15
      string literal      7.5053ms           81000      calls  json.go:13
      string literal      4.5031ms           84000      calls  json.go:11
                   ,      4.0008ms           81000      calls  json.go:13
               false      2.0018ms           85000      calls  json.go:10
                null      2.0005ms           87000      calls  json.go:8
                true       1.501ms           87000      calls  json.go:9
                   :       500.8µs           81000      calls  json.go:13
                   [            0s           13000      calls  json.go:15
                   }            0s           12000      calls  json.go:23
                   {            0s           12000      calls  json.go:23
      number literal            0s           31000      calls  json.go:12
                   ]            0s            1000      calls  json.go:15
                 Nil            0s               0      calls  profile/json.go:148
                   ,            0s            5000      calls  json.go:15
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


### prior art

Inspired by https://github.com/prataprc/goparsec

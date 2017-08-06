package parsec

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNil(t *testing.T) {
	p := Pointer{"hello world", 0}

	node, p2 := Nil(p)
	require.Equal(t, nil, node)
	require.Equal(t, p, p2)
}

func TestAnd(t *testing.T) {
	p := Pointer{"hello world", 0}

	t.Run("matches sequence", func(t *testing.T) {
		node, p2 := And("hello", WS, "world")(p)
		require.Equal(t, []Node{"hello", "world"}, node)
		require.Equal(t, "", p2.Get())
	})

	t.Run("returns errors", func(t *testing.T) {
		e, p3 := And("hello", WS, "there")(p)
		require.Equal(t, NewError(6, "Expected there"), e)
		require.Equal(t, 0, p3.pos)
	})

	t.Run("No parsers", func(t *testing.T) {
		assertNilParser(t, And())
	})
}

func TestMaybe(t *testing.T) {
	p := Pointer{"hello world", 0}

	t.Run("matches sequence", func(t *testing.T) {
		node, p2 := Maybe("hello")(p)
		require.Equal(t, "hello", node)
		require.Equal(t, " world", p2.Get())
	})

	t.Run("returns no errors", func(t *testing.T) {
		e, p3 := Maybe("world")(p)
		require.Equal(t, nil, e)
		require.Equal(t, 0, p3.pos)
	})
}

func TestAny(t *testing.T) {
	p := Pointer{"hello world!", 0}

	t.Run("Matches any", func(t *testing.T) {
		node, p2 := Any("hello", "world")(p)
		require.Equal(t, "hello", node)
		require.Equal(t, 5, p2.pos)
	})

	t.Run("Returns longest error", func(t *testing.T) {
		err, p2 := Any(
			Exact("nope"),
			And(Exact("hello"), WS, Exact("world"), Exact(".")),
			And(Exact("hello"), WS, Exact("brother")),
		)(p)
		require.Equal(t, NewError(11, "Expected ."), err)
		require.Equal(t, 0, p2.pos)
	})

	t.Run("Accepts nil matches", func(t *testing.T) {
		node, p2 := Any(Exact("ffffff"), WS)(p)
		require.Equal(t, nil, node)
		require.Equal(t, 0, p2.pos)
	})

	t.Run("No parsers", func(t *testing.T) {
		assertNilParser(t, Any())
	})
}

func TestKleene(t *testing.T) {
	p := Pointer{"a,b,c,d,e,", 0}

	t.Run("Matches sequence with sep", func(t *testing.T) {
		node, p2 := Kleene(CharRun("abcdefg"), Exact(","))(p)
		require.Equal(t, []Node{"a", "b", "c", "d", "e"}, node)
		require.Equal(t, 10, p2.pos)
	})

	t.Run("Matches sequence without sep", func(t *testing.T) {
		node, p2 := Kleene(Any(CharRun("abcdefg"), Exact(",")))(p)
		require.Equal(t, []Node{"a", ",", "b", ",", "c", ",", "d", ",", "e", ","}, node)
		require.Equal(t, 10, p2.pos)
	})

	t.Run("Stops on error", func(t *testing.T) {
		node, p2 := Kleene(CharRun("abc"), Exact(","))(p)
		require.Equal(t, []Node{"a", "b", "c"}, node)
		require.Equal(t, 6, p2.pos)
		require.Equal(t, "d,e,", p2.Get())
	})
}

func TestMany(t *testing.T) {
	p := Pointer{"a,b,c,d,e,", 0}

	t.Run("Matches sequence with sep", func(t *testing.T) {
		node, p2 := Many(CharRun("abcdefg"), Exact(","))(p)
		require.Equal(t, []Node{"a", "b", "c", "d", "e"}, node)
		require.Equal(t, 10, p2.pos)
	})

	t.Run("Matches sequence without sep", func(t *testing.T) {
		node, p2 := Many(Any(CharRun("abcdefg"), Exact(",")))(p)
		require.Equal(t, []Node{"a", ",", "b", ",", "c", ",", "d", ",", "e", ","}, node)
		require.Equal(t, 10, p2.pos)
	})

	t.Run("Stops on error", func(t *testing.T) {
		node, p2 := Many(CharRun("abc"), Exact(","))(p)
		require.Equal(t, []Node{"a", "b", "c"}, node)
		require.Equal(t, 6, p2.pos)
		require.Equal(t, "d,e,", p2.Get())
	})

	t.Run("Returns error if nothing matches", func(t *testing.T) {
		node, p2 := Many(CharRun("def"), Exact(","))(p)
		require.Equal(t, NewError(0, "Expected some of def"), node)
		require.Equal(t, 0, p2.pos)
		require.Equal(t, "a,b,c,d,e,", p2.Get())
	})
}

func TestKleeneUntil(t *testing.T) {
	p := Pointer{"a,b,c,d,e,fg", 0}

	t.Run("Matches sequence with sep", func(t *testing.T) {
		node, p2 := KleeneUntil(CharRun("abcde"), CharRun("d"), Exact(","))(p)
		require.Equal(t, []Node{"a", "b", "c"}, node)
		require.Equal(t, 6, p2.pos)
	})

	t.Run("Breaks if separator does not match", func(t *testing.T) {
		node, p2 := KleeneUntil(Char("abcdefg"), Char("y"), Exact(","))(p)
		require.Equal(t, []Node{"a", "b", "c", "d", "e", "f"}, node)
		require.Equal(t, 11, p2.pos)
	})
}

func TestManyUntil(t *testing.T) {
	p := Pointer{"a,b,c,d,e,", 0}

	t.Run("Matches sequence until", func(t *testing.T) {
		node, p2 := ManyUntil(CharRun("abcdefg"), Char("d"), Exact(","))(p)
		require.Equal(t, []Node{"a", "b", "c"}, node)
		require.Equal(t, 6, p2.pos)
	})

	t.Run("Returns error until matches early", func(t *testing.T) {
		node, p2 := ManyUntil(CharRun("abc"), Exact("a"), Exact(","))(p)
		require.Equal(t, NewError(0, "Unexpected input"), node)
		require.Equal(t, 0, p2.pos)
		require.Equal(t, "a,b,c,d,e,", p2.Get())
	})
}

type htmlTag struct {
	Name string
}

func TestMap(t *testing.T) {
	parser := Map(And("<", Range("a-zA-Z0-9"), ">"), func(n Node) Node {
		return htmlTag{n.([]Node)[1].(string)}
	})

	t.Run("sucess", func(t *testing.T) {
		result, _ := parser(Pointer{"<html>", 0})
		require.Equal(t, htmlTag{"html"}, result)
	})

	t.Run("error", func(t *testing.T) {
		result, ptr := parser(Pointer{"<html", 0})
		require.Equal(t, NewError(5, "Expected >"), result)
		require.Equal(t, 0, ptr.pos)
	})
}

func TestMerge(t *testing.T) {
	var bracer Parser
	bracer = And("(", Maybe(&bracer), ")")
	parser := Merge(bracer)

	t.Run("sucess", func(t *testing.T) {
		result, _ := parser(Pointer{"((()))", 0})
		require.Equal(t, "((()))", result)
	})

	t.Run("error", func(t *testing.T) {
		result, ptr := parser(Pointer{"((())", 0})
		require.Equal(t, NewError(5, "Expected )"), result)
		require.Equal(t, 0, ptr.pos)
	})

	require.Panics(t, func() {
		flatten(1)
	})
}

func assertNilParser(t *testing.T, parser Parser) {
	p := Pointer{"fff", 0}
	node, p2 := parser(p)
	require.Equal(t, nil, node)
	require.Equal(t, p, p2)
}

package goparsify

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNil(t *testing.T) {
	node, p2 := runParser("hello world", Nil)

	require.Nil(t, node)
	require.Equal(t, 0, p2.Pos)
	require.False(t, p2.Errored())
}

func TestNever(t *testing.T) {
	node, p2 := runParser("hello world", Never)

	require.Nil(t, node)
	require.Equal(t, 0, p2.Pos)
	require.True(t, p2.Errored())
}

func TestAnd(t *testing.T) {
	parser := And("hello", WS, "world")

	t.Run("matches sequence", func(t *testing.T) {
		node, p2 := runParser("hello world", parser)
		assertSequence(t, node, "hello", "world")
		require.Equal(t, "", p2.Get())
	})

	t.Run("returns errors", func(t *testing.T) {
		_, p2 := runParser("hello there", parser)
		require.Equal(t, "world", p2.Error.Expected)
		require.Equal(t, 6, p2.Error.pos)
		require.Equal(t, 0, p2.Pos)
	})

	t.Run("No parsers", func(t *testing.T) {
		assertNilParser(t, And())
	})
}

func TestMaybe(t *testing.T) {
	t.Run("matches sequence", func(t *testing.T) {
		node, p2 := runParser("hello world", Maybe("hello"))
		require.Equal(t, "hello", node.Token)
		require.Equal(t, " world", p2.Get())
	})

	t.Run("returns no errors", func(t *testing.T) {
		node, p3 := runParser("hello world", Maybe("world"))
		require.Nil(t, node)
		require.False(t, p3.Errored())
		require.Equal(t, 0, p3.Pos)
	})
}

func TestAny(t *testing.T) {
	t.Run("Matches any", func(t *testing.T) {
		node, p2 := runParser("hello world!", Any("hello", "world"))
		require.Equal(t, "hello", node.Token)
		require.Equal(t, 5, p2.Pos)
	})

	t.Run("Returns longest error", func(t *testing.T) {
		_, p2 := runParser("hello world!", Any(
			"nope",
			And("hello", WS, "world", "."),
			And("hello", WS, "brother"),
		))
		require.Equal(t, "offset 11: Expected .", p2.Error.Error())
		require.Equal(t, 11, p2.Error.Pos())
		require.Equal(t, 0, p2.Pos)
	})

	t.Run("Accepts nil matches", func(t *testing.T) {
		node, p2 := runParser("hello world!", Any(Exact("ffffff"), WS))
		require.Nil(t, node)
		require.Equal(t, 0, p2.Pos)
	})

	t.Run("No parsers", func(t *testing.T) {
		assertNilParser(t, Any())
	})
}

func TestKleene(t *testing.T) {
	t.Run("Matches sequence with sep", func(t *testing.T) {
		node, p2 := runParser("a,b,c,d,e,", Kleene(Chars("a-g"), ","))
		require.False(t, p2.Errored())
		assertSequence(t, node, "a", "b", "c", "d", "e")
		require.Equal(t, 10, p2.Pos)
	})

	t.Run("Matches sequence without sep", func(t *testing.T) {
		node, p2 := runParser("a,b,c,d,e,", Kleene(Any(Chars("a-g"), ",")))
		assertSequence(t, node, "a", ",", "b", ",", "c", ",", "d", ",", "e", ",")
		require.Equal(t, 10, p2.Pos)
	})

	t.Run("Stops on error", func(t *testing.T) {
		node, p2 := runParser("a,b,c,d,e,", Kleene(Chars("a-c"), ","))
		assertSequence(t, node, "a", "b", "c")
		require.Equal(t, 6, p2.Pos)
		require.Equal(t, "d,e,", p2.Get())
	})
}

func TestMany(t *testing.T) {
	t.Run("Matches sequence with sep", func(t *testing.T) {
		node, p2 := runParser("a,b,c,d,e,", Many(Chars("a-g"), Exact(",")))
		assertSequence(t, node, "a", "b", "c", "d", "e")
		require.Equal(t, 10, p2.Pos)
	})

	t.Run("Matches sequence without sep", func(t *testing.T) {
		node, p2 := runParser("a,b,c,d,e,", Many(Any(Chars("abcdefg"), Exact(","))))
		assertSequence(t, node, "a", ",", "b", ",", "c", ",", "d", ",", "e", ",")
		require.Equal(t, 10, p2.Pos)
	})

	t.Run("Stops on error", func(t *testing.T) {
		node, p2 := runParser("a,b,c,d,e,", Many(Chars("abc"), Exact(",")))
		assertSequence(t, node, "a", "b", "c")
		require.Equal(t, 6, p2.Pos)
		require.Equal(t, "d,e,", p2.Get())
	})

	t.Run("Returns error if nothing matches", func(t *testing.T) {
		_, p2 := runParser("a,b,c,d,e,", Many(Chars("def"), Exact(",")))
		require.Equal(t, "offset 0: Expected def", p2.Error.Error())
		require.Equal(t, "a,b,c,d,e,", p2.Get())
	})
}

func TestKleeneUntil(t *testing.T) {
	t.Run("Matches sequence with sep", func(t *testing.T) {
		node, p2 := runParser("a,b,c,d,e,fg", KleeneUntil(Chars("abcde"), "d", ","))
		assertSequence(t, node, "a", "b", "c")
		require.Equal(t, "d,e,fg", p2.Get())
	})

	t.Run("Breaks if separator does not match", func(t *testing.T) {
		node, p2 := runParser("a,b,c,d,e,fg", KleeneUntil(Chars("abcdefg", 1, 1), "y", ","))
		assertSequence(t, node, "a", "b", "c", "d", "e", "f")
		require.Equal(t, "g", p2.Get())
	})
}

func TestManyUntil(t *testing.T) {
	t.Run("Matches sequence until", func(t *testing.T) {
		node, p2 := runParser("a,b,c,d,e,", ManyUntil(Chars("abcdefg"), "d", ","))
		assertSequence(t, node, "a", "b", "c")
		require.Equal(t, 6, p2.Pos)
	})

	t.Run("Returns error until matches early", func(t *testing.T) {
		_, p2 := runParser("a,b,c,d,e,", ManyUntil(Chars("abc"), "a", ","))
		require.Equal(t, "offset 0: Expected something else", p2.Error.Error())
		require.Equal(t, 0, p2.Pos)
		require.Equal(t, "a,b,c,d,e,", p2.Get())
	})
}

type htmlTag struct {
	Name string
}

func TestMap(t *testing.T) {
	parser := Map(And("<", Chars("a-zA-Z0-9"), ">"), func(n *Node) *Node {
		return &Node{Result: htmlTag{n.Children[1].Token}}
	})

	t.Run("sucess", func(t *testing.T) {
		result, _ := runParser("<html>", parser)
		require.Equal(t, htmlTag{"html"}, result.Result)
	})

	t.Run("error", func(t *testing.T) {
		_, ps := runParser("<html", parser)
		require.Equal(t, "offset 5: Expected >", ps.Error.Error())
		require.Equal(t, 0, ps.Pos)
	})
}

func TestMerge(t *testing.T) {
	var bracer Parser
	bracer = And("(", Maybe(&bracer), ")")
	parser := Merge(bracer)

	t.Run("sucess", func(t *testing.T) {
		result, _ := runParser("((()))", parser)
		require.Equal(t, "((()))", result.Token)
	})

	t.Run("error", func(t *testing.T) {
		_, ps := runParser("((())", parser)
		require.Equal(t, "offset 5: Expected )", ps.Error.Error())
		require.Equal(t, 0, ps.Pos)
	})
}

func assertNilParser(t *testing.T, parser Parser) {
	node, p2 := runParser("fff", parser)
	require.Nil(t, node)
	require.Equal(t, 0, p2.Pos)
}

func assertSequence(t *testing.T, node *Node, expected ...string) {
	actual := []string{}

	for _, child := range node.Children {
		actual = append(actual, child.Token)
	}

	require.Equal(t, expected, actual)
}

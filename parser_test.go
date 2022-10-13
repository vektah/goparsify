package goparsify

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParsify(t *testing.T) {
	result := Result{}
	t.Run("strings", func(t *testing.T) {
		Parsify("ff")(NewState("ffooo"), &result)
		require.Equal(t, "ff", result.Token)
	})

	t.Run("parsers", func(t *testing.T) {
		Parsify(Chars("f"))(NewState("ffooo"), &result)
		require.Equal(t, "ff", result.Token)
	})

	t.Run("parser funcs", func(t *testing.T) {
		Parsify(func(p *State, node *Result) { node.Token = "hello" })(NewState("ffooo"), &result)

		require.Equal(t, "hello", result.Token)
	})

	t.Run("*parsers", func(t *testing.T) {
		var parser Parser
		parserfied := Parsify(&parser)
		parser = Chars("f")

		parserfied(NewState("ffooo"), &result)
		require.Equal(t, "ff", result.Token)
	})

	require.Panics(t, func() {
		Parsify(1)
	})
}

func TestParsifyAll(t *testing.T) {
	parsers := ParsifyAll("ff", "gg")

	result := Result{}
	parsers[0](NewState("ffooo"), &result)
	require.Equal(t, "ff", result.Token)

	result = Result{}
	parsers[1](NewState("ffooo"), &result)
	require.Equal(t, "", result.Token)
}

func TestExact(t *testing.T) {
	t.Run("success string", func(t *testing.T) {
		node, ps := runParser("foobar", Exact("fo"))
		require.Equal(t, "fo", node.Token)
		require.Equal(t, "obar", ps.Get())
	})

	t.Run("success char", func(t *testing.T) {
		node, ps := runParser("foobar", Exact("f"))
		require.Equal(t, "f", node.Token)
		require.Equal(t, "oobar", ps.Get())
	})

	t.Run("error", func(t *testing.T) {
		_, ps := runParser("foobar", Exact("bar"))
		require.Equal(t, "bar", ps.Error.expected)
		require.Equal(t, 0, ps.Pos)
	})

	t.Run("error char", func(t *testing.T) {
		_, ps := runParser("foobar", Exact("o"))
		require.Equal(t, "o", ps.Error.expected)
		require.Equal(t, 0, ps.Pos)
	})

	t.Run("eof char", func(t *testing.T) {
		_, ps := runParser("", Exact("o"))
		require.Equal(t, "o", ps.Error.expected)
		require.Equal(t, 0, ps.Pos)
	})
}

func TestChars(t *testing.T) {
	t.Run("full match", func(t *testing.T) {
		node, ps := runParser("foobar", Chars("a-z"))
		require.Equal(t, "foobar", node.Token)
		require.Equal(t, "", ps.Get())
		require.False(t, ps.Errored())
	})

	t.Run("partial match", func(t *testing.T) {
		node, ps := runParser("a1b2c3d4efg", Chars("1-4d-a"))
		require.Equal(t, "a1b2c3d4", node.Token)
		require.Equal(t, "efg", ps.Get())
		require.False(t, ps.Errored())
	})

	t.Run("limited match", func(t *testing.T) {
		node, ps := runParser("a1b2c3d4efg", Chars("1-4d-a", 1, 2))
		require.Equal(t, "a1", node.Token)
		require.Equal(t, "b2c3d4efg", ps.Get())
		require.False(t, ps.Errored())
	})

	t.Run("escaped hyphen", func(t *testing.T) {
		node, ps := runParser(`ab-ab\cde`, Chars(`a\-b`))
		require.Equal(t, "ab-ab", node.Token)
		require.Equal(t, `\cde`, ps.Get())
		require.False(t, ps.Errored())
	})

	t.Run("unescaped hyphen", func(t *testing.T) {
		node, ps := runParser("19-", Chars("0-9"))
		require.Equal(t, "19", node.Token)
		require.Equal(t, "-", ps.Get()) // hyphen shouldn't have been parsed
		require.False(t, ps.Errored())
	})

	t.Run("no match", func(t *testing.T) {
		_, ps := runParser("ffffff", Chars("0-9"))
		require.Equal(t, "offset 0: expected 0-9", ps.Error.Error())
		require.Equal(t, 0, ps.Pos)
	})

	t.Run("no match with min", func(t *testing.T) {
		_, ps := runParser("ffffff", Chars("0-9", 4))
		require.Equal(t, "0-9", ps.Error.expected)
		require.Equal(t, 0, ps.Pos)
	})

	t.Run("test exact matches", func(t *testing.T) {
		node, ps := runParser("aaff", Chars("abcd"))
		require.Equal(t, "aa", node.Token)
		require.Equal(t, 2, ps.Pos)
		require.False(t, ps.Errored())
	})

	t.Run("test not matches", func(t *testing.T) {
		node, ps := runParser("aaff", NotChars("ff"))
		require.Equal(t, "aa", node.Token)
		require.Equal(t, 2, ps.Pos)
		require.False(t, ps.Errored())
	})

	require.Panics(t, func() {
		Chars("a-b", 1, 2, 3)
	})
}

func TestRegex(t *testing.T) {
	t.Run("full match", func(t *testing.T) {
		node, ps := runParser("hello", Regex("[a-z]*"))
		require.Equal(t, "hello", node.Token)
		require.Equal(t, "", ps.Get())
		require.False(t, ps.Errored())
	})

	t.Run("limited match", func(t *testing.T) {
		node, ps := runParser("hello world", Regex("[a-z]*"))
		require.Equal(t, "hello", node.Token)
		require.Equal(t, " world", ps.Get())
		require.False(t, ps.Errored())
	})

	t.Run("no match", func(t *testing.T) {
		_, ps := runParser("1234", Regex("[a-z]*"))
		require.Equal(t, "offset 0: expected [a-z]*", ps.Error.Error())
		require.Equal(t, 0, ps.Pos)
	})

	t.Run("eof", func(t *testing.T) {
		_, ps := runParser("", Regex("[a-z]*"))
		require.Equal(t, "offset 0: expected [a-z]*", ps.Error.Error())
		require.Equal(t, 0, ps.Pos)
	})

	t.Run("alternates 1", func(t *testing.T) {
		_, ps := runParser("a bear in search of honey", Regex("zero|one"))
		require.Equal(t, "offset 0: expected zero|one", ps.Error.Error())
		require.Equal(t, 0, ps.Pos)
	})

	t.Run("alternates 2", func(t *testing.T) {
		_, ps := runParser("one", Regex("zero|one"))
		require.False(t, ps.Errored())
		require.Equal(t, 3, ps.Pos)
	})
}

func TestParseString(t *testing.T) {
	Y := Map("hello", func(n *Result) { n.Result = n.Token })

	t.Run("full match", func(t *testing.T) {
		result, err := Run(Y, "hello")
		require.Equal(t, "hello", result)
		require.NoError(t, err)
	})

	t.Run("partial match", func(t *testing.T) {
		result, err := Run(Y, "hello world")
		require.Equal(t, "hello", result)
		require.Error(t, err)
		require.Equal(t, "left unparsed: world", err.Error())
	})

	t.Run("error", func(t *testing.T) {
		result, err := Run(Y, "world")
		require.Nil(t, result)
		require.Error(t, err)
		require.Equal(t, "offset 0: expected hello", err.Error())
	})
}

func TestAutoWS(t *testing.T) {
	t.Run("ws is not automatically consumed", func(t *testing.T) {
		_, ps := runParser(" hello", NoAutoWS("hello"))
		require.Equal(t, "offset 0: expected hello", ps.Error.Error())
	})

	t.Run("ws is can be explicitly consumed ", func(t *testing.T) {
		result, ps := runParser(" hello", NoAutoWS(Seq(ASCIIWhitespace, "hello")))
		require.Equal(t, "hello", result.Child[1].Token)
		require.Equal(t, " hello", result.Token)
		require.Equal(t, "", ps.Get())
	})

	t.Run("unicode whitespace", func(t *testing.T) {
		result, ps := runParser(" \u202f hello", NoAutoWS(Seq(UnicodeWhitespace, "hello")))
		require.Equal(t, "hello", result.Child[1].Token)
		require.Equal(t, "", ps.Get())
		require.False(t, ps.Errored())
	})
}

func TestUntil(t *testing.T) {
	parser := Until("world", ".")

	t.Run("success", func(t *testing.T) {
		result, ps := runParser("this is the end of the world", parser)
		require.Equal(t, "this is the end of the ", result.Token)
		require.Equal(t, "world", ps.Get())
	})

	t.Run("eof", func(t *testing.T) {
		result, ps := runParser("this is the end of it all", parser)
		require.Equal(t, "this is the end of it all", result.Token)
		require.Equal(t, "", ps.Get())
	})
}

func runParser(input string, parser Parser) (Result, *State) {
	ps := NewState(input)
	result := Result{}
	parser(ps, &result)
	return result, ps
}

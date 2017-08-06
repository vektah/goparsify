package parsec

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParsify(t *testing.T) {
	p := Pointer{"ffooo", 0}

	t.Run("strings", func(t *testing.T) {
		node, _ := Parsify("ff")(p)
		require.Equal(t, "ff", node)
	})

	t.Run("parsers", func(t *testing.T) {
		node, _ := Parsify(CharRun("f"))(p)
		require.Equal(t, "ff", node)
	})

	t.Run("parser funcs", func(t *testing.T) {
		node, _ := Parsify(func(p Pointer) (Node, Pointer) {
			return "hello", p
		})(p)
		require.Equal(t, "hello", node)
	})

	t.Run("*parsers", func(t *testing.T) {
		var parser Parser
		parserfied := Parsify(&parser)
		parser = CharRun("f")

		node, _ := parserfied(p)
		require.Equal(t, "ff", node)
	})

	require.Panics(t, func() {
		Parsify(1)
	})
}

func TestParsifyAll(t *testing.T) {
	parsers := ParsifyAll("ff", "gg")

	result, _ := parsers[0](Pointer{"ffooo", 0})
	require.Equal(t, "ff", result)

	result, _ = parsers[1](Pointer{"ffooo", 0})
	require.Equal(t, NewError(0, "Expected gg"), result)
}

func TestExact(t *testing.T) {
	p := Pointer{"fooo", 0}

	t.Run("success", func(t *testing.T) {
		node, p2 := Exact("fo")(p)
		require.Equal(t, "fo", node)
		require.Equal(t, p.Advance(2), p2)
	})

	t.Run("error", func(t *testing.T) {
		node, p2 := Exact("bar")(p)
		require.Equal(t, NewError(0, "Expected bar"), node)
		require.Equal(t, 0, p2.pos)
	})
}

func TestChar(t *testing.T) {
	p := Pointer{"foobar", 0}

	t.Run("success", func(t *testing.T) {
		node, p2 := Char("fo")(p)
		require.Equal(t, "f", node)
		require.Equal(t, p.Advance(1), p2)
	})

	t.Run("error", func(t *testing.T) {
		node, p2 := Char("bar")(p)
		require.Equal(t, NewError(0, "Expected one of bar"), node)
		require.Equal(t, 0, p2.pos)
	})
}

func TestCharRun(t *testing.T) {
	p := Pointer{"foobar", 0}

	t.Run("success", func(t *testing.T) {
		node, p2 := CharRun("fo")(p)
		require.Equal(t, "foo", node)
		require.Equal(t, p.Advance(3), p2)
	})

	t.Run("error", func(t *testing.T) {
		node, p2 := CharRun("bar")(p)
		require.Equal(t, NewError(0, "Expected some of bar"), node)
		require.Equal(t, 0, p2.pos)
	})
}

func TestCharUntil(t *testing.T) {
	p := Pointer{"foobar", 0}

	t.Run("success", func(t *testing.T) {
		node, p2 := CharRunUntil("z")(p)
		require.Equal(t, "foobar", node)
		require.Equal(t, p.Advance(6), p2)
	})

	t.Run("error", func(t *testing.T) {
		node, p2 := CharRunUntil("f")(p)
		require.Equal(t, NewError(0, "Expected some of f"), node)
		require.Equal(t, 0, p2.pos)
	})
}

func TestWS(t *testing.T) {
	p := Pointer{"  fooo", 0}

	node, p2 := WS(p)
	require.Equal(t, nil, node)
	require.Equal(t, p.Advance(2), p2)
}

func TestRange(t *testing.T) {
	t.Run("full match", func(t *testing.T) {
		node, p := Range("a-z")(Pointer{"foobar", 0})
		require.Equal(t, "foobar", node)
		require.Equal(t, "", p.Get())
	})

	t.Run("partial match", func(t *testing.T) {
		node, p := Range("1-4d-a")(Pointer{"a1b2c3d4efg", 0})
		require.Equal(t, "a1b2c3d4", node)
		require.Equal(t, "efg", p.Get())
	})

	t.Run("limited match", func(t *testing.T) {
		node, p := Range("1-4d-a", 1, 2)(Pointer{"a1b2c3d4efg", 0})
		require.Equal(t, "a1", node)
		require.Equal(t, "b2c3d4efg", p.Get())
	})

	t.Run("no match", func(t *testing.T) {
		node, p := Range("0-9")(Pointer{"ffffff", 0})
		require.Equal(t, NewError(0, "Expected at least 1 more of 0-9"), node)
		require.Equal(t, 0, p.pos)
	})

	t.Run("no match with min", func(t *testing.T) {
		node, p := Range("0-9", 4)(Pointer{"ffffff", 0})
		require.Equal(t, NewError(0, "Expected at least 4 more of 0-9"), node)
		require.Equal(t, 0, p.pos)
	})

	require.Panics(t, func() {
		Range("abcd")
	})

	require.Panics(t, func() {
		Range("a-b", 1, 2, 3)
	})
}

func TestParseString(t *testing.T) {
	t.Run("partial match", func(t *testing.T) {
		result, remaining, err := ParseString("hello", "hello world")
		require.Equal(t, "hello", result)
		require.Equal(t, " world", remaining)
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		result, remaining, err := ParseString("world", "hello world")
		require.Equal(t, nil, result)
		require.Equal(t, "hello world", remaining)
		require.Error(t, err)
		require.Equal(t, "offset 0: Expected world", err.Error())
	})
}

func TestString(t *testing.T) {
	t.Run("test basic match", func(t *testing.T) {
		result, p := String('"')(Pointer{`"hello"`, 0})
		require.Equal(t, `hello`, result)
		require.Equal(t, "", p.Get())
	})

	t.Run("test non match", func(t *testing.T) {
		result, p := String('"')(Pointer{`1`, 0})
		require.Equal(t, NewError(0, `Expected "`), result)
		require.Equal(t, `1`, p.Get())
	})

	t.Run("test unterminated string", func(t *testing.T) {
		result, p := String('"')(Pointer{`"hello `, 0})
		require.Equal(t, NewError(0, `Unterminated string`), result)
		require.Equal(t, `"hello `, p.Get())
	})

	t.Run("test escaping", func(t *testing.T) {
		result, p := String('"')(Pointer{`"hello \"world\""`, 0})
		require.Equal(t, `hello "world"`, result)
		require.Equal(t, ``, p.Get())
	})
}

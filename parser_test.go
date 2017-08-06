package parsec

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParsify(t *testing.T) {
	p := Pointer{"ffooo", 0}

	t.Run("strings", func(t *testing.T) {
		node, _ := Parsify("ff")(p)
		require.Equal(t, NewToken(0, "ff"), node)
	})

	t.Run("parsers", func(t *testing.T) {
		node, _ := Parsify(CharRun("f"))(p)
		require.Equal(t, NewToken(0, "ff"), node)
	})

	t.Run("*parsers", func(t *testing.T) {
		var parser Parser
		parserfied := Parsify(&parser)
		parser = CharRun("f")

		node, _ := parserfied(p)
		require.Equal(t, NewToken(0, "ff"), node)
	})
}

func TestExact(t *testing.T) {
	p := Pointer{"fooo", 0}

	t.Run("success", func(t *testing.T) {
		node, p2 := Exact("fo")(p)
		require.Equal(t, NewToken(0, "fo"), node)
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
		require.Equal(t, NewToken(0, "f"), node)
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
		require.Equal(t, NewToken(0, "foo"), node)
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
		require.Equal(t, NewToken(0, "foobar"), node)
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
	require.Equal(t, "abcdefg", Range("a-g"))
	require.Equal(t, "01234abcd", Range("0-4a-d"))
}

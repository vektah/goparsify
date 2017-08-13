package goparsify

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestState_Advance(t *testing.T) {
	ps := NewState("fooo")
	require.Equal(t, 0, ps.Pos)
	ps.Advance(2)
	require.Equal(t, 2, ps.Pos)
	ps.Advance(1)
	require.Equal(t, 3, ps.Pos)
}

func TestState_Get(t *testing.T) {
	ps := NewState("fooo")
	require.Equal(t, "fooo", ps.Get())
	ps.Advance(1)
	require.Equal(t, "ooo", ps.Get())
	ps.Advance(4)
	require.Equal(t, "", ps.Get())
	ps.Advance(10)
	require.Equal(t, "", ps.Get())
}

func TestState_Errors(t *testing.T) {
	ps := NewState("fooo")

	ps.ErrorHere("hello")
	require.Equal(t, "offset 0: expected hello", ps.Error.Error())
	require.Equal(t, 0, ps.Error.Pos())
	require.True(t, ps.Errored())

	ps.Recover()
	require.False(t, ps.Errored())

	ps.Advance(2)
	ps.ErrorHere("hello2")
	require.Equal(t, "offset 2: expected hello2", ps.Error.Error())
	require.Equal(t, 2, ps.Error.Pos())
	require.True(t, ps.Errored())
}

func TestState_Preview(t *testing.T) {
	require.Equal(t, "", NewState("").Preview(10))
	require.Equal(t, "asdf", NewState("asdf").Preview(10))
	require.Equal(t, "asdfasdfas", NewState("asdfasdfasdf").Preview(10))
}

func TestWhitespaces(t *testing.T) {
	p := Many(Any("hello", "world", "!"))

	_, err := Run(p, "hello world\u2005!", ASCIIWhitespace)
	require.Equal(t, "left unparsed: \u2005!", err.Error())

	_, err = Run(p, "hello world\u2005!", UnicodeWhitespace)
	require.NoError(t, err)
}

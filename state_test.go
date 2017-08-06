package parsec

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestState_Advance(t *testing.T) {
	ps := InputString("fooo")
	require.Equal(t, 0, ps.Pos)
	ps.Advance(2)
	require.Equal(t, 2, ps.Pos)
	ps.Advance(1)
	require.Equal(t, 3, ps.Pos)
}

func TestState_Get(t *testing.T) {
	ps := InputString("fooo")
	require.Equal(t, "fooo", ps.Get())
	ps.Advance(1)
	require.Equal(t, "ooo", ps.Get())
	ps.Advance(4)
	require.Equal(t, "", ps.Get())
	ps.Advance(10)
	require.Equal(t, "", ps.Get())
}

func TestState_Errors(t *testing.T) {
	ps := InputString("fooo")

	ps.ErrorHere("hello")
	require.Equal(t, "offset 0: Expected hello", ps.Error.Error())
	require.Equal(t, 0, ps.Error.Pos())
	require.True(t, ps.Errored())

	ps.ClearError()
	require.False(t, ps.Errored())

	ps.Advance(2)
	ps.ErrorHere("hello2")
	require.Equal(t, "offset 2: Expected hello2", ps.Error.Error())
	require.Equal(t, 2, ps.Error.Pos())
	require.True(t, ps.Errored())
}

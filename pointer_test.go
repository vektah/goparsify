package parsec

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPointer(t *testing.T) {
	p := Pointer{"fooo", 0}

	t.Run("Advances", func(t *testing.T) {
		p2 := p.Advance(2)
		require.Equal(t, Pointer{"fooo", 2}, p2)
		require.Equal(t, Pointer{"fooo", 0}, p)
		require.Equal(t, Pointer{"fooo", 3}, p2.Advance(1))
	})

	t.Run("Get", func(t *testing.T) {
		require.Equal(t, "fooo", p.Get())
		require.Equal(t, "ooo", p.Advance(1).Get())
	})

	t.Run("Remaining", func(t *testing.T) {
		require.Equal(t, 4, p.Remaining())
		require.Equal(t, 0, p.Advance(4).Remaining())
		require.Equal(t, 0, p.Advance(10).Remaining())
	})

	t.Run("Next takes one character", func(t *testing.T) {
		s, p2 := p.Next()
		require.Equal(t, p.Advance(1), p2)
		require.Equal(t, 'f', s)
	})

	t.Run("Next handles EOF", func(t *testing.T) {
		s, p2 := p.Advance(5).Next()
		require.Equal(t, p.Advance(5), p2)
		require.Equal(t, EOF, s)
	})

	t.Run("HasPrefix", func(t *testing.T) {
		require.True(t, p.HasPrefix("fo"))
		require.False(t, p.HasPrefix("ooo"))
		require.True(t, p.Advance(1).HasPrefix("ooo"))
		require.False(t, p.Advance(1).HasPrefix("oooo"))
	})

	t.Run("Accept", func(t *testing.T) {
		s, p2 := p.Accept("abcdef")
		require.Equal(t, "f", s)
		require.Equal(t, p.Advance(1), p2)

		s, p2 = p.Accept("ooooo")
		require.Equal(t, "", s)
		require.Equal(t, p.Advance(0), p2)

		s, p2 = p.Advance(4).Accept("ooooo")
		require.Equal(t, "", s)
		require.Equal(t, p.Advance(4), p2)
	})

	t.Run("AcceptRun", func(t *testing.T) {
		s, p2 := p.AcceptRun("f")
		require.Equal(t, "f", s)
		require.Equal(t, p.Advance(1), p2)

		s, p3 := p.AcceptRun("fo")
		require.Equal(t, "fooo", s)
		require.Equal(t, p.Advance(4), p3)

		s, p4 := p3.AcceptRun("fo")
		require.Equal(t, "", s)
		require.Equal(t, p.Advance(4), p4)
	})

	t.Run("AcceptUntil", func(t *testing.T) {
		s, p2 := p.AcceptUntil("o")
		require.Equal(t, "f", s)
		require.Equal(t, p.Advance(1), p2)

		s, p3 := p2.AcceptRun("o")
		require.Equal(t, "ooo", s)
		require.Equal(t, p.Advance(4), p3)
	})
}

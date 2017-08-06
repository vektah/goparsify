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
		require.Equal(t, "", p.Advance(4).Get())
		require.Equal(t, "", p.Advance(10).Get())
	})
}

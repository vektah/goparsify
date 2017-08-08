package goparsify

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	parser := StringLit(`"'`)
	t.Run("test double match", func(t *testing.T) {
		result, p := runParser(`"hello"`, parser)
		require.Equal(t, `hello`, result.Token)
		require.Equal(t, "", p.Get())
	})

	t.Run("test single match", func(t *testing.T) {
		result, p := runParser(`"hello"`, parser)
		require.Equal(t, `hello`, result.Token)
		require.Equal(t, "", p.Get())
	})

	t.Run("test nested quotes", func(t *testing.T) {
		result, p := runParser(`"hello 'world'"`, parser)
		require.Equal(t, `hello 'world'`, result.Token)
		require.Equal(t, "", p.Get())
	})

	t.Run("test non match", func(t *testing.T) {
		_, p := runParser(`1`, parser)
		require.Equal(t, `"'`, p.Error.Expected)
		require.Equal(t, `1`, p.Get())
	})

	t.Run("test unterminated string", func(t *testing.T) {
		_, p := runParser(`"hello `, parser)
		require.Equal(t, `"`, p.Error.Expected)
		require.Equal(t, `"hello `, p.Get())
	})

	t.Run("test unmatched quotes", func(t *testing.T) {
		_, p := runParser(`"hello '`, parser)
		require.Equal(t, `"`, p.Error.Expected)
		require.Equal(t, 0, p.Pos)
	})

	t.Run("test unterminated escape", func(t *testing.T) {
		_, p := runParser(`"hello \`, parser)
		require.Equal(t, `"`, p.Error.Expected)
		require.Equal(t, 0, p.Pos)
	})

	t.Run("test escaping", func(t *testing.T) {
		result, p := runParser(`"hello \"world\""`, parser)
		require.Equal(t, `hello "world"`, result.Token)
		require.Equal(t, ``, p.Get())
	})

	t.Run("test escaped unicode", func(t *testing.T) {
		result, p := runParser(`"hello \ubeef cake"`, parser)
		require.Equal(t, "", p.Error.Expected)
		require.Equal(t, "hello \uBEEF cake", result.Token)
		require.Equal(t, ``, p.Get())
	})

	t.Run("test invalid escaped unicode", func(t *testing.T) {
		_, p := runParser(`"hello \ucake"`, parser)
		require.Equal(t, "offset 9: Expected [a-f0-9]", p.Error.Error())
		require.Equal(t, 0, p.Pos)
	})

	t.Run("test incomplete escaped unicode", func(t *testing.T) {
		_, p := runParser(`"hello \uca"`, parser)
		require.Equal(t, "offset 9: Expected [a-f0-9]{4}", p.Error.Error())
		require.Equal(t, 0, p.Pos)
	})
}

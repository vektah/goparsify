package json

import (
	stdlibJson "encoding/json"
	"testing"

	parsecJson "github.com/prataprc/goparsec/json"
	"github.com/stretchr/testify/require"
	"github.com/vektah/goparsify"
)

func TestUnmarshal(t *testing.T) {
	t.Run("basic types", func(t *testing.T) {
		result, err := Unmarshal(`true`)
		require.NoError(t, err)
		require.Equal(t, true, result)

		result, err = Unmarshal(`false`)
		require.NoError(t, err)
		require.Equal(t, false, result)

		result, err = Unmarshal(`null`)
		require.NoError(t, err)
		require.Equal(t, nil, result)

		result, err = Unmarshal(`"true"`)
		require.NoError(t, err)
		require.Equal(t, "true", result)
	})

	t.Run("array", func(t *testing.T) {
		result, err := Unmarshal(`[true, null, false]`)
		require.NoError(t, err)
		require.Equal(t, []interface{}{true, nil, false}, result)
	})

	t.Run("object", func(t *testing.T) {
		result, err := Unmarshal(`{"true":true, "false":false, "null": null} `)
		require.NoError(t, err)
		require.Equal(t, map[string]interface{}{"true": true, "false": false, "null": nil}, result)
	})
}

const benchmarkString = `{"true":true, "false":false, "null": null}`

func BenchmarkUnmarshalParsec(b *testing.B) {
	bytes := []byte(benchmarkString)

	for i := 0; i < b.N; i++ {
		scanner := parsecJson.NewJSONScanner(bytes)
		_, remaining := parsecJson.Y(scanner)

		require.True(b, remaining.Endof())
	}
}

func BenchmarkUnmarshalParsify(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Unmarshal(benchmarkString)
		require.NoError(b, err)
	}
	goparsify.DumpDebugStats()
}

func BenchmarkUnmarshalStdlib(b *testing.B) {
	bytes := []byte(benchmarkString)
	var result interface{}
	for i := 0; i < b.N; i++ {
		err := stdlibJson.Unmarshal(bytes, &result)
		require.NoError(b, err)
	}
}

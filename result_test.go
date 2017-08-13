package goparsify

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResult_String(t *testing.T) {
	require.Equal(t, "Hello", Result{Token: "Hello"}.String())
	require.Equal(t, "[Hello,World]", Result{Child: []Result{{Token: "Hello"}, {Token: "World"}}}.String())
	require.Equal(t, "10", Result{Result: 10}.String())
	require.Equal(t, "10", Result{Result: big.NewInt(10)}.String())
}

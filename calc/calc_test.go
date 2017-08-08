package calc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNumbers(t *testing.T) {
	result, err := Calc(`1`)
	require.NoError(t, err)
	require.EqualValues(t, 1, result)
}

func TestAddition(t *testing.T) {
	result, err := Calc(`1+1`)
	require.NoError(t, err)
	require.EqualValues(t, 2, result)
}

func TestSubtraction(t *testing.T) {
	result, err := Calc(`1-1`)
	require.NoError(t, err)
	require.EqualValues(t, 0, result)
}

func TestDivision(t *testing.T) {
	result, err := Calc(`1/2`)
	require.NoError(t, err)
	require.EqualValues(t, .5, result)
}

func TestMultiplication(t *testing.T) {
	result, err := Calc(`1*2`)
	require.NoError(t, err)
	require.EqualValues(t, 2, result)
}

func TestOrderOfOperations(t *testing.T) {
	result, err := Calc(`1+10*2`)
	require.NoError(t, err)
	require.EqualValues(t, 21, result)
}
func TestParenthesis(t *testing.T) {
	result, err := Calc(`(1+10)*2`)
	require.NoError(t, err)
	require.EqualValues(t, 22, result)
}

func TestRecursive(t *testing.T) {
	result, err := Calc(`(1+(2*(3-(4/(5)))))`)
	require.NoError(t, err)
	require.EqualValues(t, 5.4, result)
}

package debug

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegex(t *testing.T) {
	tests := map[string]string{
		"attrs": `	attrs = Map(Some(attr), func(node Result) Result {`,
		"_value": `	_value = Any(_null, _true, _false, _string, _number, _array, _object)`,
		"_object": `_object = Map(Seq("{", Cut, _properties, "}"), func(n Result) Result {`,
		"expr":    `var expr = Exact("foo")`,
		"number":  `number := NumberLit()`,
	}
	for expected, input := range tests {
		t.Run(input, func(t *testing.T) {
			matches := varRegex.FindStringSubmatch(input)
			require.NotNil(t, matches)
			require.Equal(t, expected, matches[1])
		})
	}
}

package html

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	result, err := parse(`<body>hello <p color="blue">world</p></body>`)
	require.NoError(t, err)
	require.Equal(t, htmlTag{Name: "body", Attributes: map[string]string{}, Body: []interface{}{
		"hello ",
		htmlTag{Name: "p", Attributes: map[string]string{"color": "blue"}, Body: []interface{}{"world"}},
	}}, result)
}

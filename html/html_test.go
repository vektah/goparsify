package html

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	result, err := Parse(`<body>hello <p color="blue">world</p></body>`)
	require.NoError(t, err)
	require.Equal(t, Tag{Name: "body", Attributes: map[string]string{}, Body: []interface{}{
		"hello ",
		Tag{Name: "p", Attributes: map[string]string{"color": "blue"}, Body: []interface{}{"world"}},
	}}, result)
}

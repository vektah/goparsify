package html

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vektah/goparsify"
)

func TestParse(t *testing.T) {
	goparsify.EnableLogging(os.Stdout)
	result, err := parse(`<body>hello <p color="blue">world</p></body>`)
	require.NoError(t, err)
	require.Equal(t, htmlTag{Name: "body", Attributes: map[string]string{}, Body: []interface{}{
		"hello ",
		htmlTag{Name: "p", Attributes: map[string]string{"color": "blue"}, Body: []interface{}{"world"}},
	}}, result)
}

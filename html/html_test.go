package html

import (
	"testing"

	"github.com/stretchr/testify/require"
	. "github.com/vektah/goparsify"
)

func TestParse(t *testing.T) {
	result, _, err := Parse(`<body>hello <p color="blue">world</p></body>`)
	require.NoError(t, err)
	require.Equal(t, Tag{Name: "body", Attributes: map[string]string{}, Body: []Node{
		"hello ",
		Tag{Name: "p", Attributes: map[string]string{"color": "blue"}, Body: []Node{"world"}},
	}}, result)
}

package html

import (
	"testing"

	"github.com/stretchr/testify/require"
	. "github.com/vektah/goparsify"
)

func TestParse(t *testing.T) {
	result, _, err := Parse("<body>hello <b>world</b></body>")
	require.NoError(t, err)
	require.Equal(t, Tag{Name: "body", Attributes: map[string]string{}, Body: []Node{
		"hello ",
		Tag{Name: "b", Attributes: map[string]string{}, Body: []Node{"world"}},
	}}, result)
}

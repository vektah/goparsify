package goparsify

import (
	"fmt"
	"strings"
)

var TrashResult = &Result{}

// Result is the output of a parser. Usually only one of its fields will be set and should be though of
// more as a union type. having it avoids interface{} littered all through the parsing code and makes
// the it easy to do the two most common operations, getting a token and finding a child.
type Result struct {
	Token  string
	Child  []Result
	Result interface{}
}

func (r Result) String() string {
	if r.Result != nil {
		if rs, ok := r.Result.(fmt.Stringer); ok {
			return rs.String()
		}
		return fmt.Sprintf("%#v", r.Result)
	}

	if len(r.Child) > 0 {
		children := []string{}
		for _, child := range r.Child {
			children = append(children, child.String())
		}
		return "[" + strings.Join(children, ",") + "]"
	}

	return r.Token
}

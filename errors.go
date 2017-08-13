package goparsify

import "fmt"

// Error represents a parse error. These will often be set, the parser will back up a little and
// find another viable path. In general when combining errors the longest error should be returned.
type Error struct {
	pos      int
	expected string
}

// Pos is the offset into the document the error was found
func (e *Error) Pos() int { return e.pos }

// Error satisfies the golang error interface
func (e *Error) Error() string { return fmt.Sprintf("offset %d: expected %s", e.pos, e.expected) }

// UnparsedInputError is returned by Run when not all of the input was consumed. There may still be a valid result
type UnparsedInputError struct {
	remaining string
}

// Error satisfies the golang error interface
func (e UnparsedInputError) Error() string {
	return "left unparsed: " + e.remaining
}

package parsec

import "fmt"

type Node interface {
}

type Error struct {
	pos     int
	Message string
}

func (e Error) Pos() int      { return e.pos }
func (e Error) Error() string { return fmt.Sprintf("offset %d: %s", e.pos, e.Message) }

func NewError(pos int, message string) Error {
	return Error{pos, message}
}

func IsError(n interface{}) bool {
	_, isErr := n.(Error)
	return isErr
}

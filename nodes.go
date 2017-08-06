package parsec

type Node interface {
	Pos() int
}

type Token struct {
	pos   int
	Value string
}

func (e Token) Pos() int { return e.pos }

func NewToken(pos int, value string) Token {
	return Token{pos, value}
}

type Error struct {
	pos   int
	Error string
}

func (e Error) Pos() int { return e.pos }

func NewError(pos int, message string) Error {
	return Error{pos, message}
}

func IsError(n Node) bool {
	_, isErr := n.(Error)
	return isErr
}

type Sequence struct {
	pos   int
	Nodes []Node
}

func (e Sequence) Pos() int { return e.pos }

func NewSequence(pos int, n ...Node) Sequence {
	return Sequence{pos, n}
}

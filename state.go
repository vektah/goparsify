package parsec

import "fmt"

type Error struct {
	pos      int
	Expected string
}

func (e Error) Pos() int      { return e.pos }
func (e Error) Error() string { return fmt.Sprintf("offset %d: Expected %s", e.pos, e.Expected) }

type State struct {
	Input string
	Pos   int
	Error Error
}

func (s *State) Advance(i int) {
	s.Pos += i
}

func (s *State) Get() string {
	if s.Pos > len(s.Input) {
		return ""
	}
	return s.Input[s.Pos:]
}

func (s *State) ErrorHere(expected string) {
	s.Error.pos = s.Pos
	s.Error.Expected = expected
}

func (s *State) ClearError() {
	s.Error.Expected = ""
}

func (s *State) Errored() bool {
	return s.Error.Expected != ""
}

func InputString(input string) *State {
	return &State{Input: input}
}

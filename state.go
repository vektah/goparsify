package goparsify

import (
	"fmt"
)

// Error represents a parse error. These will often be set, the parser will back up a little and
// find another viable path. In general when combining errors the longest error should be returned.
type Error struct {
	pos      int
	expected string
}

// Pos is the offset into the document the error was found
func (e Error) Pos() int { return e.pos }

// Error satisfies the golang error interface
func (e Error) Error() string { return fmt.Sprintf("offset %d: expected %s", e.pos, e.expected) }

// WSFunc matches a byte and returns true if it is whitespace
type WSFunc func(c byte) bool

// State is the current parse state. It is entirely public because parsers are expected to mutate it during the parse.
type State struct {
	// The full input string
	Input string
	// An offset into the string, pointing to the current tip
	Pos int
	// Error is a secondary return channel from parsers, but used so heavily
	// in backtracking that it has been inlined to avoid allocations.
	Error Error
	// Called to determine what to ignore when WS is called, or when AutoWS fires
	WSFunc   WSFunc
	NoAutoWS bool
}

// NewState creates a new State from a string
func NewState(input string) *State {
	return &State{
		Input: input,
		WSFunc: func(b byte) bool {
			switch b {
			case '\t', '\n', '\v', '\f', '\r', ' ':
				return true
			}
			return false
		},
	}
}

// Advance the Pos along by i bytes
func (s *State) Advance(i int) {
	s.Pos += i
}

// AutoWS consumes all whitespace and advances Pos but can be disabled by the NoAutWS() parser.
func (s *State) AutoWS() {
	if s.NoAutoWS {
		return
	}
	s.WS()
}

// WS consumes all whitespace and advances Pos.
func (s *State) WS() {
	for s.Pos < len(s.Input) && s.WSFunc(s.Input[s.Pos]) {
		s.Pos++
	}
}

// Get the remaining input.
func (s *State) Get() string {
	if s.Pos > len(s.Input) {
		return ""
	}
	return s.Input[s.Pos:]
}

// ErrorHere raises an error at the current position.
func (s *State) ErrorHere(expected string) {
	s.Error.pos = s.Pos
	s.Error.expected = expected
}

// Recover from the current error. Often called by combinators that can match
// when one of their children succeed, but others have failed.
func (s *State) Recover() {
	s.Error.expected = ""
}

// Errored returns true if the current parser has failed.
func (s *State) Errored() bool {
	return s.Error.expected != ""
}

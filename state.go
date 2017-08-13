package goparsify

import (
	"strconv"
	"unicode"
	"unicode/utf8"
)

// State is the current parse state. It is entirely public because parsers are expected to mutate it during the parse.
type State struct {
	// The full input string
	Input string
	// An offset into the string, pointing to the current tip
	Pos int
	// Do not backtrack past this point
	Cut int
	// Error is a secondary return channel from parsers, but used so heavily
	// in backtracking that it has been inlined to avoid allocations.
	Error Error
	// Called to determine what to ignore when WS is called, or when WS fires
	WS VoidParser
}

// ASCIIWhitespace matches any of the standard whitespace characters. It is faster
// than the UnicodeWhitespace parser as it does not need to decode unicode runes.
func ASCIIWhitespace(s *State) {
	for s.Pos < len(s.Input) {
		switch s.Input[s.Pos] {
		case '\t', '\n', '\v', '\f', '\r', ' ':
			s.Pos++
		default:
			return
		}
	}
}

// UnicodeWhitespace matches any unicode space character. Its a little slower
// than the ascii parser because it matches a rune at a time.
func UnicodeWhitespace(s *State) {
	for s.Pos < len(s.Input) {
		r, w := utf8.DecodeRuneInString(s.Get())
		if !unicode.IsSpace(r) {
			return
		}
		s.Pos += w
	}
}

// NoWhitespace disables automatic whitespace matching
func NoWhitespace(s *State) {

}

// NewState creates a new State from a string
func NewState(input string) *State {
	return &State{
		Input: input,
		WS:    UnicodeWhitespace,
	}
}

// Advance the Pos along by i bytes
func (s *State) Advance(i int) {
	s.Pos += i
}

// Get the remaining input.
func (s *State) Get() string {
	if s.Pos > len(s.Input) {
		return ""
	}
	return s.Input[s.Pos:]
}

// Preview of the the next x characters
func (s *State) Preview(x int) string {
	if s.Pos >= len(s.Input) {
		return ""
	}

	quoted := strconv.Quote(s.Get())
	quoted = quoted[1 : len(quoted)-1]
	if len(quoted) >= x {
		return quoted[0:x]
	}

	return quoted
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

package goparsify

import (
	"fmt"
	"unicode"
	"unicode/utf8"
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
	// Called to determine what to ignore when WS is called, or when AutoWS fires
	WS       VoidParser
	NoAutoWS bool
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

// AutoWS consumes all whitespace and advances Pos but can be disabled by the NoAutWS() parser.
func (s *State) AutoWS() {
	if s.NoAutoWS {
		return
	}
	s.WS(s)
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
	if len(s.Input)-s.Pos >= x {
		return s.Input[s.Pos : s.Pos+x]
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

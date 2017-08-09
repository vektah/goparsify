// +build !debug

package goparsify

// NewParser should be called around the creation of every Parser.
// It does nothing normally and should incur no runtime overhead, but when building with -tags debug
// it will instrument every parser to collect valuable timing information displayable with DumpDebugStats.
func NewParser(description string, p Parser) Parser {
	return p
}

// DumpDebugStats will print out the curring timings for each parser if built with -tags debug
func DumpDebugStats() {}

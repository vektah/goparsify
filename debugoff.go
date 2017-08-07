// +build !debug

package goparsify

func NewParser(description string, p Parser) Parser {
	return p
}

func DumpDebugStats() {

}

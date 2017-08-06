package parsec

type Pointer struct {
	input string
	pos   int
}

func (p Pointer) Advance(i int) Pointer {
	return Pointer{p.input, p.pos + i}
}

func (p Pointer) Get() string {
	if p.pos > len(p.input) {
		return ""
	}
	return p.input[p.pos:]
}

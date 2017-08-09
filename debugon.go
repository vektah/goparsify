// +build debug

package goparsify

import (
	"fmt"
	"runtime"
	"sort"
	"strings"
	"time"
)

var parsers []*debugParser

type debugParser struct {
	Description string
	Caller      string
	Next        Parser
	Time        time.Duration
	Calls       int
}

func (dp *debugParser) Parse(ps *State) Result {
	start := time.Now()

	ret := dp.Next(ps)

	dp.Time = dp.Time + time.Since(start)
	dp.Calls++

	return ret
}

func getPackageName(f *runtime.Func) string {
	parts := strings.Split(f.Name(), ".")
	pl := len(parts)

	if parts[pl-2][0] == '(' {
		return strings.Join(parts[0:pl-2], ".")
	} else {
		return strings.Join(parts[0:pl-1], ".")
	}
}

// NewParser should be called around the creation of every Parser.
// It does nothing normally and should incur no runtime overhead, but when building with -tags debug
// it will instrument every parser to collect valuable timing information displayable with DumpDebugStats.
func NewParser(description string, p Parser) Parser {
	fpcs := make([]uintptr, 1)
	caller := ""

	for i := 1; i < 10; i++ {
		n := runtime.Callers(i, fpcs)

		if n != 0 {
			frame := runtime.FuncForPC(fpcs[0] - 1)
			pkg := getPackageName(frame)

			if pkg != "github.com/vektah/goparsify" {
				file, line := frame.FileLine(fpcs[0] - 1)
				caller = fmt.Sprintf("%s:%d", file, line)
				break
			}
		}
	}

	dp := &debugParser{
		Description: description,
		Next:        p,
		Caller:      caller,
	}

	parsers = append(parsers, dp)
	return dp.Parse
}

// DumpDebugStats will print out the curring timings for each parser if built with -tags debug
func DumpDebugStats() {
	sort.Slice(parsers, func(i, j int) bool {
		return parsers[i].Time >= parsers[j].Time
	})

	fmt.Println("Parser stats:")
	for _, parser := range parsers {
		fmt.Printf("%20s\t%10s\t%10d\tcalls\t%s\n", parser.Description, parser.Time.String(), parser.Calls, parser.Caller)
	}
}

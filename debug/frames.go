package debug

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var varRegex = regexp.MustCompile(`(?:var)?\s*(\w*)\s*:?=`)

func getPackageName(f runtime.Frame) string {
	parts := strings.Split(f.Func.Name(), ".")
	pl := len(parts)

	if parts[pl-2][0] == '(' {
		return strings.Join(parts[0:pl-2], ".")
	}

	return strings.Join(parts[0:pl-1], ".")
}

func getVarName(filename string, lineNo int) string {
	f, err := os.Open(filename)
	if err != nil {
		return ""
	}

	scanner := bufio.NewScanner(f)
	for i := 0; i < lineNo; i++ {
		scanner.Scan()
	}

	line := scanner.Text()
	if matches := varRegex.FindStringSubmatch(line); matches != nil {
		return matches[1]
	}
	return ""
}

// GetDefinition returns the name of the variable and location this parser was defined by walking up the stack
func GetDefinition() (varName string, location string) {
	pc := make([]uintptr, 64)
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])

	var frame runtime.Frame
	more := true
	for more {
		frame, more = frames.Next()
		pkg := getPackageName(frame)
		if pkg == "github.com/ijt/goparsify" || pkg == "github.com/ijt/goparsify/debug" {
			continue
		}

		varName := getVarName(frame.File, frame.Line)
		if varName != "" {
			return varName, fmt.Sprintf("%s:%d", filepath.Base(frame.File), frame.Line)
		}
	}

	return "", ""
}

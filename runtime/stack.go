package runtime

import (
	"runtime"
	"strings"
)

const (
	maxStackDepth = 30
	callerSkip    = 2
)

// A StackFrame contains all necessary information about to generate a line
// in a call stack.
type StackFrame struct {
	// The path to the file containing this ProgramCounter
	File string
	// The LineNumber in that file
	LineNumber int
	// The Name of the function that contains this ProgramCounter
	Name string
	// The Package that contains this function
	Package string
	// The underlying ProgramCounter
	ProgramCounter uintptr
}

// Stack produces a stack trace for the current caller.
func Stack() []StackFrame {
	stack := make([]uintptr, maxStackDepth)
	length := runtime.Callers(callerSkip, stack)

	frames := []StackFrame{}

	for _, pc := range stack[:length] {
		frame := StackFrame{ProgramCounter: pc}
		if frame.Func() == nil {
			frames = append(frames, frame)
			continue
		}

		// pc -1 because the program counters we use are usually return addresses,
		// and we want to show the line that corresponds to the function call
		frame.File, frame.LineNumber = frame.Func().FileLine(pc - 1)
		frame.Package, frame.Name = packageAndName(frame.Func())

		frames = append(frames, frame)
	}

	return frames
}

// Func returns the function that contained this frame.
func (frame *StackFrame) Func() *runtime.Func {
	if frame.ProgramCounter == 0 {
		return nil
	}

	return runtime.FuncForPC(frame.ProgramCounter)
}

func packageAndName(fn *runtime.Func) (name, pkg string) {
	name = fn.Name()

	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//  runtime/debug.*T·ptrmethod
	// and want
	//  *T.ptrmethod
	// Since the package path might contains dots (e.g. code.google.com/...),
	// we first remove the path prefix if there is one.
	if lastSlash := strings.LastIndex(name, "/"); lastSlash >= 0 {
		pkg += name[:lastSlash] + "/"
		name = name[lastSlash+1:]
	}
	if period := strings.Index(name, "."); period >= 0 {
		pkg += name[:period]
		name = name[period+1:]
	}

	name = strings.ReplaceAll(name, "·", ".")
	return pkg, name
}

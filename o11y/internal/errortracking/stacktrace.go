package errortracking

import (
	"runtime"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/posthog/posthog-go"
)

const goKit = "github.com/lcnascimento/go-kit/"

var mainModule = sync.OnceValue(func() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}

	return info.Main.Path
})

func stackTraceExtractor() posthog.DefaultStackTraceExtractor {
	return posthog.DefaultStackTraceExtractor{InAppDecider: inApp}
}

// stackTrace returns the stack up to the go-kit frame that logged the error, outermost frame first.
func stackTrace(extractor posthog.DefaultStackTraceExtractor) *posthog.ExceptionStacktrace {
	stack := extractor.GetStackTrace(2) //nolint:mnd // Skips runtime.Callers and GetStackTrace.

	frames := stack.Frames
	for len(frames) > 0 && fromGoKit(frames[len(frames)-1].Function) {
		frames = frames[:len(frames)-1]
	}
	stack.Frames = frames

	return stack
}

//nolint:gocritic // posthog.InAppDecider takes the frame by value.
func inApp(frame runtime.Frame) bool {
	main := mainModule()
	pkg := packageOf(frame.Function)

	return main != "" && (pkg == main || pkg == main+"_test" || strings.HasPrefix(pkg, main+"/"))
}

// fromGoKit reports whether function belongs to a go-kit package, test packages excluded.
func fromGoKit(function string) bool {
	pkg := packageOf(function)

	return strings.HasPrefix(pkg, goKit) && !strings.HasSuffix(pkg, "_test")
}

// packageOf returns the import path of a fully qualified function name, such as "example.com/a/b.(*T).M".
func packageOf(function string) string {
	slash := strings.LastIndex(function, "/")
	if dot := strings.Index(function[slash+1:], "."); dot >= 0 {
		return function[:slash+1+dot]
	}

	return function
}

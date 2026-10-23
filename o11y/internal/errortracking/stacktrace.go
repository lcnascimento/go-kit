package errortracking

import (
	"runtime"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/posthog/posthog-go"
)

const (
	goKit = "github.com/lcnascimento/go-kit/"

	skipCallersAndExtractor = 2
)

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

func stackTrace(extractor posthog.DefaultStackTraceExtractor) *posthog.ExceptionStacktrace {
	stack := extractor.GetStackTrace(skipCallersAndExtractor)

	frames := stack.Frames
	for len(frames) > 0 && fromGoKitPackage(frames[len(frames)-1].Function) {
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

func fromGoKitPackage(function string) bool {
	pkg := packageOf(function)

	return strings.HasPrefix(pkg, goKit) && !strings.HasSuffix(pkg, "_test")
}

func packageOf(function string) string {
	slash := strings.LastIndex(function, "/")
	if dot := strings.Index(function[slash+1:], "."); dot >= 0 {
		return function[:slash+1+dot]
	}

	return function
}

package main

import (
	"fmt"

	"github.com/lcnascimento/go-kit/runtime"
)

func main() {
	// ctx := runtime.ContextWithOSSignalCancellation()
	// ctx = ctx

	err := func3()
	fmt.Println(err)
}

func func1() error {
	for _, f := range runtime.Stack() {
		fmt.Printf("%s:%d (%s)\n", f.File, f.LineNumber, f.Name)
	}

	return fmt.Errorf("fake error")
}

func func2() error {
	return func1()
}

func func3() error {
	return func2()
}

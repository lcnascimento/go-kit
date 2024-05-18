package main

import (
	"fmt"

	"github.com/lcnascimento/go-kit/runtime"
)

func main() {
	fmt.Println("### Stack ####")
	func3()

	fmt.Println("\n\n### Caller ####")
	fmt.Println(runtime.Caller().String())
}

func func1() {
	for _, f := range runtime.Stack() {
		fmt.Println(f.String())
	}
}

func func2() {
	func1()
}

func func3() {
	func2()
}

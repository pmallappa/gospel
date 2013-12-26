package util

import (
	"fmt"
	"runtime"
)

func PrintMe() {
	pc, file, line, _ := runtime.Caller(1)
	fmt.Printf("%s=>%20s:%d\n", runtime.FuncForPC(pc).Name(), file, line)
}

func printMyFunc(s string, n int) {
	pc, _, line, _ := runtime.Caller(n)
	fmt.Printf("%s=>%s:%d\n", s, runtime.FuncForPC(pc).Name(), line)
}

func Entered() {
	if dbg > 0 {
		printMyFunc("Entering", 2)
	}
}

func Exiting() {
	if dbg > 0 {
		printMyFunc("Exiting", 2)
	}
}

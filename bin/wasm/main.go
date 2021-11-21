package main

import (
	"fmt"
	"monkey/repl"
	"syscall/js"
)

func StartStepWrapper(this js.Value, s []js.Value) interface{} {
	if len(s) == 0 {
		return js.ValueOf("")
	}

	output := repl.StartStep(s[0].String())
	return js.ValueOf(output)
}

func main() {
	c := make(chan struct{}, 0)
	fmt.Println("Initializing wasm")
	js.Global().Set("runReplStep", js.FuncOf(StartStepWrapper))
	<-c
}

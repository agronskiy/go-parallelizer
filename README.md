# go-parallelizer
Go package that allows to perform parallelized computations using Go paradigm of input-output channel, and gracefully close workers.

# Usage

Let's assume you want to build a pipeline of two consecutive steps, each being paralLelized over `N` workers.

To use this simple package, all you need is to wrap those tasks into functions which take `gopar.Msg`
(which is just a rename of `interface{}`) and output `gopar.Msg`.

Those can be later channeled one after another.

```go
package main

import (
	"math"

	"github/agronskiy/go-parallelizer"
)

func main() {
	// imagine we want to first square the input and second take a square root of it.
	firstTaskFun := func(i Msg) Msg {
		ii := i.(int)
		return (float64)(ii * ii)
	}

	secondTaskFun := func(i Msg) Msg {
		ii := i.(float64)
		return math.Sqrt(ii)
	}

	inputQueue := make(chan Msg)
	outputQueue := MakeRunner(
		MakeRunner(
			inputQueue,
			firstTaskFun,
			runtime.NumCPU()),
		secondTaskFun,
		runtime.NumCPU(),
	)
	go func() {
		for i := 0; i < 20; i++ {
			inputQueue <- Msg(i)
		}

		close(inputQueue)
	}()

	// consuming final output
	processOutput := func(outputQueue <-chan Msg) {
		for out := range outputQueue {
			fmt.Printf("%v\n", out.(float64))
		}
	}
	processOutput(outputQueue)
}

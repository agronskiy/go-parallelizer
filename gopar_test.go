package gopar

import (
	"fmt"
	"math"
	"runtime"
	"testing"
)

func TestBasicChanneling(t *testing.T) {

	inputQueue := make(chan Msg)
	taskFun := func(i Msg) Msg {
		ii := i.(int)
		return ii * ii
	}

	processOutput := func(outputQueue <-chan Msg) {
		for out := range outputQueue {
			fmt.Printf("%v\n", out.(int))
		}
	}

	outputQueue := MakeRunner(inputQueue, taskFun, runtime.NumCPU())
	go func() {
		for i := 0; i < 20; i++ {
			inputQueue <- Msg(i)
		}

		close(inputQueue)
	}()

	processOutput(outputQueue)

}

func TestConsecutiveChanneling(t *testing.T) {
	// imagine we want to first square the input and second take a square root of it.
	inputQueue := make(chan Msg)
	firstTaskFun := func(i Msg) Msg {
		ii := i.(int)
		return (float64)(ii * ii)
	}

	secondTaskFun := func(i Msg) Msg {
		ii := i.(float64)
		return math.Sqrt(ii)
	}

	processOutput := func(outputQueue <-chan Msg) {
		for out := range outputQueue {
			fmt.Printf("%v\n", out.(float64))
		}
	}

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

	processOutput(outputQueue)
}

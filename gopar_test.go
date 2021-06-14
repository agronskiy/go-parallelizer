package gopar

import (
	"fmt"
	"runtime"
	"testing"
)

func TestBasicChanneling(t *testing.T) {

	inputQueue := make(chan InputTask)
	taskFun := func(i InputTask) OutputResult {
		ii := i.(int)
		return OutputResult(ii * ii)
	}

	processOutput := func(outputQueue <-chan OutputResult) {
		for out := range outputQueue {
			fmt.Printf("%v\n", out.(int))
		}
	}

	outputQueue := MakeRunner(inputQueue, taskFun, runtime.NumCPU())

	go func() {
		for i := 0; i < 20; i++ {
			inputQueue <- InputTask(i)
		}

		close(inputQueue)
	}()

	processOutput(outputQueue)

}

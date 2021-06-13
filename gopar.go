// TODOs:
// Generalize and make it possible to call
// MakeRunner(input chan, func, numWorkers) -> output chan.

package main

import (
	"fmt"
	"runtime"
)

type (
	InputTask    interface{}
	OutputResult interface{}
)

func MakeRunner(numWorkers int) (chan<- InputTask, <-chan OutputResult) {
	var (
		inputQueue  = make(chan InputTask, numWorkers)
		outputQueue = make(chan OutputResult, numWorkers)

		counterCh = make(chan int)

		numOpenWorkers = 0
	)

	// Create actual runner. It:
	// 1. spawns workers
	// 2. listens to their increment/decrement counts how many are still working
	// 	  and closes output channel
	go func() {
		for i := 0; i < numWorkers; i++ {
			go worker(inputQueue, outputQueue, counterCh)
		}

		for n := range counterCh {
			numOpenWorkers += n
			if numOpenWorkers > 0 {
				continue
			}
			close(outputQueue)
		}
	}()

	return inputQueue, outputQueue
}

func worker(
	inputQueue <-chan InputTask,
	outputQueue chan<- OutputResult,
	counterCh chan<- int,
) {
	// worker sends +1 to the main counter when starting, and -1 when closing.
	counterCh <- 1
	defer func() {
		counterCh <- -1
	}()

	// Worker will
	for i := range inputQueue {
		// WIP, TODO figure out the task logic
		i := i.(int)
		var out OutputResult = i * i
		outputQueue <- out
	}
}

func processOutput(outputQueue <-chan OutputResult) {

	// This will create files under 'tags/...'
	for out := range outputQueue {
		fmt.Printf("%v\n", out.(int))
	}
}

func main() {
	inputQueue, outputQueue := MakeRunner(runtime.NumCPU())

	go func() {
		for i := 0; i < 20; i++ {
			inputQueue <- InputTask(i)
		}

		close(inputQueue)
	}()

	processOutput(outputQueue)
}

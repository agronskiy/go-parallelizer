package main

import (
	"runtime"
)

type (
	InputTask    interface{}
	OutputResult interface{}
)

func MakeRunner(numWorkers int) (chan<- InputTask, <-chan OutputResult) {
	var (
		inputQueue        = make(chan InputTask, numWorkers)
		intermediateQueue = make(chan OutputResult, numWorkers)
		outputQueue       = make(chan OutputResult, numWorkers)

		counterCh = make(chan int)

		numOpenWorkers = 0
	)

	stopGracefully := func() {
		for {
			// First, drain remaining results, and only then stop.
			select {
			case out := <-intermediateQueue:
				outputQueue <- out
			default:
				close(outputQueue)
				return
			}
		}
	}

	// Create actual runner. It:
	// 1. spawns workers
	// 2a. listens to their output
	// 2b. does bookkeeping, counts how many are still working and closes output channel
	go func() {
		for i := 0; i < numWorkers; i++ {
			go worker(inputQueue, intermediateQueue, counterCh)
		}

		for {
			select {
			case out := <-intermediateQueue:
				outputQueue <- out
			case n := <-counterCh:
				numOpenWorkers += n
				if numOpenWorkers > 0 {
					continue
				}
				stopGracefully()
			}
		}
	}()

	return inputQueue, outputQueue
}

func worker(
	inputQueue <-chan InputTask,
	intermediateQueue chan<- OutputResult,
	counterCh chan<- int,
) {
	counterCh <- 1
	defer func() {
		counterCh <- -1
	}()

	for range inputQueue {
		// WIP, TODO figure out the task logic
		var out OutputResult = 10
		intermediateQueue <- out
	}
}

func processOutput(outputQueue <-chan OutputResult) {

	// This will create files under 'tags/...'
	for range outputQueue {
		// WIP, TODO do smth here
	}
}

func main() {
	inputQueue, outputQueue := MakeRunner(runtime.NumCPU())

	go func() {
		for i := 0; i < 10; i++ {
			inputQueue <- InputTask(i)
		}

		close(inputQueue)
	}()

	processOutput(outputQueue)
}

// TODOs:
// Generalize and make it possible to call
// MakeRunner(input chan, func, numWorkers) -> output chan.

package gopar

type (
	InputTask    interface{}
	OutputResult interface{}
)

func MakeRunner(
	inputQueue <-chan InputTask,
	taskFun func(InputTask) OutputResult,
	numWorkers int,
) <-chan OutputResult {
	var (
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
			go worker(inputQueue, outputQueue, taskFun, counterCh)
		}

		for n := range counterCh {
			numOpenWorkers += n
			if numOpenWorkers > 0 {
				continue
			}
			close(outputQueue)
		}
	}()

	return outputQueue
}

func worker(
	inputQueue <-chan InputTask,
	outputQueue chan<- OutputResult,
	taskFun func(InputTask) OutputResult,
	counterCh chan<- int,
) {
	// worker sends +1 to the main counter when starting, and -1 when closing.
	counterCh <- 1
	defer func() {
		counterCh <- -1
	}()

	// Worker will
	for i := range inputQueue {
		outputQueue <- taskFun(i)
	}
}

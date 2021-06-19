package gopar

type (
	Msg interface{}
)

func MakeRunner(
	inputQueue <-chan Msg,
	taskFun func(Msg) Msg,
	numWorkers int,
) <-chan Msg {
	var (
		outputQueue    = make(chan Msg, numWorkers)
		counterCh      = make(chan int)
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
	inputQueue <-chan Msg,
	outputQueue chan<- Msg,
	taskFun func(Msg) Msg,
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

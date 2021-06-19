package gopar

import (
	"bytes"
	"log"
	"math"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func captureOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}

func TestBasicChanneling(t *testing.T) {
	str := captureOutput(func() {
		inputQueue := make(chan Msg)
		taskFun := func(i Msg) Msg {
			ii := i.(int)
			return ii * ii
		}

		processOutput := func(outputQueue <-chan Msg) {
			for out := range outputQueue {
				log.Printf("%v\n", out.(int))
			}
		}

		outputQueue := MakeRunner(inputQueue, taskFun, runtime.NumCPU())
		go func() {
			for i := 0; i < 5; i++ {
				inputQueue <- Msg(i)
			}

			close(inputQueue)
		}()

		processOutput(outputQueue)
	})

	assert.Contains(t, str, "0")
	assert.Contains(t, str, "1")
	assert.Contains(t, str, "4")
	assert.Contains(t, str, "9")
	assert.Contains(t, str, "16")
}

func TestConsecutiveChanneling(t *testing.T) {
	str := captureOutput(func() {
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
				log.Printf("%v\n", out.(float64))
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
			for i := 0; i < 5; i++ {
				inputQueue <- Msg(i)
			}

			close(inputQueue)
		}()

		processOutput(outputQueue)
	})
	assert.Contains(t, str, "0")
	assert.Contains(t, str, "1")
	assert.Contains(t, str, "2")
	assert.Contains(t, str, "3")
	assert.Contains(t, str, "4")
}

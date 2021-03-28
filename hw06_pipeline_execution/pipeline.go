package hw06pipelineexecution

import (
	"sync"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

// Mutex for locking/unlocking isDone flag.
var mutex = &sync.RWMutex{}

// Flag to determine is terminating flag sent.
var isDoneFlag bool = false

// IsDone flag setter.
func setDoneFlag(value bool) {
	mutex.Lock()
	isDoneFlag = value
	mutex.Unlock()
}

// IsDone flag getter.
func isDone() bool {
	defer mutex.RUnlock()
	mutex.RLock()
	return isDoneFlag
}

// Returning Stage that closes next Stage input Channel, if done signal has been sent.
func getStageSplitter() func(in In) Out {
	return func(in In) Out {
		out := make(Bi)
		go func() {
			for v := range in {
				if isDone() {
					// If terminating signal has been send, closing output channel,
					// to break the next Stage input Channel reading loop
					close(out)
					return
				}
				out <- v
			}
			close(out)
		}()
		return out
	}
}

// Executes pipeline by given input channel and stages slice.
func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// IsDone flag is false by default
	setDoneFlag(false)

	// Launching a goroutine for terminating channel listening
	if done != nil {
		go func() {
			<-done
			setDoneFlag(true)
		}()
	}

	// Splitting up the stages slice by splitter stage
	splitStages := []Stage{}
	for _, stage := range stages {
		splitStages = append(splitStages, getStageSplitter())
		splitStages = append(splitStages, stage)
	}

	// Building channels chain, connected by appropriate channels
	resultChannel := in
	for _, stage := range splitStages {
		resultChannel = stage(resultChannel)
	}

	return resultChannel
}

package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

// Executes pipeline by given input channel and stages slice.
func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// Splitter stage closes next Stage input Channel, if done signal has been sent
	stageSplitter := func(in In) Out {
		out := make(Bi)
		go func() {
			// Closing "out" channel anyway, when returning from the function
			defer close(out)
			for {
				// Listening both "done" and "in" channels
				select {
				// In case of "done" channel closing, return from stage and close "out" channel
				case <-done:
					return
				// In case of previous stage wrote in the "in" channel, or closed them
				case v, ok := <-in:
					// Returning if "in" channel is closed (and close "out" channel as well)
					if !ok {
						return
					}
					// If "in" channel is not closed, trying to write received value to "out" channel
					// Using "select" syntax, to prevent locking, if the receiving stage is already not listening,
					// but the "out" channel is still not closed
					select {
					// We need to listen "done" channel twice, in case of "input" value already received
					// but the receiver stage is not responding
					case <-done:
						return
					// And finally writing to the "out" channel
					case out <- v:
					}
				}
			}
		}()
		return out
	}

	// Splitting up the stages slice by splitter stage
	splitStages := make([]Stage, 0, len(stages)*2)
	for _, stage := range stages {
		splitStages = append(splitStages, stageSplitter)
		splitStages = append(splitStages, stage)
	}

	// Building channels chain, connected by appropriate channels
	resultChannel := in
	for _, stage := range splitStages {
		resultChannel = stage(resultChannel)
	}

	return resultChannel
}

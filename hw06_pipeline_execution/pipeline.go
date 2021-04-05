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
			for {
				select {
				case <-done:
					close(out)
					return
				case v, ok := <-in:
					if !ok {
						close(out)
						return
					}
					out <- v
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

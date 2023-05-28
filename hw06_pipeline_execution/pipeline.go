package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	result := in
	for _, stage := range stages {
		chain := make(Bi)
		go func(input In, chain Bi) {
			defer close(chain)
			for {
				select {
				case <-done:
					return
				case value, ok := <-input:
					if !ok {
						return
					}
					chain <- value
				}
			}
		}(result, chain)
		result = stage(chain)
	}
	return result
}

package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

//nolint:revive
func ExecutePipeline(in In, done In, stages ...Stage) Out {
	mapOut := make(map[int]Bi, len(stages))
	for i := 0; i < len(stages); i++ {
		mapOut[i] = make(Bi)
	}

	for i, stage := range stages {
		var inLocal In
		if i == 0 {
			inLocal = in
		} else {
			inLocal = mapOut[i-1]
		}

		go func(in In, stage Stage, out Bi) {
		loop:
			for v := range stage(in) {
				select {
				case _, ok := <-done:
					if !ok {
						break loop
					}
				default:
					out <- v
				}
			}
			close(out)
			for range stage(in) {
			}
		}(inLocal, stage, mapOut[i])
	}

	return mapOut[len(stages)-1]
}

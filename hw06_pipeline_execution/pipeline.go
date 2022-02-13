package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func checkDone(in, done In) Out {
	out := make(Bi)
	go func() {
		defer func() {
			close(out)
			for range in {
			}
		}()
		for {
			select {
			case <-done:
				return
			default:
			}
			select {
			case <-done:
				return
			case inVal, ok := <-in:
				if !ok {
					return
				}
				out <- inVal
			}
		}
	}()
	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	if stages == nil {
		return out
	}
	for _, stage := range stages {
		out = checkDone(stage(out), done)
	}

	return out
}

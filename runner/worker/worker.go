package worker

import (
	"github.com/Sharktheone/mcp262/runner/rebuild"
	"github.com/Sharktheone/mcp262/runner/results"
	"github.com/Sharktheone/mcp262/runner/test"
	"sync"
)

type Job struct {
	FullPath     string
	RelativePath string
}

func Worker(id int, jobs <-chan Job, results chan<- results.Result, wg *sync.WaitGroup, loc *rebuild.EngineLocation) {
	defer wg.Done()

	for job := range jobs {
		engine := loc.GetPath()

		res := test.RunTest(job.RelativePath, job.FullPath, engine)

		results <- res
	}
}

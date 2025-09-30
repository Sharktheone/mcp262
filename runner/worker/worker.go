package worker

import (
	"sync"

	"github.com/Sharktheone/mcp262/runner/rebuild"
	"github.com/Sharktheone/mcp262/runner/results"
	"github.com/Sharktheone/mcp262/runner/test"
)

type Job struct {
	FullPath     string
	RelativePath string
}

func Worker(id int, root string, jobs <-chan Job, results chan<- results.Result, wg *sync.WaitGroup, loc *rebuild.EngineLocation) {
	defer wg.Done()

	for job := range jobs {
		engine := loc.GetPath()

		res := test.RunTest(job.RelativePath, job.FullPath, engine, root)

		results <- res
	}
}

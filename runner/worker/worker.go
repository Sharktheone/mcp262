package worker

import (
	"sync"

	"github.com/Sharktheone/mcp262/runner/results"
	"github.com/Sharktheone/mcp262/runner/test"
)

func Worker(id int, jobs <-chan string, results chan<- results.Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for path := range jobs {
		res := test.RunTest(path)

		results <- res
	}
}

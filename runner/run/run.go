package run

import (
	"github.com/Sharktheone/mcp262/runner/rebuild"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Sharktheone/mcp262/runner/results"
	"github.com/Sharktheone/mcp262/runner/status"
	"github.com/Sharktheone/mcp262/runner/worker"
)

var SKIP = []string{
	"intl402",
	"staging",
}

func RunTestsInDir(testRoot string, testDir string, repoRoot string, workers int, rebuildEngine bool) *results.TestResults {
	num := countTests(filepath.Join(testRoot, testDir))

	loc, cancel, err := rebuild.RebuildEngine(repoRoot, num, rebuildEngine)

	if err != nil {
		log.Printf("Failed to rebuild engine: %v", err)
		return nil
	}

	res := testsInDir(testRoot, testDir, workers, loc, num)

	cancel()

	return res
}

func testsInDir(testRoot string, testDir string, workers int, loc *rebuild.EngineLocation, num uint32) *results.TestResults {
	jobs := make(chan worker.Job, workers*8)
	testsDir := filepath.Join(testRoot, testDir)

	resultsChan := make(chan results.Result, workers*8)

	wg := &sync.WaitGroup{}

	wg.Add(workers)

	for i := range workers {
		go worker.Worker(i, jobs, resultsChan, wg, loc)
	}

	testResults := results.New(num)

	go func() {
		for res := range resultsChan {
			testResults.Add(res)
		}
	}()

	now := time.Now()
	_ = filepath.Walk(testsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			//log.Printf("Failed to get file info for %s: %v", path, err)
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if strings.Contains(path, "_FIXTURE") {
			return nil
		}

		p, err := filepath.Rel(testRoot, path)
		if err != nil {
			log.Printf("Failed to get relative path for %s: %v", path, err)
			return nil
		}

		for _, skip := range SKIP {
			if strings.HasPrefix(p, skip) {
				resultsChan <- results.Result{
					Status:   status.SKIP,
					Msg:      "skip",
					Path:     p,
					MemoryKB: 0,
					Duration: 0,
				}

				return nil
			}
		}

		jobs <- worker.Job{
			FullPath:     path,
			RelativePath: p,
		}

		return nil
	})

	close(jobs)

	wg.Wait()
	log.Printf("Finished running %d tests in %s", num, time.Since(now).String())

	close(resultsChan)

	return testResults
}

func countTests(path string) uint32 {
	var num uint32 = 0

	_ = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info == nil {
			log.Printf("Failed to get file info for %s", path)
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if strings.Contains(path, "_FIXTURE") {
			return nil
		}

		num++

		return nil
	})

	return num
}

package run

import (
	"github.com/Sharktheone/mcp262/runner/rebuild"
	"github.com/Sharktheone/mcp262/runner/test"
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

func RunTestsInDir(testRoot string, testDir string, repoRoot string, workers int, rebuildEngine bool) (*results.TestResults, error) {
	num := countTests(filepath.Join(testRoot, testDir))

	loc, cancel, err := rebuild.RebuildEngine(repoRoot, num, rebuildEngine)

	if err != nil {
		return nil, err
	}

	res := testsInDir(testRoot, testDir, workers, loc, num)

	cancel()

	return res, nil
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

func RunSingleTest(testRoot string, testPath string, repoRoot string, rebuildEngine bool) (results.Result, error) {
	loc, cancel, err := rebuild.RebuildEngine(repoRoot, 1, rebuildEngine)

	log.Printf("Engine located at: %s", loc.GetPath())

	cancel()

	if err != nil {
		return results.Result{}, err
	}

	engine := loc.GetPath()

	fullPath := filepath.Join(testRoot, testPath)

	return test.RunTest(testPath, fullPath, engine), nil
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

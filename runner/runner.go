package runner

import (
	"errors"
	"github.com/Sharktheone/mcp262/provider"
	"github.com/Sharktheone/mcp262/runner/run"
)

type Runner struct {
	testRoot string
	repoRoot string
	workers  int
}

func New(testRoot, repoRoot string, workers int) *Runner {
	return &Runner{
		testRoot: testRoot,
		repoRoot: repoRoot,
		workers:  workers,
	}
}

func (r *Runner) RerunTest(testPath string, rebuild bool) (provider.TestResult, error) {
	res, err := run.RunSingleTest(r.testRoot, testPath, r.repoRoot, rebuild)
	if err != nil {
		return provider.TestResult{}, err
	}

	return provider.TestResult{
		TestPath: res.Path,
		Status:   res.Status.String(),
		Output:   res.Msg,
		Duration: res.Duration.String(),
	}, nil

}

func (r *Runner) RerunTestsInDir(dir string, rebuild bool) (map[string]provider.TestResult, error) {
	results, err := run.RunTestsInDir(r.testRoot, dir, r.repoRoot, r.workers, rebuild)
	if err != nil {
		return nil, err
	}

	testResults := make(map[string]provider.TestResult)
	for _, res := range results.TestResults {
		testResults[res.Path] = provider.TestResult{
			TestPath: res.Path,
			Status:   string(res.Status),
			Output:   res.Msg,
			Duration: res.Duration.String(),
		}
	}

	return testResults, nil
}

func (r *Runner) RerunFailedTestsInDir(dir string, rebuild bool) (map[string]provider.TestResult, error) {
	return nil, errors.New("not implemented")
}

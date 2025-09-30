package runner

import (
	"errors"

	"github.com/Sharktheone/mcp262/provider"
	"github.com/Sharktheone/mcp262/runner/ci"
	"github.com/Sharktheone/mcp262/runner/results"
	"github.com/Sharktheone/mcp262/runner/run"
)

type Runner struct {
	testRoot string
	repoRoot string
	workers  int

	prev *results.TestResults
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
	tres, err := run.RunTestsInDir(r.testRoot, dir, r.repoRoot, r.workers, rebuild)
	if err != nil {
		return nil, err
	}

	testResults := make(map[string]provider.TestResult)
	for _, res := range tres.TestResults {
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

func (r *Runner) RerunTestsInDirChanges(dir string, rebuild bool) ([]provider.TestDiff, error) {
	tres, err := run.RunTestsInDir(r.testRoot, dir, r.repoRoot, r.workers, rebuild)
	if err != nil {
		return nil, err
	}

	prev, err := r.getPrevResults()

	if err != nil {
		return nil, err
	}

	diff := tres.ComputeDiff(prev)

	diffs := make([]provider.TestDiff, 0, len(diff))

	for d, items := range diff {
		paths := make([]string, len(items))
		for i, item := range items {
			paths[i] = item.Path()
		}

		td := provider.TestDiff{
			From:  d.From.String(),
			To:    d.To.String(),
			Items: paths,
		}

		diffs = append(diffs, td)
	}

	return diffs, nil
}

func (r *Runner) RerunFailedTestsInDirChanges(dir string, rebuild bool) ([]provider.TestDiff, error) {
	return nil, errors.New("not implemented")
}

func (r *Runner) getPrevResults() (*results.TestResults, error) {
	if r.prev != nil {
		return r.prev, nil
	}

	prev, err := ci.LoadPrevCiGithub()

	if err != nil {
		return nil, err
	}

	r.prev = prev

	return prev, nil
}

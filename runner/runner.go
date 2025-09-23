package runner

import (
	"github.com/Sharktheone/mcp262/provider"
)

type Runner struct {
}

func (r Runner) RerunTest(testPath string, rebuild bool) (provider.TestResult, error) {
	//TODO implement me
	panic("implement me")
}

func (r Runner) RerunTestsInDir(dir string, rebuild bool) (map[string]provider.TestResult, error) {
	//TODO implement me
	panic("implement me")
}

func (r Runner) RerunFailedTestsInDir(dir string, rebuild bool) (map[string]provider.TestResult, error) {
	//TODO implement me
	panic("implement me")
}

package provider

type TestRunner interface {
	RerunTest(testPath string, rebuild bool) (TestResult, error)
	RerunTestsInDir(dir string, rebuild bool) (map[string]TestResult, error)
	RerunFailedTestsInDir(dir string, rebuild bool) (map[string]TestResult, error)
}

type TestResult struct {
	TestPath string `json:"test_path"`
	Status   string `json:"status"`
	Output   string `json:"output"`
	Duration string `json:"duration"`
}

var Runner TestRunner

func SetRunner(r TestRunner) {
	Runner = r
}

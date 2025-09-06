package provider

type TestProvider interface {
	NumTests() int
	NumTestsInDir(dir string) (int, error)
	NumTestInDirRec(dir string) (int, error)

	GetTestStatus(testPath string) (string, error)
	GetTestStatusesInDir(dir string) (map[string]string, error)
	GetTestsWithStatusInDir(dir string, status string) ([]string, error)
	GetFailedTestsInDir(dir string) ([]string, error)

	GetTestOutput(testPath string) (string, error)
}

var Provider TestProvider

func SetProvider(p TestProvider) {
	Provider = p
}

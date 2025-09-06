package provider

type TestProvider interface {
	NumTests() int
	NumTestsInDir(dir string) (int, error)
	NumTestsInDirRec(dir string) (int, error)

	GetTestsInDir(dir string) ([]string, error)
	GetTestsInDirRec(dir string) ([]string, error)

	GetTestStatus(testPath string) (string, error)
	GetTestStatusesInDir(dir string) (map[string]string, error)
	GetTestStatusesInDirRec(dir string) (map[string]string, error)

	GetTestsWithStatusInDir(dir string, status string) ([]string, error)
	GetTestsWithStatusInDirRec(dir string, status string) ([]string, error)

	GetFailedTestsInDir(dir string) ([]string, error)
	GetFailedTestsInDirRec(dir string) ([]string, error)

	SearchDir(query string) ([]string, error)
	SearchDirIn(dir string, query string) ([]string, error)

	SearchTest(query string) ([]string, error)
	SearchTestInDir(dir string, query string) ([]string, error)

	GetTestOutput(testPath string) (string, string, error)
}

var Provider TestProvider

func SetProvider(p TestProvider) {
	Provider = p
}

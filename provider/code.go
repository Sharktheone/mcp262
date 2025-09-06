package provider

type TestCodeProvider interface {
	GetTestCode(testPath string) (string, error)
	GetHarnessForTest(testPath string) (map[string]string, error)
	GetHarness() (map[string]string, error)
	GetHarnessCode(filePath string) (string, error)
	GetHaressFiles() ([]string, error)
	GetHarnessFilesForTest(testPath string) ([]string, error)

	SetTestCode(testPath string, code string) error
	SetHarnessCode(filePath string, code string) error

	ResetEdits() error
}

var CodeProvider TestCodeProvider

func SetCodeProvider(p TestCodeProvider) {
	CodeProvider = p
}

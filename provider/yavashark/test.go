package yavashark

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"slices"

	"github.com/Sharktheone/mcp262/testtree"
)

const RESULTS_URL = "https://raw.githubusercontent.com/Sharktheone/yavashark-data/refs/heads/main/results.json"
const BASE_RESULT_URL = "https://raw.githubusercontent.com/Sharktheone/yavashark-data/refs/heads/main/results"

var FailedStatuses = []string{"FAIL", "TIMEOUT", "CRASH", "NOT_IMPLEMENTED", "RUNNER_ERROR", "ERROR"}

type YavasharkResult struct {
	Status   string `json:"status"`
	Message  string `json:"msg"`
	Path     string `json:"path"`
	MemoryKB int    `json:"memory_kb"`
	Duration int    `json:"duration"`
}

type YavasharkTestProvider struct {
	*testtree.TestTree
}

type YavasharkTestResult struct {
	Path   string `json:"p"`
	Status string `json:"s"`
}

func NewYavasharkTestProvider() (*YavasharkTestProvider, error) {
	tree := testtree.NewTestTree()

	// Load results from the RESULTS_URL

	res, err := http.Get(RESULTS_URL)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, err
	}

	var results []YavasharkTestResult
	err = json.NewDecoder(res.Body).Decode(&results)

	if err != nil {
		return nil, err
	}

	for _, result := range results {
		tree.AddFile(result.Path, expandShortStatus(result.Status))
	}

	return &YavasharkTestProvider{
		tree,
	}, nil
}

func (yt *YavasharkTestProvider) GetFailedTestsInDir(dir string) ([]string, error) {
	if d, exists := yt.Directories[dir]; exists {
		out := make([]string, 0)
		for p, f := range d.Files {
			if slices.Contains(FailedStatuses, f.Status) {
				out = append(out, p)
			}
		}
		return out, nil
	}
	return nil, errors.New("directory not found")
}

func (yt *YavasharkTestProvider) GetFailedTestsInDirRec(dir string) ([]string, error) {
	if d, exists := yt.Directories[dir]; exists {
		out := make([]string, 0)
		var walk func(*testtree.TestTreeDir)
		walk = func(td *testtree.TestTreeDir) {
			for p, f := range td.Files {
				if slices.Contains(FailedStatuses, f.Status) {
					out = append(out, p)
				}
			}
			for _, sub := range td.Directories {
				walk(sub)
			}
		}
		walk(d)
		return out, nil
	}
	return nil, errors.New("directory not found")
}

func (yt *YavasharkTestProvider) GetTestOutput(testPath string) (string, string, error) {
	urlString, err := url.JoinPath(BASE_RESULT_URL, testPath+".json")
	if err != nil {
		return "", "", err
	}

	res, err := http.Get(urlString)
	if err != nil {
		return "", "", err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", "", errors.New("test output not found")
	}

	var result YavasharkResult
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return "", "", err
	}

	return result.Message, result.Status, nil
}

func expandShortStatus(status string) string {
	switch status {
	case "P":
		return "PASS"
	case "F":
		return "FAIL"
	case "S":
		return "SKIP"
	case "T":
		return "TIMEOUT"
	case "C":
		return "CRASH"
	case "O":
		return "PARSE_ERROR"
	case "PF":
		return "NOT_IMPLEMENTED"
	case "N":
		return "RUNNER_ERROR"
	default:
		return "ERROR"
	}
}

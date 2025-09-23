package ci

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/Sharktheone/mcp262/runner/results"
)

const CI_GITHUB_URL = "https://raw.githubusercontent.com/Sharktheone/yavashark-data/main/results.json"

func LoadPrevCi(path string) (*results.TestResults, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}

	var resultsCI []results.CIResult
	err = json.Unmarshal(contents, &resultsCI)

	if err != nil {
		return nil, err
	}

	res := results.ConvertResultsFromCI(resultsCI)

	return results.FromResults(res), nil
}

func LoadPrevCiGithub() (*results.TestResults, error) {
	res, err := http.Get(CI_GITHUB_URL)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, err
	}

	var resultsCI []results.CIResult
	err = json.NewDecoder(res.Body).Decode(&resultsCI)

	if err != nil {
		return nil, err
	}

	r := results.ConvertResultsFromCI(resultsCI)

	return results.FromResults(r), nil
}

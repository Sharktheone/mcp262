package github

import (
	"errors"
	"io"
	"net/http"
	"net/url"
)

const BaseUrl = "https://raw.githubusercontent.com/tc39/test262/refs/heads/main/"

var HarnessFiles = []string{
	"assert.js",
	"sta.js",
}

type GithubTest262CodeProvider struct{}

func (g GithubTest262CodeProvider) GetTestCode(testPath string) (string, error) {
	codeUrl, err := url.JoinPath(BaseUrl, "test", testPath)
	if err != nil {
		return "", err
	}

	return fetchCodeFromUrl(codeUrl)
}

func (g GithubTest262CodeProvider) GetHarnessForTest(testPath string) (map[string]string, error) {
	return nil, errors.New("not implemented")
}

func (g GithubTest262CodeProvider) GetHarness() (map[string]string, error) {
	harnessCode := make(map[string]string, 2)
	for _, file := range HarnessFiles {
		codeUrl, err := url.JoinPath(BaseUrl, "harness", file)
		if err != nil {
			return nil, err
		}
		code, err := fetchCodeFromUrl(codeUrl)
		if err != nil {
			return nil, err
		}
		harnessCode[file] = code
	}

	return harnessCode, nil
}

func (g GithubTest262CodeProvider) GetHaressFiles() ([]string, error) {
	return HarnessFiles, nil
}

func (g GithubTest262CodeProvider) GetHarnessCode(filePath string) (string, error) {
	codeUrl, err := url.JoinPath(BaseUrl, "harness", filePath)
	if err != nil {
		return "", err
	}

	return fetchCodeFromUrl(codeUrl)
}

func (g GithubTest262CodeProvider) GetHarnessFilesForTest(testPath string) ([]string, error) {
	return nil, errors.New("not implemented")
}

func (g GithubTest262CodeProvider) SetTestCode(testPath string, code string) error {
	return errors.New("SetTestCode not available with GithubTest262CodeProvider")
}

func (g GithubTest262CodeProvider) SetHarnessCode(filePath string, code string) error {
	return errors.New("SetHarnessCode not available with GithubTest262CodeProvider")
}

func (g GithubTest262CodeProvider) ResetEdits() error {
	return nil
}

func NewGithubTest262CodeProvider() GithubTest262CodeProvider {
	return GithubTest262CodeProvider{}
}

func fetchCodeFromUrl(codeUrl string) (string, error) {
	res, err := http.Get(codeUrl)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", err
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

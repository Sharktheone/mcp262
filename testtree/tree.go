package testtree

import (
	"errors"
	"path"
	"strings"
)

type TestTree struct {
	Files       map[string]*TestTreeFile
	Directories map[string]*TestTreeDir
}

func (tt *TestTree) NumTests() int {
	return len(tt.Files)
}

func (tt *TestTree) NumTestsInDir(dir string) (int, error) {
	d := normalizeDir(dir)
	if d, exists := tt.Directories[d]; exists {
		return len(d.Files), nil
	}

	return 0, errors.New("directory not found")
}

func (tt *TestTree) NumTestsInDirRec(dir string) (int, error) {
	d := normalizeDir(dir)
	if d, exists := tt.Directories[d]; exists {
		count := 0
		var walk func(*TestTreeDir)
		walk = func(td *TestTreeDir) {
			count += len(td.Files)
			for _, sub := range td.Directories {
				walk(sub)
			}
		}
		walk(d)
		return count, nil
	}

	return 0, errors.New("directory not found")
}

func (tt *TestTree) GetTestsInDir(dir string) ([]string, error) {
	d := normalizeDir(dir)
	if d, exists := tt.Directories[d]; exists {
		out := make([]string, 0)
		for p, _ := range d.Files {
			out = append(out, p)
		}

		for p, _ := range d.Directories {
			out = append(out, p)
		}

		return out, nil
	}
	return nil, errors.New("directory not found")
}

func (tt *TestTree) GetTestsInDirRec(dir string) ([]string, error) {
	d := normalizeDir(dir)
	if d, exists := tt.Directories[d]; exists {
		out := make([]string, 0)
		var walk func(*TestTreeDir)
		walk = func(td *TestTreeDir) {
			for p, _ := range td.Files {
				out = append(out, p)
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

func (tt *TestTree) GetTestStatus(testPath string) (string, error) {
	if f, exists := tt.Files[testPath]; exists {
		return f.Status, nil
	}
	return "", errors.New("test not found")
}

func (tt *TestTree) GetTestStatusesInDir(dir string) (map[string]string, error) {
	d := normalizeDir(dir)
	if d, exists := tt.Directories[d]; exists {
		res := make(map[string]string, len(d.Files))
		for p, f := range d.Files {
			res[p] = f.Status
		}
		return res, nil
	}
	return nil, errors.New("directory not found")
}

func (tt *TestTree) GetTestStatusesInDirRec(dir string) (map[string]string, error) {
	d := normalizeDir(dir)
	if d, exists := tt.Directories[d]; exists {
		res := make(map[string]string)
		var walk func(*TestTreeDir)
		walk = func(td *TestTreeDir) {
			for p, f := range td.Files {
				res[p] = f.Status
			}
			for _, sub := range td.Directories {
				walk(sub)
			}
		}
		walk(d)
		return res, nil
	}
	return nil, errors.New("directory not found")
}

func (tt *TestTree) GetTestsWithStatusInDir(dir string, status string) ([]string, error) {
	d := normalizeDir(dir)
	if d, exists := tt.Directories[d]; exists {
		out := make([]string, 0)
		for p, f := range d.Files {
			if strings.EqualFold(f.Status, status) {
				out = append(out, p)
			}
		}
		return out, nil
	}
	return nil, errors.New("directory not found")
}

func (tt *TestTree) GetTestsWithStatusInDirRec(dir string, status string) ([]string, error) {
	d := normalizeDir(dir)
	if d, exists := tt.Directories[d]; exists {
		out := make([]string, 0)
		var walk func(*TestTreeDir)
		walk = func(td *TestTreeDir) {
			for p, f := range td.Files {
				if strings.EqualFold(f.Status, status) {
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

type TestTreeFile struct {
	Path   string
	Status string
}

type TestTreeDir struct {
	Path        string
	Files       map[string]*TestTreeFile
	Directories map[string]*TestTreeDir
}

func NewTestTree() *TestTree {
	return &TestTree{
		Files:       make(map[string]*TestTreeFile),
		Directories: make(map[string]*TestTreeDir),
	}
}

func NewTestTreeSize(files int, dirs int) *TestTree {
	return &TestTree{
		Files:       make(map[string]*TestTreeFile, files),
		Directories: make(map[string]*TestTreeDir, dirs),
	}
}

// normalizeDir cleans a directory path and strips any trailing slash.
// It returns an empty string for the repository root.
func normalizeDir(p string) string {
	if p == "" {
		return ""
	}
	// Use path.Clean to remove trailing slashes and normalize
	cleaned := path.Clean(p)
	if cleaned == "." || cleaned == "/" {
		return ""
	}
	return cleaned
}

func (tt *TestTree) AddFile(p string, status string) {
	f := &TestTreeFile{
		Path:   p,
		Status: status,
	}

	tt.Files[p] = f

	// Determine directory for the file and normalize it
	dirPath := normalizeDir(path.Dir(p))
	tt.AddDir(dirPath)

	dir := tt.Directories[dirPath]
	dir.Files[p] = f
}

func (tt *TestTree) AddDir(p string) {
	np := normalizeDir(p)
	if _, exists := tt.Directories[np]; !exists {
		dir := &TestTreeDir{
			Path:        np,
			Files:       make(map[string]*TestTreeFile),
			Directories: make(map[string]*TestTreeDir),
		}

		tt := tt
		tt.Directories[np] = dir

		// If this is not the root, attach it to its parent (which will also be created)
		if np != "" {
			parent := normalizeDir(path.Dir(np))
			if parent == np {
				// shouldn't happen, but guard against infinite recursion
				parent = ""
			}
			tt.AddDir(parent)
			parentDir := tt.Directories[parent]
			parentDir.Directories[np] = dir
		}
	}
}

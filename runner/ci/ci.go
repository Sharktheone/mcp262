package ci

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Sharktheone/mcp262/runner/results"
)

func computeAggregate(dir string, summaries map[string]*DirectorySummary) DirectorySummary {
	base, exists := summaries[dir]
	var agg DirectorySummary
	if exists {
		agg = *base
	} else {
		agg = DirectorySummary{Directory: dir}
	}

	for k, ds := range summaries {
		if k == dir {
			continue
		}
		if dir == "" {
			if k != "" {
				agg.Passed += ds.Passed
				agg.Failed += ds.Failed
				agg.Skipped += ds.Skipped
				agg.NotImplemented += ds.NotImplemented
				agg.RunnerError += ds.RunnerError
				agg.Crashed += ds.Crashed
				agg.Timeout += ds.Timeout
				agg.ParseError += ds.ParseError
				agg.Total += ds.Total
			}
		} else {
			prefix := dir + string(filepath.Separator)
			if strings.HasPrefix(k, prefix) {
				agg.Passed += ds.Passed
				agg.Failed += ds.Failed
				agg.Skipped += ds.Skipped
				agg.NotImplemented += ds.NotImplemented
				agg.RunnerError += ds.RunnerError
				agg.Crashed += ds.Crashed
				agg.Timeout += ds.Timeout
				agg.ParseError += ds.ParseError
				agg.Total += ds.Total
			}
		}
	}
	return agg
}

func fmtCiDiff(path string, testResults *results.TestResults, root string) (string, error) {
	var sb strings.Builder

	prev, _ := LoadPrevCi(path)
	if prev != nil {
		d := testResults.ComputeDiffRoot(prev, root)

		_, err := fmt.Fprintln(&sb, "<=== DIFF ===>")
		if err != nil {
			return "", err
		}
		d.FmtGrouped(&sb)

		_, err = fmt.Fprintln(&sb)
		if err != nil {
			return "", err
		}
		_, err = fmt.Fprintln(&sb)
		if err != nil {
			return "", err
		}
		_, err = fmt.Fprintln(&sb)
		if err != nil {
			return "", err
		}
		_, err = fmt.Fprintln(&sb, "<=== Test Results ===>")
		if err != nil {
			return "", err
		}
		_, err = fmt.Fprintln(&sb)
		if err != nil {
			return "", err
		}

		testResults.FmtResults(&sb)

		_, err = fmt.Fprintln(&sb)
		if err != nil {
			return "", err
		}
		_, err = fmt.Fprintln(&sb, "<=== Comparison ===>")
		if err != nil {
			return "", err
		}

		testResults.Compare(prev)
	}

	return sb.String(), nil
}

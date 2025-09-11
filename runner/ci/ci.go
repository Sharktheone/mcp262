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

func printCiDiff(path string, testResults *results.TestResults, root string) bool {
	prev, _ := LoadPrevCi(path)
	if prev != nil {
		d := testResults.ComputeDiffRoot(prev, root)

		fmt.Println("<=== DIFF ===>")
		d.PrintGrouped()

		fmt.Println()
		fmt.Println()
		fmt.Println("<=== Test Results ===>")

		testResults.PrintResults()

		fmt.Println()
		fmt.Println("<=== Comparison ===>")

		testResults.Compare(prev)

		return d.NewTestFailures() > MAX_NEW_TEST_FAILURES
	}

	return false
}

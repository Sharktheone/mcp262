package tools

import (
	"context"
	"errors"

	"github.com/Sharktheone/mcp262/provider"
	"github.com/Sharktheone/mcp262/utils"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type RerunTestParams struct {
	TestPath string `json:"test_path" jsonschema:"Path to the single test file (e.g. /test262/test/language/...)"`
	Rebuild  bool   `json:"rebuild" jsonschema:"Whether to rebuild before running the test"`
}

type RerunTestsInDirParams struct {
	Dir     string `json:"dir" jsonschema:"Directory path to run tests in"`
	Rebuild bool   `json:"rebuild" jsonschema:"Whether to rebuild before running the tests"`
}

type RerunFailedTestsInDirParams struct {
	Dir     string `json:"dir" jsonschema:"Directory path to run failed tests in"`
	Rebuild bool   `json:"rebuild" jsonschema:"Whether to rebuild before running the tests"`
}

func RerunTest(ctx context.Context, req *mcp.CallToolRequest, args RerunTestParams) (*mcp.CallToolResult, any, error) {
	runner, err := getRunner()
	if err != nil {
		return nil, nil, err
	}
	p := utils.ResolvePath(args.TestPath)
	result, err := runner.RerunTest(p, args.Rebuild)
	if err != nil {
		return nil, nil, err
	}
	return utils.RespondWith(map[string]any{
		"test_result": result,
	}), nil, nil
}

func RerunTestsInDir(ctx context.Context, req *mcp.CallToolRequest, args RerunTestsInDirParams) (*mcp.CallToolResult, any, error) {
	runner, err := getRunner()
	if err != nil {
		return nil, nil, err
	}
	p := utils.ResolvePath(args.Dir)
	results, err := runner.RerunTestsInDirChanges(p, args.Rebuild)
	if err != nil {
		return nil, nil, err
	}

	return utils.RespondWith(map[string]any{
		"results": results,
	}), nil, nil
}

func RerunFailedTestsInDir(ctx context.Context, req *mcp.CallToolRequest, args RerunFailedTestsInDirParams) (*mcp.CallToolResult, any, error) {
	runner, err := getRunner()
	if err != nil {
		return nil, nil, err
	}
	p := utils.ResolvePath(args.Dir)
	results, err := runner.RerunFailedTestsInDirChanges(p, args.Rebuild)
	if err != nil {
		return nil, nil, err
	}
	return utils.RespondWith(map[string]any{
		"results": results,
	}), nil, nil
}

func AddRunnerTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "RerunTest",
		Description: "Rerun a single test",
	}, RerunTest)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "RerunTestsInDir",
		Description: "Rerun all tests in a directory",
	}, RerunTestsInDir)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "RerunFailedTestsInDir",
		Description: "Rerun failed tests in a directory",
	}, RerunFailedTestsInDir)
}

// helper to validate runner is set
func getRunner() (provider.TestRunner, error) {
	if provider.Runner == nil {
		return nil, errors.New("runner not set")
	}
	return provider.Runner, nil
}

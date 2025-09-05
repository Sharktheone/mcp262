package test

import (
	"context"

	"github.com/Sharktheone/mcp262/utils"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type NumTestParams struct {
	Path string `json:"path" jsonschema:"Path of the directory, starting from /test262/test/{builtins,language,...}/..., /test/{builtins,language,...}/... or just /{builtins,language,...}/..."`
}

func NumTests(ctx context.Context, req *mcp.CallToolRequest, args NumTestParams) (*mcp.CallToolResult, any, error) {
	n, err := getNumTests(args.Path)
	if err != nil {
		return nil, nil, err
	}

	res := map[string]any{
		"num_tests": n,
		"path":      args.Path,
	}

	return utils.RespondWith(res), nil, nil
}

func AddTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "num_tests",
		Description: "Get the number of tests in a given directory",
	}, NumTests)
}

func getNumTests(path string) (int, error) {
	return 42, nil
}

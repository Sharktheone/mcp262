package test

import (
	"context"
	"errors"
	"strings"

	"github.com/Sharktheone/mcp262/provider"
	"github.com/Sharktheone/mcp262/utils"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type GetTestCodeParams struct {
	TestPath string `json:"test_path" jsonschema:"Path to the single test file (e.g. /test262/test/language/...)"`
}

type GetHarnessForTestParams struct {
	TestPath string `json:"test_path" jsonschema:"Path to the single test file (e.g. /test262/test/language/...)"`
}

type GetHarnessCodeParams struct {
	HarnessPath string `json:"harness_path" jsonschema:"Path to the single harness file (e.g. sta.js, assert.js, etc.)"`
}

type GetHarnessFilesForTestParams struct {
	TestPath string `json:"test_path" jsonschema:"Path to the single test file (e.g. /test262/test/language/...)"`
}

type SetTestCodeParams struct {
	TestPath string `json:"test_path" jsonschema:"Path to the single test file (e.g. /test262/test/language/...)"`
	Code     string `json:"code" jsonschema:"New code for the test"`
}

type SetHarnessCodeParams struct {
	FilePath string `json:"file_path" jsonschema:"Path to the harness file (relative to harness root)"`
	Code     string `json:"code" jsonschema:"New code for the harness file"`
}

func GetTestCode(ctx context.Context, req *mcp.CallToolRequest, args GetTestCodeParams) (*mcp.CallToolResult, any, error) {
	pv, err := getCodeProvider()
	if err != nil {
		return nil, nil, err
	}
	p := utils.ResolvePath(args.TestPath)
	if !strings.HasSuffix(p, ".js") {
		p += ".js"
	}
	code, err := pv.GetTestCode(p)
	if err != nil {
		return nil, nil, err
	}
	return utils.RespondWith(map[string]any{"test_path": args.TestPath, "code": code}), nil, nil
}

func GetHarnessForTest(ctx context.Context, req *mcp.CallToolRequest, args GetHarnessForTestParams) (*mcp.CallToolResult, any, error) {
	pv, err := getCodeProvider()
	if err != nil {
		return nil, nil, err
	}
	p := utils.ResolvePath(args.TestPath)
	h, err := pv.GetHarnessForTest(p)
	if err != nil {
		return nil, nil, err
	}
	return utils.RespondWith(map[string]any{"test_path": args.TestPath, "harness": h}), nil, nil
}

func GetHarness(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
	pv, err := getCodeProvider()
	if err != nil {
		return nil, nil, err
	}
	h, err := pv.GetHarness()
	if err != nil {
		return nil, nil, err
	}
	return utils.RespondWith(map[string]any{"harness": h}), nil, nil
}

func GetHarnessCode(ctx context.Context, req *mcp.CallToolRequest, args GetHarnessCodeParams) (*mcp.CallToolResult, any, error) {
	pv, err := getCodeProvider()
	if err != nil {
		return nil, nil, err
	}
	p := utils.ResolvePath(args.HarnessPath)
	if !strings.HasSuffix(p, ".js") {
		p += ".js"
	}
	code, err := pv.GetHarnessCode(p)
	if err != nil {
		return nil, nil, err
	}
	return utils.RespondWith(map[string]any{"harness_path": args.HarnessPath, "code": code}), nil, nil
}

func GetHaressFiles(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
	pv, err := getCodeProvider()
	if err != nil {
		return nil, nil, err
	}
	files, err := pv.GetHaressFiles()
	if err != nil {
		return nil, nil, err
	}
	return utils.RespondWith(map[string]any{"files": files}), nil, nil
}

func GetHarnessFilesForTest(ctx context.Context, req *mcp.CallToolRequest, args GetHarnessFilesForTestParams) (*mcp.CallToolResult, any, error) {
	pv, err := getCodeProvider()
	if err != nil {
		return nil, nil, err
	}
	p := utils.ResolvePath(args.TestPath)
	files, err := pv.GetHarnessFilesForTest(p)
	if err != nil {
		return nil, nil, err
	}
	return utils.RespondWith(map[string]any{"test_path": args.TestPath, "files": files}), nil, nil
}

func SetTestCode(ctx context.Context, req *mcp.CallToolRequest, args SetTestCodeParams) (*mcp.CallToolResult, any, error) {
	pv, err := getCodeProvider()
	if err != nil {
		return nil, nil, err
	}
	p := utils.ResolvePath(args.TestPath)
	if !strings.HasSuffix(p, ".js") {
		p += ".js"
	}
	if err := pv.SetTestCode(p, args.Code); err != nil {
		return nil, nil, err
	}
	return utils.RespondWith(map[string]any{"test_path": args.TestPath, "updated": true}), nil, nil
}

func SetHarnessCode(ctx context.Context, req *mcp.CallToolRequest, args SetHarnessCodeParams) (*mcp.CallToolResult, any, error) {
	pv, err := getCodeProvider()
	if err != nil {
		return nil, nil, err
	}
	p := utils.ResolvePath(args.FilePath)
	if err := pv.SetHarnessCode(p, args.Code); err != nil {
		return nil, nil, err
	}
	return utils.RespondWith(map[string]any{"file_path": args.FilePath, "updated": true}), nil, nil
}

func ResetEdits(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
	pv, err := getCodeProvider()
	if err != nil {
		return nil, nil, err
	}
	if err := pv.ResetEdits(); err != nil {
		return nil, nil, err
	}
	return utils.RespondWith(map[string]any{"reset": true}), nil, nil
}

func AddCodeTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "GetTestCode",
		Description: "Get the source code for a single test",
	}, GetTestCode)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "GetHarnessForTest",
		Description: "Get harness files (map path->code) required by a test",
	}, GetHarnessForTest)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "GetHarness",
		Description: "Get all harness files (map path->code)",
	}, GetHarness)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "GetHarnessCode",
		Description: "Get the source code for a single harness file",
	}, GetHarnessCode)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "GetHaressFiles",
		Description: "List harness file paths ",
	}, GetHaressFiles)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "GetHarnessFilesForTest",
		Description: "List harness files used by a single test",
	}, GetHarnessFilesForTest)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "SetTestCode",
		Description: "Replace the source code for a single test",
	}, SetTestCode)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "SetHarnessCode",
		Description: "Replace the source code for a harness file",
	}, SetHarnessCode)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "ResetEdits",
		Description: "Reset any in-memory edits to tests/harness",
	}, ResetEdits)
}

// helper to validate provider is set
func getCodeProvider() (provider.TestCodeProvider, error) {
	if provider.CodeProvider == nil {
		return nil, errors.New("code provider not set")
	}
	return provider.CodeProvider, nil
}

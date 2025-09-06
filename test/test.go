package test

import (
	"context"
	"errors"
	"sort"

	"github.com/Sharktheone/mcp262/provider"
	"github.com/Sharktheone/mcp262/utils"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const DefaultPageSize = 20

type NumTestsRecursiveParams struct {
	Path string `json:"path" jsonschema:"Path of the directory, starting from /test262/test/{builtins,language,...}/..., /test/{builtins,language,...}/... or just /{builtins,language,...}/..."`
}

type GetTestStatusParams struct {
	TestPath string `json:"test_path" jsonschema:"Path to the single test file (e.g. /test262/test/language/...)"`
}

type Pagination struct {
	Page     int `json:"page" jsonschema:"Page number starting from 1; defaults to 1"`
	PageSize int `json:"page_size" jsonschema:"Items per page; defaults to DefaultPageSize if omitted"`
	Max      int `json:"max" jsonschema:"Optional global maximum number of items to include across all pages; 0 means no limit"`
}

type GetStatusesInDirParams struct {
	Path string `json:"path" jsonschema:"Directory path to list statuses for"`
	Pagination
}

type GetTestsWithStatusInDirParams struct {
	Path   string `json:"path" jsonschema:"Directory path to filter tests by status"`
	Status string `json:"status" jsonschema:"Status to filter by (e.g. PASS, FAIL, SKIP, TIMEOUT, CRASH, PARSE_ERROR, NOT_IMPLEMENTED, RUNNER_ERROR)"`
	Pagination
}

type GetFailedTestsInDirParams struct {
	Path string `json:"path" jsonschema:"Directory path to list failed tests for"`
	Pagination
}

type GetTestOutputParams struct {
	TestPath string `json:"test_path" jsonschema:"Path to the single test file (e.g. /test262/test/language/...)"`
}

// Tool handlers

func NumTestsTotal(ctx context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
	prov, err := getProvider()
	if err != nil {
		return nil, nil, err
	}
	n := prov.NumTests()
	return utils.RespondWith(map[string]any{"num_tests": n}), nil, nil
}

func NumTestsInDirRecursive(ctx context.Context, req *mcp.CallToolRequest, args NumTestsRecursiveParams) (*mcp.CallToolResult, any, error) {
	prov, err := getProvider()
	if err != nil {
		return nil, nil, err
	}
	p := utils.ResolvePath(args.Path)
	n, err := prov.NumTestInDirRec(p)
	if err != nil {
		return nil, nil, err
	}
	return utils.RespondWith(map[string]any{"num_tests": n, "path": args.Path, "recursive": true}), nil, nil
}

func GetTestStatus(ctx context.Context, req *mcp.CallToolRequest, args GetTestStatusParams) (*mcp.CallToolResult, any, error) {
	prov, err := getProvider()
	if err != nil {
		return nil, nil, err
	}
	p := utils.ResolvePath(args.TestPath)
	status, err := prov.GetTestStatus(p)
	if err != nil {
		return nil, nil, err
	}
	return utils.RespondWith(map[string]any{"test_path": args.TestPath, "status": status}), nil, nil
}

func GetTestStatusesInDir(ctx context.Context, req *mcp.CallToolRequest, args GetStatusesInDirParams) (*mcp.CallToolResult, any, error) {
	prov, err := getProvider()
	if err != nil {
		return nil, nil, err
	}
	p := utils.ResolvePath(args.Path)
	statuses, err := prov.GetTestStatusesInDir(p)
	if err != nil {
		return nil, nil, err
	}
	// stable order by test path
	keys := make([]string, 0, len(statuses))
	for k := range statuses {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	page, pageSize := normalizePage(args.Page, args.PageSize)
	items, remaining, total := paginateStrings(keys, page, pageSize, args.Max)
	paged := make(map[string]string, len(items))
	for _, k := range items {
		paged[k] = statuses[k]
	}
	res := map[string]any{
		"path":      args.Path,
		"page":      page,
		"page_size": pageSize,
		"returned":  len(items),
		"remaining": remaining,
		"total":     total,
		"statuses":  paged,
	}
	return utils.RespondWith(res), nil, nil
}

func GetTestStatusesInDirRec(ctx context.Context, req *mcp.CallToolRequest, args GetStatusesInDirParams) (*mcp.CallToolResult, any, error) {
	prov, err := getProvider()
	if err != nil {
		return nil, nil, err
	}
	p := utils.ResolvePath(args.Path)
	statuses, err := prov.GetTestStatusesInDirRec(p)
	if err != nil {
		return nil, nil, err
	}
	// stable order by test path
	keys := make([]string, 0, len(statuses))
	for k := range statuses {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	page, pageSize := normalizePage(args.Page, args.PageSize)
	items, remaining, total := paginateStrings(keys, page, pageSize, args.Max)
	paged := make(map[string]string, len(items))
	for _, k := range items {
		paged[k] = statuses[k]
	}
	res := map[string]any{
		"path":      args.Path,
		"page":      page,
		"page_size": pageSize,
		"returned":  len(items),
		"remaining": remaining,
		"total":     total,
		"statuses":  paged,
	}
	return utils.RespondWith(res), nil, nil
}

func GetTestsWithStatusInDir(ctx context.Context, req *mcp.CallToolRequest, args GetTestsWithStatusInDirParams) (*mcp.CallToolResult, any, error) {
	prov, err := getProvider()
	if err != nil {
		return nil, nil, err
	}
	p := utils.ResolvePath(args.Path)
	tests, err := prov.GetTestsWithStatusInDir(p, args.Status)
	if err != nil {
		return nil, nil, err
	}
	sort.Strings(tests)
	page, pageSize := normalizePage(args.Page, args.PageSize)
	items, remaining, total := paginateStrings(tests, page, pageSize, args.Max)
	res := map[string]any{
		"path":      args.Path,
		"status":    args.Status,
		"page":      page,
		"page_size": pageSize,
		"returned":  len(items),
		"remaining": remaining,
		"total":     total,
		"tests":     items,
	}
	return utils.RespondWith(res), nil, nil
}

func GetTestsWithStatusInDirRec(ctx context.Context, req *mcp.CallToolRequest, args GetTestsWithStatusInDirParams) (*mcp.CallToolResult, any, error) {
	prov, err := getProvider()
	if err != nil {
		return nil, nil, err
	}
	p := utils.ResolvePath(args.Path)
	tests, err := prov.GetTestsWithStatusInDirRec(p, args.Status)
	if err != nil {
		return nil, nil, err
	}
	sort.Strings(tests)
	page, pageSize := normalizePage(args.Page, args.PageSize)
	items, remaining, total := paginateStrings(tests, page, pageSize, args.Max)
	res := map[string]any{
		"path":      args.Path,
		"status":    args.Status,
		"page":      page,
		"page_size": pageSize,
		"returned":  len(items),
		"remaining": remaining,
		"total":     total,
		"tests":     items,
	}
	return utils.RespondWith(res), nil, nil
}

func GetFailedTestsInDir(ctx context.Context, req *mcp.CallToolRequest, args GetFailedTestsInDirParams) (*mcp.CallToolResult, any, error) {
	prov, err := getProvider()
	if err != nil {
		return nil, nil, err
	}
	p := utils.ResolvePath(args.Path)
	tests, err := prov.GetFailedTestsInDir(p)
	if err != nil {
		return nil, nil, err
	}
	sort.Strings(tests)
	page, pageSize := normalizePage(args.Page, args.PageSize)
	items, remaining, total := paginateStrings(tests, page, pageSize, args.Max)
	res := map[string]any{
		"path":      args.Path,
		"page":      page,
		"page_size": pageSize,
		"returned":  len(items),
		"remaining": remaining,
		"total":     total,
		"tests":     items,
	}
	return utils.RespondWith(res), nil, nil
}

func GetFailedTestsInDirRec(ctx context.Context, req *mcp.CallToolRequest, args GetFailedTestsInDirParams) (*mcp.CallToolResult, any, error) {
	prov, err := getProvider()
	if err != nil {
		return nil, nil, err
	}
	p := utils.ResolvePath(args.Path)
	tests, err := prov.GetFailedTestsInDir(p)
	if err != nil {
		return nil, nil, err
	}
	sort.Strings(tests)
	page, pageSize := normalizePage(args.Page, args.PageSize)
	items, remaining, total := paginateStrings(tests, page, pageSize, args.Max)
	res := map[string]any{
		"path":      args.Path,
		"page":      page,
		"page_size": pageSize,
		"returned":  len(items),
		"remaining": remaining,
		"total":     total,
		"tests":     items,
	}
	return utils.RespondWith(res), nil, nil
}

func GetTestOutput(ctx context.Context, req *mcp.CallToolRequest, args GetTestOutputParams) (*mcp.CallToolResult, any, error) {
	prov, err := getProvider()
	if err != nil {
		return nil, nil, err
	}
	p := utils.ResolvePath(args.TestPath)
	out, err := prov.GetTestOutput(p)
	if err != nil {
		return nil, nil, err
	}
	return utils.RespondWith(map[string]any{"test_path": args.TestPath, "output": out}), nil, nil
}

func AddTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "NumTestsTotal",
		Description: "Get the total number of tests",
	}, NumTestsTotal)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "NumTestsInDirRecursive",
		Description: "Get the number of tests in a directory recursively",
	}, NumTestsInDirRecursive)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "GetTestStatus",
		Description: "Get the status of a single test",
	}, GetTestStatus)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "GetTestStatusesInDirRec",
		Description: "List statuses for tests in a directory recursively (paginated)",
	}, GetTestStatusesInDirRec)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "GetTestStatusesInDir",
		Description: "List statuses for tests in a directory (paginated)",
	}, GetTestStatusesInDir)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "GetTestsWithStatusInDirRec",
		Description: "List tests with a specific status in a directory recursively (paginated)",
	}, GetTestsWithStatusInDirRec)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "GetTestsWithStatusInDir",
		Description: "List tests with a specific status in a directory (paginated)",
	}, GetTestsWithStatusInDir)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "GetFailedTestsInDirRec",
		Description: "List failed tests in a directory recursively (paginated)",
	}, GetFailedTestsInDirRec)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "GetFailedTestsInDir",
		Description: "List failed tests in a directory (paginated)",
	}, GetFailedTestsInDir)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "GetTestsWithStatusInDirRec",
		Description: "List tests with a specific status in a directory recursively (paginated)",
	}, GetTestsWithStatusInDirRec)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "GetTestOutput",
		Description: "Get the output of a single test",
	}, GetTestOutput)
}

func getProvider() (provider.TestProvider, error) {
	if provider.Provider == nil {
		return nil, errors.New("test provider not set")
	}
	return provider.Provider, nil
}

func normalizePage(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	return page, pageSize
}

func paginateStrings(items []string, page, pageSize, max int) ([]string, int, int) {
	total := len(items)
	limit := total
	if max > 0 && max < limit {
		limit = max
	}
	start := (page - 1) * pageSize
	if start > limit {
		start = limit
	}
	end := start + pageSize
	if end > limit {
		end = limit
	}
	slice := items[start:end]
	remaining := limit - end
	return slice, remaining, total
}

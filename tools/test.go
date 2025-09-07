package tools

import (
	"context"
	"errors"
	"sort"
	"strings"

	"github.com/Sharktheone/mcp262/provider"
	"github.com/Sharktheone/mcp262/utils"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const DefaultPageSize = 20

type NumTestsRecursiveParams struct {
	Path string `json:"path" jsonschema:"Path of the directory, starting from /test262/test/{built-ins,language,...}/..., /test/{built-ins,language,...}/... or just /{built-ins,language,...}/..."`
}

type GetTestStatusParams struct {
	TestPath string `json:"test_path" jsonschema:"Path to the single test file (e.g. /test262/test/language/...)"`
}

//type Pagination struct {
//	Page     int `json:"page" jsonschema:"Page number starting from 1; defaults to 1"`
//	PageSize int `json:"page_size" jsonschema:"Items per page; defaults to DefaultPageSize if omitted"`
//	Max      int `json:"max" jsonschema:"Optional global maximum number of items to include across all pages; 0 means no limit"`
//}

type GetTestsInDirParams struct {
	Path     string `json:"path" jsonschema:"Directory path to filter tests by status"`
	Page     int    `json:"page" jsonschema:"Page number starting from 1; defaults to 1"`
	PageSize int    `json:"page_size" jsonschema:"Items per page; defaults to DefaultPageSize if omitted"`
	Max      int    `json:"max" jsonschema:"Optional global maximum number of items to include across all pages; 0 means no limit"`
}

type GetStatusesInDirParams struct {
	Path     string `json:"path" jsonschema:"Directory path to list statuses for"`
	Page     int    `json:"page" jsonschema:"Page number starting from 1; defaults to 1"`
	PageSize int    `json:"page_size" jsonschema:"Items per page; defaults to DefaultPageSize if omitted"`
	Max      int    `json:"max" jsonschema:"Optional global maximum number of items to include across all pages; 0 means no limit"`
}

type GetTestsWithStatusInDirParams struct {
	Path     string `json:"path" jsonschema:"Directory path to filter tests by status"`
	Status   string `json:"status" jsonschema:"Status to filter by (e.g. PASS, FAIL, SKIP, TIMEOUT, CRASH, PARSE_ERROR, NOT_IMPLEMENTED, RUNNER_ERROR)"`
	Page     int    `json:"page" jsonschema:"Page number starting from 1; defaults to 1"`
	PageSize int    `json:"page_size" jsonschema:"Items per page; defaults to DefaultPageSize if omitted"`
	Max      int    `json:"max" jsonschema:"Optional global maximum number of items to include across all pages; 0 means no limit"`
}

type GetFailedTestsInDirParams struct {
	Path     string `json:"path" jsonschema:"Directory path to list failed tests for"`
	Page     int    `json:"page" jsonschema:"Page number starting from 1; defaults to 1"`
	PageSize int    `json:"page_size" jsonschema:"Items per page; defaults to DefaultPageSize if omitted"`
	Max      int    `json:"max" jsonschema:"Optional global maximum number of items to include across all pages; 0 means no limit"`
}

type GetTestOutputParams struct {
	TestPath string `json:"test_path" jsonschema:"Path to the single test file (e.g. /test262/test/language/...)"`
}

// Added search parameter structs
type SearchDirParams struct {
	Query    string `json:"query" jsonschema:"Search query"`
	Page     int    `json:"page" jsonschema:"Page number starting from 1; defaults to 1"`
	PageSize int    `json:"page_size" jsonschema:"Items per page; defaults to DefaultPageSize if omitted"`
	Max      int    `json:"max" jsonschema:"Optional global maximum number of items to include across all pages; 0 means no limit"`
}

type SearchDirInParams struct {
	Dir      string `json:"dir" jsonschema:"Directory path to restrict the search to"`
	Query    string `json:"query" jsonschema:"Search query"`
	Page     int    `json:"page" jsonschema:"Page number starting from 1; defaults to 1"`
	PageSize int    `json:"page_size" jsonschema:"Items per page; defaults to DefaultPageSize if omitted"`
	Max      int    `json:"max" jsonschema:"Optional global maximum number of items to include across all pages; 0 means no limit"`
}

type SearchTestParams struct {
	Query    string `json:"query" jsonschema:"Search query for tests"`
	Page     int    `json:"page" jsonschema:"Page number starting from 1; defaults to 1"`
	PageSize int    `json:"page_size" jsonschema:"Items per page; defaults to DefaultPageSize if omitted"`
	Max      int    `json:"max" jsonschema:"Optional global maximum number of items to include across all pages; 0 means no limit"`
}

type SearchTestInDirParams struct {
	Dir      string `json:"dir" jsonschema:"Directory path to restrict the test search to"`
	Query    string `json:"query" jsonschema:"Search query for tests"`
	Page     int    `json:"page" jsonschema:"Page number starting from 1; defaults to 1"`
	PageSize int    `json:"page_size" jsonschema:"Items per page; defaults to DefaultPageSize if omitted"`
	Max      int    `json:"max" jsonschema:"Optional global maximum number of items to include across all pages; 0 means no limit"`
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

func NumTestsInDir(ctx context.Context, req *mcp.CallToolRequest, args NumTestsRecursiveParams) (*mcp.CallToolResult, any, error) {
	prov, err := getProvider()
	if err != nil {
		return nil, nil, err
	}
	p := utils.ResolvePath(args.Path)
	n, err := prov.NumTestsInDir(p)
	if err != nil {
		return nil, nil, err
	}
	return utils.RespondWith(map[string]any{"num_tests": n, "path": args.Path, "recursive": true}), nil, nil

}

func NumTestsInDirRecursive(ctx context.Context, req *mcp.CallToolRequest, args NumTestsRecursiveParams) (*mcp.CallToolResult, any, error) {
	prov, err := getProvider()
	if err != nil {
		return nil, nil, err
	}
	p := utils.ResolvePath(args.Path)
	n, err := prov.NumTestsInDirRec(p)
	if err != nil {
		return nil, nil, err
	}
	return utils.RespondWith(map[string]any{"num_tests": n, "path": args.Path, "recursive": true}), nil, nil
}

func GetTestsInDir(ctx context.Context, req *mcp.CallToolRequest, args GetTestsInDirParams) (*mcp.CallToolResult, any, error) {
	prov, err := getProvider()
	if err != nil {
		return nil, nil, err
	}
	p := utils.ResolvePath(args.Path)
	tests, err := prov.GetTestsInDir(p)
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

func GetTestsInDirRec(ctx context.Context, req *mcp.CallToolRequest, args GetTestsInDirParams) (*mcp.CallToolResult, any, error) {
	prov, err := getProvider()
	if err != nil {
		return nil, nil, err
	}
	p := utils.ResolvePath(args.Path)
	tests, err := prov.GetTestsInDirRec(p)
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

	if !strings.HasSuffix(p, ".js") {
		p += ".js"
	}

	out, status, err := prov.GetTestOutput(p)
	if err != nil {
		return nil, nil, err
	}
	return utils.RespondWith(map[string]any{"test_path": args.TestPath, "output": out, "status": status}), nil, nil
}

func SearchDir(ctx context.Context, req *mcp.CallToolRequest, args SearchDirParams) (*mcp.CallToolResult, any, error) {
	prov, err := getProvider()
	if err != nil {
		return nil, nil, err
	}
	results, err := prov.SearchDir(args.Query)
	if err != nil {
		return nil, nil, err
	}
	sort.Strings(results)
	page, pageSize := normalizePage(args.Page, args.PageSize)
	items, remaining, total := paginateStrings(results, page, pageSize, args.Max)
	res := map[string]any{
		"query":     args.Query,
		"page":      page,
		"page_size": pageSize,
		"returned":  len(items),
		"remaining": remaining,
		"total":     total,
		"results":   items,
	}
	return utils.RespondWith(res), nil, nil
}

func SearchDirIn(ctx context.Context, req *mcp.CallToolRequest, args SearchDirInParams) (*mcp.CallToolResult, any, error) {
	prov, err := getProvider()
	if err != nil {
		return nil, nil, err
	}
	dir := utils.ResolvePath(args.Dir)
	results, err := prov.SearchDirIn(dir, args.Query)
	if err != nil {
		return nil, nil, err
	}
	sort.Strings(results)
	page, pageSize := normalizePage(args.Page, args.PageSize)
	items, remaining, total := paginateStrings(results, page, pageSize, args.Max)
	res := map[string]any{
		"dir":       args.Dir,
		"query":     args.Query,
		"page":      page,
		"page_size": pageSize,
		"returned":  len(items),
		"remaining": remaining,
		"total":     total,
		"results":   items,
	}
	return utils.RespondWith(res), nil, nil
}

func SearchTest(ctx context.Context, req *mcp.CallToolRequest, args SearchTestParams) (*mcp.CallToolResult, any, error) {
	prov, err := getProvider()
	if err != nil {
		return nil, nil, err
	}
	results, err := prov.SearchTest(args.Query)
	if err != nil {
		return nil, nil, err
	}
	sort.Strings(results)
	page, pageSize := normalizePage(args.Page, args.PageSize)
	items, remaining, total := paginateStrings(results, page, pageSize, args.Max)
	res := map[string]any{
		"query":     args.Query,
		"page":      page,
		"page_size": pageSize,
		"returned":  len(items),
		"remaining": remaining,
		"total":     total,
		"tests":     items,
	}
	return utils.RespondWith(res), nil, nil
}

func SearchTestInDir(ctx context.Context, req *mcp.CallToolRequest, args SearchTestInDirParams) (*mcp.CallToolResult, any, error) {
	prov, err := getProvider()
	if err != nil {
		return nil, nil, err
	}
	dir := utils.ResolvePath(args.Dir)
	results, err := prov.SearchTestInDir(dir, args.Query)
	if err != nil {
		return nil, nil, err
	}
	sort.Strings(results)
	page, pageSize := normalizePage(args.Page, args.PageSize)
	items, remaining, total := paginateStrings(results, page, pageSize, args.Max)
	res := map[string]any{
		"dir":       args.Dir,
		"query":     args.Query,
		"page":      page,
		"page_size": pageSize,
		"returned":  len(items),
		"remaining": remaining,
		"total":     total,
		"tests":     items,
	}
	return utils.RespondWith(res), nil, nil
}

func AddTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "NumTestsTotal",
		Description: "Get the total number of tests (results from last CI run)",
	}, NumTestsTotal)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "NumTestsInDir",
		Description: "Get the number of tests in a directory (results from last CI run)",
	}, NumTestsInDir)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "NumTestsInDirRecursive",
		Description: "Get the number of tests in a directory recursively (results from last CI run)",
	}, NumTestsInDirRecursive)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "GetTestsInDirRecursive",
		Description: "List tests in a directory recursively (paginated) (results from last CI run)",
	}, GetTestsInDirRec)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "GetTestsInDir",
		Description: "List tests in a directory (paginated) (results from last CI run)",
	}, GetTestsInDir)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "GetTestStatus",
		Description: "Get the status of a single test (results from last CI run)",
	}, GetTestStatus)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "GetTestStatusesInDirRecursive",
		Description: "List statuses for tests in a directory recursively (paginated) (results from last CI run)",
	}, GetTestStatusesInDirRec)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "GetTestStatusesInDir",
		Description: "List statuses for tests in a directory (paginated) (results from last CI run)",
	}, GetTestStatusesInDir)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "GetTestsWithStatusInDirRecursive",
		Description: "List tests with a specific status in a directory recursively (paginated) (results from last CI run)",
	}, GetTestsWithStatusInDirRec)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "GetTestsWithStatusInDir",
		Description: "List tests with a specific status in a directory (paginated) (results from last CI run)",
	}, GetTestsWithStatusInDir)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "GetFailedTestsInDirRecursive",
		Description: "List failed tests in a directory recursively (paginated) (results from last CI run)",
	}, GetFailedTestsInDirRec)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "GetFailedTestsInDir",
		Description: "List failed tests in a directory (paginated) (results from last CI run)",
	}, GetFailedTestsInDir)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "GetTestsWithStatusInDirRecursive",
		Description: "List tests with a specific status in a directory recursively (paginated) (results from last CI run)",
	}, GetTestsWithStatusInDirRec)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "GetTestOutput",
		Description: "Get the output of a single test (results from last CI run)",
	}, GetTestOutput)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "SearchDir",
		Description: "Search repository paths by query (paginated) (results from last CI run)",
	}, SearchDir)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "SearchDirIn",
		Description: "Search repository paths within a directory by query (paginated) (results from last CI run)",
	}, SearchDirIn)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "SearchTest",
		Description: "Search tests by query (paginated) (results from last CI run)",
	}, SearchTest)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "SearchTestInDir",
		Description: "Search tests within a directory by query (paginated) (results from last CI run)",
	}, SearchTestInDir)
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

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

type GetSpecParams struct {
	SpecPath string `json:"spec_path" jsonschema:"Path to the specification file or section (e.g. array.prototype.tostring)"`
}

type SpecForIntrinsicParams struct {
	Intrinsic string `json:"intrinsic" jsonschema:"Intrinsic name to lookup in the spec (e.g. %Array.prototype.toString%)"`
}

type SearchSpecParams struct {
	Query    string `json:"query" jsonschema:"Search query"`
	Page     int    `json:"page" jsonschema:"Page number starting from 1; defaults to 1"`
	PageSize int    `json:"page_size" jsonschema:"Items per page; defaults to DefaultPageSize if omitted"`
	Max      int    `json:"max" jsonschema:"Optional global maximum number of items to include across all pages; 0 means no limit"`
}

type SearchSectionsParams struct {
	Query    string `json:"query" jsonschema:"Search query for spec sections"`
	Page     int    `json:"page" jsonschema:"Page number starting from 1; defaults to 1"`
	PageSize int    `json:"page_size" jsonschema:"Items per page; defaults to DefaultPageSize if omitted"`
	Max      int    `json:"max" jsonschema:"Optional global maximum number of items to include across all pages; 0 means no limit"`
}

func GetSpec(ctx context.Context, req *mcp.CallToolRequest, args GetSpecParams) (*mcp.CallToolResult, any, error) {
	s, err := getSpecProvider()
	if err != nil {
		return nil, nil, err
	}
	specPath := utils.ResolvePath(args.SpecPath)

	specPath = strings.TrimPrefix(specPath, "sec-")

	content, err := s.GetSpec(specPath)
	if err != nil {
		return nil, nil, err
	}
	return utils.RespondWith(map[string]any{"spec_path": args.SpecPath, "content": content}), nil, nil
}

func SpecForIntrinsic(ctx context.Context, req *mcp.CallToolRequest, args SpecForIntrinsicParams) (*mcp.CallToolResult, any, error) {
	s, err := getSpecProvider()
	if err != nil {
		return nil, nil, err
	}
	section, err := s.SpecForIntrinsic(args.Intrinsic)
	if err != nil {
		return nil, nil, err
	}
	return utils.RespondWith(map[string]any{"intrinsic": args.Intrinsic, "section": section}), nil, nil
}

func SearchSpec(ctx context.Context, req *mcp.CallToolRequest, args SearchSpecParams) (*mcp.CallToolResult, any, error) {
	s, err := getSpecProvider()
	if err != nil {
		return nil, nil, err
	}
	results, err := s.SearchSpec(args.Query)
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

func SearchSections(ctx context.Context, req *mcp.CallToolRequest, args SearchSectionsParams) (*mcp.CallToolResult, any, error) {
	s, err := getSpecProvider()
	if err != nil {
		return nil, nil, err
	}
	results, err := s.SearchSections(args.Query)
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
		"sections":  items,
	}
	return utils.RespondWith(res), nil, nil
}

func AddSpecTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "GetSpec",
		Description: "Get the content of a specification file or section",
	}, GetSpec)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "SpecForIntrinsic",
		Description: "Lookup the spec section for a given intrinsic",
	}, SpecForIntrinsic)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "SearchSpec",
		Description: "Search the specification content for a query (paginated)",
	}, SearchSpec)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "SearchSections",
		Description: "Search spec section titles/ids for a query (paginated)",
	}, SearchSections)
}

func getSpecProvider() (provider.SpecProvider, error) {
	if provider.Spec == nil {
		return nil, errors.New("spec provider not set")
	}
	return provider.Spec, nil
}

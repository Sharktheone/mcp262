package utils

import (
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func RespondWith(res any) *mcp.CallToolResult {
	b, _ := json.Marshal(res)

	return &mcp.CallToolResult{
		StructuredContent: json.RawMessage(b),
		Content:           []mcp.Content{&mcp.TextContent{Text: string(b)}},
	}
}

func ResolvePath(path string) string {
	if len(path) >= 6 && path[:6] == "/test/" {
		return path[5:]
	} else if len(path) >= 10 && path[:10] == "/test262/" {
		return path[9:]
	} else if len(path) >= 1 && path[0] == '/' {
		return path[1:]
	} else {
		return path
	}
}

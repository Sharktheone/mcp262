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

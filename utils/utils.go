package utils

import (
	"encoding/json"
	"log"
	"strings"

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
	log.Println(path)
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimPrefix(path, "test262/")

	if path == "test" {
		return ""

	}
	path = strings.TrimPrefix(path, "test/")

	return path
}

func SplitPath(p string) ([]string, string) {
	elems := strings.Split(p, "/")

	last := elems[len(elems)-1]

	if strings.Contains(last, ".") {
		return elems[:len(elems)-1], last
	}

	return elems, ""
}

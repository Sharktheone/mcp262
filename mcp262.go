package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Sharktheone/mcp262/runner"

	"github.com/Sharktheone/mcp262/provider"
	"github.com/Sharktheone/mcp262/provider/github"
	"github.com/Sharktheone/mcp262/provider/yavashark"
	"github.com/Sharktheone/mcp262/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	p, err := yavashark.NewYavasharkTestProvider()
	if err != nil {
		log.Fatalf("Failed to create YavasharkTestProvider: %v", err)
		return
	}

	provider.SetProvider(p)
	provider.SetCodeProvider(github.NewGithubTest262CodeProvider())
	provider.SetSpecProvider(github.NewGithubSpecProvider())

	//url := "0.0.0.0:8080"

	server := mcp.NewServer(&mcp.Implementation{Name: "mcp262", Version: "v1.0.0", Title: "mcp262"}, nil)

	tools.AddTools(server)
	tools.AddCodeTools(server)
	tools.AddSpecTools(server)
	tools.AddRunnerTools(server)

	//handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
	//	return server
	//}, nil)

	//handler := mcp.NewSSEHandler(func(req *http.Request) *mcp.Server {
	//	return server
	//})
	//
	//handlerWithLogging := loggingHandler(handler)
	//
	//log.Printf("MCP server listening on %s", url)
	//
	//if err := http.ListenAndServe(url, handlerWithLogging); err != nil {
	//	log.Fatalf("Server failed: %v", err)
	//}

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func loggingHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		log.Printf("[REQUEST] %s | %s | %s %s",
			start.Format(time.RFC3339),
			r.RemoteAddr,
			r.Method,
			r.URL.Path)

		handler.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		log.Printf("[RESPONSE] %s | %s | %s %s | Status: %d | Duration: %v",
			time.Now().Format(time.RFC3339),
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			wrapped.statusCode,
			duration)
	})
}

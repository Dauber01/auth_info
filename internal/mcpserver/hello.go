package mcpserver

import (
	"context"
	"net/http"
	"strings"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

	bizhello "auth_info/internal/biz/hello"
)

type helloInput struct {
	Name string `json:"name,omitempty" jsonschema:"Name to greet. Defaults to World when empty."`
}

type helloOutput struct {
	Message string `json:"message" jsonschema:"Greeting message returned by the hello service."`
}

func NewHelloMCPServer(uc *bizhello.UseCase) *sdkmcp.Server {
	server := sdkmcp.NewServer(&sdkmcp.Implementation{
		Name:    "auth-info",
		Title:   "Auth Info",
		Version: "1.0.0",
	}, nil)

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "hello",
		Title:       "Hello",
		Description: "Call the same capability exposed by GET /api/v1/hello and return a greeting message.",
	}, func(ctx context.Context, _ *sdkmcp.CallToolRequest, input helloInput) (*sdkmcp.CallToolResult, helloOutput, error) {
		message := uc.SayHello(ctx, strings.TrimSpace(input.Name))
		return nil, helloOutput{Message: message}, nil
	})

	return server
}

func NewHelloMCPHandler(uc *bizhello.UseCase) http.Handler {
	server := NewHelloMCPServer(uc)
	return sdkmcp.NewStreamableHTTPHandler(func(*http.Request) *sdkmcp.Server {
		return server
	}, nil)
}

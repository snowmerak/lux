package mcp

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// GoDocInput is the input schema for go doc lookup.
type GoDocInput struct {
	Target string `json:"target" jsonschema:"the package, function, or type to look up. e.g. 'net/http', 'net/http.Client', 'fmt.Printf'"`
}

// GoDocOutput is the output schema for go doc lookup.
type GoDocOutput struct {
	Content string `json:"content" jsonschema:"the documentation result from go doc"`
}

// GetGoDocTool provides documentation results by running 'go doc'.
func GetGoDocTool(ctx context.Context, req *mcp.CallToolRequest, input GoDocInput) (
	*mcp.CallToolResult,
	GoDocOutput,
	error,
) {
	// target이 비어있으면 에러 반환
	if strings.TrimSpace(input.Target) == "" {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Target is required.",
				},
			},
		}, GoDocOutput{}, nil
	}

	// go doc -all <target> 실행
	cmd := exec.CommandContext(ctx, "go", "doc", "-all", input.Target)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("go doc failed: %s\nOutput: %s", err, string(out)),
				},
			},
		}, GoDocOutput{}, nil
	}

	return nil, GoDocOutput{
		Content: string(out),
	}, nil
}

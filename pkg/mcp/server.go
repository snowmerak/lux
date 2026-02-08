package mcp

import (
	"context"
	"fmt"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/snowmerak/lux/v4/pkg/knowledge"
)

// SearchInput은 지식 검색을 위한 입력 스키마입니다.
type SearchInput struct {
	Tags []string `json:"tags" jsonschema:"the list of tags to search for"`
}

// SearchOutput은 지식 검색 결과(타이틀 목록)를 담는 출력 스키마입니다.
type SearchOutput struct {
	Results []string `json:"results" jsonschema:"the list of knowledge titles and IDs prioritized by tag matching"`
}

// ContentInput은 특정 지식 조회를 위한 입력 스키마입니다.
type ContentInput struct {
	ID string `json:"id" jsonschema:"the unique ID of the knowledge to retrieve"`
}

// ContentOutput은 지식 상세 내용을 담는 출력 스키마입니다.
type ContentOutput struct {
	Title   string `json:"title" jsonschema:"the title of the knowledge"`
	Content string `json:"content" jsonschema:"the full content of the knowledge"`
}

// SearchKnowledgeTool은 태그 매칭 기반으로 상위 10개의 지식 정보를 제공합니다.
func SearchKnowledgeTool(ctx context.Context, req *mcp.CallToolRequest, input SearchInput) (
	*mcp.CallToolResult,
	SearchOutput,
	error,
) {
	results, err := knowledge.SearchKnowledge(ctx, input.Tags)
	log.Printf("Search for tags %v found %d results", input.Tags, len(results))
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Search failed: %s", err),
				},
			},
		}, SearchOutput{}, nil
	}

	output := SearchOutput{
		Results: []string{},
	}

	for _, res := range results {
		output.Results = append(output.Results, fmt.Sprintf("[%s] %s (Score: %.2f)", res.Item.ID, res.Item.Title, res.Score))
	}

	return nil, output, nil
}

// GetKnowledgeContentTool은 ID를 기반으로 지식의 전체 데이터를 제공합니다.
func GetKnowledgeContentTool(ctx context.Context, req *mcp.CallToolRequest, input ContentInput) (
	*mcp.CallToolResult,
	ContentOutput,
	error,
) {
	item, ok := knowledge.GetContentByID(ctx, input.ID)
	if !ok {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Knowledge not found for the given ID.",
				},
			},
		}, ContentOutput{}, nil
	}

	return nil, ContentOutput{
		Title:   item.Title,
		Content: item.Content,
	}, nil
}

// NewServer는 사용자 예제 패턴을 따라 MCP 서버를 초기화합니다.
func NewServer() *mcp.Server {
	s := mcp.NewServer(
		&mcp.Implementation{
			Name:    "lux-knowledge",
			Version: "v1.0.0",
		},
		nil,
	)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "search_knowledge",
		Description: "Search for the top 10 knowledge titles based on tag matching. More tag matches mean higher priority.",
	}, SearchKnowledgeTool)

	mcp.AddTool(s, &mcp.Tool{
		Name:        "get_knowledge_content",
		Description: "Get the full content and tags of a knowledge item by its ID.",
	}, GetKnowledgeContentTool)

	return s
}

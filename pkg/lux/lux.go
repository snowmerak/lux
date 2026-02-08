package lux

import (
	"context"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	luxmcp "github.com/snowmerak/lux/v4/pkg/mcp"
)

// App은 Lux 프레임워크의 메인 구조체입니다.
type App struct {
	mcpServer *mcp.Server
	server    *http.Server
}

// New는 새로운 Lux 애플리케이션을 생성합니다.
func New() *App {
	return &App{
		mcpServer: luxmcp.NewServer(),
		server:    &http.Server{},
	}
}

// Run은 서버를 시작합니다.
func (a *App) Run(ctx context.Context, addr string) error {
	// Stdio를 통해 MCP 서버 실행 (사용자 예제 패턴)
	return a.mcpServer.Run(ctx, &mcp.StdioTransport{})
}

package datasource

import "context"

// Client는 데이터베이스 또는 외부 서비스 클라이언트의 공통 인터페이스입니다.
type Client interface {
	Name() string
	Close() error
	Ping(ctx context.Context) error
}

// Manager는 활성화된 데이터소스 클라이언트들을 관리합니다.
type Manager struct {
	clients map[string]Client
}

func NewManager() *Manager {
	return &Manager{
		clients: make(map[string]Client),
	}
}

func (m *Manager) Add(c Client) {
	m.clients[c.Name()] = c
}

func (m *Manager) Get(name string) (Client, bool) {
	c, ok := m.clients[name]
	return c, ok
}

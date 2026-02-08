# Taskfile Automation Guide

This guide defines the standard for automating project workflows (build, test, indexing, deployment) using [Task](https://taskfile.dev/). A well-defined `Taskfile.yml` ensures consistent development environments and CI/CD pipelines.

### Core Principles

1.  **Tool Dependency Management (Go 1.24+)**: It is mandatory to manage `task` as a versioned tool dependency.
    - **Installation**: Run `go get -tool github.com/go-task/task/v3/cmd/task@latest`.
    - **Usage**: Use `go tool task` to execute the local version defined in `go.mod`.
2.  **Atomicity**: Each task should perform a single, logical unit of work.
3.  **Dependency Management**: Use `deps` or `cmds` with `task:` to chain related operations (e.g., `build` should depend on `index`).
3.  **Descriptive**: Every task **must** have a `desc` field for discoverability (`task --list`).
4.  **Environment Agnostic**: Use environment variables for paths and configurations where possible.

### Standard Task Templates

#### 1. Core Workflow (Build & Install)

```yaml
version: '3'

tasks:
  install:
    desc: Install dependencies and prepare the environment
    cmds:
      - go mod tidy
      - go mod download

  build:
    desc: Build the application binary
    cmds:
      - task: index
      - go build -o {{.APP_NAME}}{{.EXE}} ./cmd/{{.APP_NAME}}
    vars:
      APP_NAME: lux
      EXE: '{{if eq OS "windows"}}.exe{{else}}{{end}}'
```

#### 2. Development & Testing

```yaml
tasks:
  test:
    desc: Run unit tests
    cmds:
      - go test -v -race ./pkg/... ./internal/...

  lint:
    desc: Run golangci-lint
    cmds:
      - golangci-lint run ./...
```

#### 3. Project Specific (Indexing)

```yaml
tasks:
  index:
    desc: Generate search index from markdown files
    cmds:
      - go run cmd/indexer/main.go
    sources:
      - knowledge/**/*.md
    generates:
      - pkg/knowledge/data.bleve/**/*
```

### Integration with Project Structure

- **Location**: The `Taskfile.yml` must reside in the project root.
- **Sub-tasks**: For complex sub-systems, you can use `includes` to split the Taskfile, but keep the main entry points in the root file.

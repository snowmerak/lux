# Lux

Lux is a tool for managing a knowledge repository (.md files) and providing it to AI agents via the [Model Context Protocol (MCP)](https://modelcontextprotocol.io/). This project features the ability to index and search highly structured knowledge data.

## Project Structure

```text
.
├── cmd/
│   ├── indexer/          # Tool to read Markdown files and generate a search index (Bleve)
│   └── lux/              # Entry point for running the MCP server
├── knowledge/            # Markdown-based knowledge repository (Raw data)
├── pkg/
│   ├── knowledge/        # Knowledge data modeling and Bleve search engine logic
│   ├── lux/              # Lux application framework core
│   └── mcp/              # MCP server implementation and Tool definitions
├── Taskfile.yml          # Build and task automation configuration
└── go.mod                # Go module dependency definition
```

## Key Features

- **Knowledge Indexing**: Analyzes Markdown files in the `knowledge/` directory to create a `pkg/knowledge/data.bleve` index.
- **MCP Server**: Provides a standard MCP interface so AI agents (e.g., Claude Desktop) can search and retrieve knowledge.
- **Knowledge Search (Tools)**:
  - `search_knowledge`: Searches for a list of knowledge items based on tag matching.
  - `get_knowledge_content`: Retrieves the detailed content of a specific knowledge item by its ID.
- **Embedded Index**: The generated search index is embedded into the binary for easy distribution and deployment.

## Getting Started

### Prerequisites

- Go 1.25 or higher
- [Go Task](https://taskfile.dev/) (Optional, if using `Taskfile.yml`)

### Installation & Build

1. Install dependencies and generate the index:
   ```bash
   task index
   ```

2. Build the binary:
   ```bash
   task build
   ```
   or install binary via go install to `$GOPATH/bin/lux`:
   ```bash
   task install
   ```

### Running

Run the MCP server using the Stdio transport:
```bash
./lux.exe
```

## Knowledge Guides

This project includes various guides for development. Refer to these when adding new modules or extending functionality:

- [Project Structure](knowledge/project_structure.md): Project layering and architectural principles.
- [Entry Point Guide](knowledge/entry_point.md): Guide for the `cmd` layer and using Fx.
- [Module Generation Guide](knowledge/module.md): Guide for creating Fx modules and managing environment variables.

## License

This project is licensed under the [MIT License](LICENSE).

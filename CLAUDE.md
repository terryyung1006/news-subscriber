# Antigravity AI Agent Rules

This project uses modular Claude Code rules located in `.claude/rules/`:

- `general.md` - Code reuse, clean hierarchy, documentation limits
- `grpc-workflow.md` - gRPC/proto file editing and code generation
- `project-structure.md` - Root and service directory organization
- `feature-workflow.md` - Feature branch and documentation workflow

You must follow all rules defined in the `.claude/rules/` directory when working on this codebase.

## Custom Skills

Available slash commands (`.claude/skills/`):

- `/proto-sync [validate]` - Regenerate proto code across backend and frontend
- `/feature <name>` - Create feature branch with documentation structure
- `/services <up|down|status|logs>` - Manage Docker Compose services
- `/db <migrate|reset|seed|fresh>` - Database operations
- `/db-query <postgres|chroma|redis> [query]` - Quick database query
- `/test <service> [flags]` - Run tests for a specific service
- `/test-gen <service> [file-path]` - Generate test file for a source file

## Custom Agents

Available subagents (`.claude/agents/`):

- `grpc-validator` - Validates gRPC stack after proto changes
- `multi-service-test` - Runs tests across all services
- `kafka-debug` - Debug Kafka event pipeline issues
- `db-inspector` - Multi-database inspection across PostgreSQL, ChromaDB, and Redis
- `test-scaffold-go` - Go test generator for backend and news-fetcher
- `test-scaffold-python` - Python test generator for ETL and inference workers
- `data-pipeline-tracer` - End-to-end data flow tracer across the pipeline

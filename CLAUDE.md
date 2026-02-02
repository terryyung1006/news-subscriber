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

## Custom Agents

Available subagents (`.claude/agents/`):

- `grpc-validator` - Validates gRPC stack after proto changes
- `multi-service-test` - Runs tests across all services
- `kafka-debug` - Debug Kafka event pipeline issues

# Project Structure

## Clean Root Directory

The root folder should ONLY contain:
- `CLAUDE.md` - Minimal project context file
- `.claude/` - Claude Code rules directory
- `README.md` - Project overview with links
- `ARCHITECTURE.md` - System architecture documentation
- Individual service directories (e.g., `news-subscriber-core`, `news-subscriber-frontend`)

All other files (docker-compose, Makefiles, setup guides, etc.) belong in their respective service directories.

## Service Directory Organization

Each service should be self-contained with its own:
- Configuration files
- Build scripts
- Docker configuration
- Service-specific documentation

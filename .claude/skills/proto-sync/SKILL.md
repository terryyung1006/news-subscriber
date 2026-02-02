---
name: proto-sync
description: Regenerate proto code across backend and frontend services after proto file changes. Use when proto files have been modified or when gRPC code needs updating.
allowed-tools: Bash(make *), Bash(npm run *), Bash(cd *), Bash(go build *)
---

# Proto Sync Workflow

Regenerate gRPC code from proto definitions in both backend and frontend.

## Steps

1. **Regenerate Backend Proto Code**
   ```bash
   cd /Users/terryyung/SourceCode/news-subscriber/backend && make gen-proto
   ```

2. **Regenerate Frontend Proto Code**
   ```bash
   cd /Users/terryyung/SourceCode/news-subscriber/frontend && npm run gen-proto
   ```

3. **Validate Backend Compilation** (if $ARGUMENTS contains "validate")
   ```bash
   cd /Users/terryyung/SourceCode/news-subscriber/backend && go build ./...
   ```

4. **Validate Frontend Types** (if $ARGUMENTS contains "validate")
   ```bash
   cd /Users/terryyung/SourceCode/news-subscriber/frontend && npx tsc --noEmit
   ```

## Proto Locations
- Source protos: `protos/proto/*.proto`
- Backend generated: `backend/api/proto/v1/`
- Frontend generated: `frontend/src/lib/proto/`

## Usage
- `/proto-sync` - Regenerate code only
- `/proto-sync validate` - Regenerate and validate compilation

Report any errors encountered during generation or validation.

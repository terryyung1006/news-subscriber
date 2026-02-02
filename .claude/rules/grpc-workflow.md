# gRPC & Proto Workflow

All API communication between services must be implemented using gRPC.

## Proto File Editing Rules

- **NEVER edit proto files directly** in `news-subscriber-frontend/proto-repo/` or `news-subscriber-core/proto-repo/` - these are git submodules
- **ALWAYS edit proto files** in the `news-subscriber-grpc-proto` repository first

## Submodule Update Commands

After updating `news-subscriber-grpc-proto`, update submodules in both frontend and core:

```bash
# Frontend
cd news-subscriber-frontend && git submodule update --remote proto-repo

# Core
cd news-subscriber-core && git submodule update --remote proto-repo
```

## Code Generation Commands

After updating submodules, regenerate proto code in each service:

```bash
# Frontend
npm run gen-proto

# Core
make gen-proto
```

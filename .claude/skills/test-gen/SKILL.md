---
name: test-gen
description: Generate a test file for a specific source file. Use for incremental test writing during development.
argument-hint: <service> [file-path]
allowed-tools: Bash(go *), Bash(uv run *), Bash(python *), Read, Glob, Grep, Write
---

# Generate Test File

Generate a test file for a source file based on `$ARGUMENTS`.

The first word of `$ARGUMENTS` specifies the service. The optional second argument is the relative path to the source file. If no file path is given, prompt the user to specify one.

## Process

1. **Read the target source file** to understand functions, classes, and dependencies
2. **Check for existing test conventions** in the service's test directory
3. **Generate a test file** with 3-5 test cases per function (happy path + error paths)
4. **Validate** that generated tests compile/import correctly

## Go Services (`backend`, `news-fetcher`)

- Generate `_test.go` in the same directory as the source file
- Use table-driven tests with `t.Run()` subtests
- Create mock structs for interfaces inline
- Include happy path, error path, and edge case subtests

Validation:
```bash
cd <service-dir> && go build ./...
```

### Example
For `backend src/service/auth.go`, generate `backend/src/service/auth_test.go`.

## Python Services (`etl-worker`, `inference-worker`)

- Generate `tests/test_<module>.py` in the service's test directory
- Use pytest with `unittest.mock` for mocking
- Apply `@pytest.mark.unit` markers
- Use fixtures for shared setup
- If `conftest.py` doesn't exist, create it with common fixtures

Validation:
```bash
# ETL worker
cd <service-dir> && uv run python -m pytest tests/<test_file> --collect-only

# Inference worker
cd <service-dir> && python -m pytest tests/<test_file> --collect-only
```

### Example
For `etl-worker src/storage/chroma_storage.py`, generate `etl-worker/tests/test_chroma_storage.py`.

## Service Directories

| Service | Root Directory |
|---------|---------------|
| backend | `news-subscriber-core/backend/` |
| news-fetcher | `news-subscriber-core/news-fetcher/` |
| etl-worker | `news-subscriber-core/etl-worker/` |
| inference-worker | `news-subscriber-core/inference-worker/` |
| frontend | `news-subscriber-frontend/` |

## Usage
- `/test-gen backend src/service/auth.go` — Generate tests for auth service
- `/test-gen etl-worker src/storage/chroma_storage.py` — Generate tests for ChromaDB storage
- `/test-gen news-fetcher internal/fetcher/rss.go` — Generate tests for RSS fetcher
- `/test-gen inference-worker src/worker.py` — Generate tests for inference worker

Report: generated file path, test count, and compilation/import status.

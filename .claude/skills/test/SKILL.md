---
name: test
description: Run tests for a specific service with optional flags. Use for targeted testing of individual services.
argument-hint: <backend|frontend|etl-worker|news-fetcher|inference-worker> [flags]
allowed-tools: Bash(go test *), Bash(make *), Bash(uv run *), Bash(python -m pytest *), Bash(npm run *), Bash(cd *)
---

# Run Service Tests

Run tests for the specified service based on `$ARGUMENTS`.

The first word of `$ARGUMENTS` specifies the service. Remaining arguments are passed as flags to the test command.

## Service Commands

### `backend [flags]` — Go Tests
```bash
cd /Users/terryyung/SourceCode/news-subscriber-branches/subagents-and-skills/news-subscriber/news-subscriber-core/backend && go test ./... $FLAGS
```
Common flags: `-v`, `-run TestXxx`, `-count=1`, `-cover`, `-race`

### `frontend [flags]` — Lint/Type Check
```bash
cd /Users/terryyung/SourceCode/news-subscriber-branches/subagents-and-skills/news-subscriber/news-subscriber-frontend && npm run lint $FLAGS
```

### `etl-worker [flags]` — Pytest
```bash
cd /Users/terryyung/SourceCode/news-subscriber-branches/subagents-and-skills/news-subscriber/news-subscriber-core/etl-worker && uv run pytest tests/ $FLAGS
```
Common flags: `-v`, `-k expression`, `-m unit`, `--cov`, `-x` (stop on first failure)

### `news-fetcher [flags]` — Go Tests
```bash
cd /Users/terryyung/SourceCode/news-subscriber-branches/subagents-and-skills/news-subscriber/news-subscriber-core/news-fetcher && go test -v ./... $FLAGS
```
Common flags: `-run TestXxx`, `-count=1`, `-cover`, `-race`

### `inference-worker [flags]` — Pytest
```bash
cd /Users/terryyung/SourceCode/news-subscriber-branches/subagents-and-skills/news-subscriber/news-subscriber-core/inference-worker && python -m pytest tests/ $FLAGS
```
Common flags: `-v`, `-k expression`, `-m unit`, `--cov`, `-x`

## Usage
- `/test backend` — Run all backend tests
- `/test backend -run TestAuth -v` — Run specific Go test verbosely
- `/test etl-worker -m unit` — Run only unit-marked ETL tests
- `/test etl-worker --cov` — Run ETL tests with coverage
- `/test frontend` — Run frontend linting
- `/test news-fetcher -race` — Run news-fetcher tests with race detection
- `/test inference-worker -k test_worker` — Run specific inference tests

Report test results with pass/fail count and any error output.

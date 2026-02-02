# Feature Development Workflow

## Feature Branch Workflow

1. ALWAYS checkout `master` and pull latest before starting a new feature
2. Create a new feature branch: `git checkout -b feature/[feature-name]`
3. Create the `spec/feature/[feature-name]/` documentation FIRST
4. Implement the feature
5. Verify implementation against `plan.md`

## Feature Documentation Structure

For new features, create a subfolder in `spec/feature/[branch-name]/` (e.g., `spec/feature/inference-worker/`).

Each feature spec folder contains ONLY TWO files:
- `reference.md` - Quick reference (commands, key files, config)
- `plan.md` - Detailed implementation plan and progress

This structure serves as both planning documentation and implementation record. Keep documentation minimal and focused.

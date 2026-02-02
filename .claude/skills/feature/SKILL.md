---
name: feature
description: Create a new feature branch with documentation structure. Use when starting work on a new feature.
argument-hint: <feature-name>
disable-model-invocation: true
allowed-tools: Bash(git *), Write
---

# Feature Branch Setup

Create a new feature branch and documentation structure for `$ARGUMENTS`.

## Steps

1. **Update from master**
   ```bash
   cd /Users/terryyung/SourceCode/news-subscriber && git checkout master && git pull origin master
   ```

2. **Create feature branch**
   ```bash
   git checkout -b feature/$ARGUMENTS
   ```

3. **Create feature documentation directory**
   Create `spec/feature/$ARGUMENTS/` with two files:

   **reference.md:**
   ```markdown
   # $ARGUMENTS - Quick Reference

   ## Key Files
   - TBD

   ## Commands
   - TBD

   ## Configuration
   - TBD
   ```

   **plan.md:**
   ```markdown
   # $ARGUMENTS Implementation Plan

   ## Overview
   Brief description of the feature.

   ## Goals
   - [ ] Goal 1
   - [ ] Goal 2

   ## Implementation Steps
   1. Step 1
   2. Step 2

   ## Progress
   - [ ] Planning complete
   - [ ] Implementation started
   - [ ] Tests added
   - [ ] Documentation updated
   ```

4. **Initial commit**
   ```bash
   git add spec/feature/$ARGUMENTS/
   git commit -m "docs: Add feature spec for $ARGUMENTS"
   ```

## Usage
`/feature user-auth` - Creates feature/user-auth branch with spec docs

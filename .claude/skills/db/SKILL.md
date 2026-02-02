---
name: db
description: Database operations for news-subscriber backend including migrations, reset, and seeding. Use for database management tasks.
argument-hint: <migrate|reset|seed|fresh>
allowed-tools: Bash(make *), Bash(docker *)
---

# Database Operations

Execute database operation based on `$ARGUMENTS`.

All commands run from the backend directory:
```bash
cd /Users/terryyung/SourceCode/news-subscriber/backend
```

## Commands

### `migrate` - Run Pending Migrations
```bash
make db-migrate
```
Applies pending database migrations.

### `reset` - Reset Database Schema
```bash
make db-reset
```
Drops and recreates all tables.

### `seed` - Seed Database
```bash
make db-seed
```
Populates database with initial/test data.

### `fresh` - Complete Reset Cycle
```bash
make db-fresh
```
Full reset: stops containers, removes volumes, restarts, migrates, and seeds.

## Database Connection Info
- Host: `localhost`
- Port: `5433`
- Database: `news_subscriber`
- User: `user`
- Password: `password`

## Prerequisites
- Docker containers must be running (`/services up`)
- Go toolchain for migration tool

## Usage
- `/db migrate` - Apply pending migrations
- `/db reset` - Drop and recreate tables
- `/db seed` - Populate with test data
- `/db fresh` - Complete database reset

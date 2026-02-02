---
name: services
description: Manage Docker Compose services across all project services. Use to start, stop, check status, or view logs of infrastructure.
argument-hint: <up|down|status|logs [service]>
allowed-tools: Bash(docker *), Bash(docker-compose *)
---

# Service Management

Manage Docker Compose services based on `$ARGUMENTS`.

## Commands

### `up` - Start All Services
```bash
cd /Users/terryyung/SourceCode/news-subscriber/backend && docker-compose up -d
cd /Users/terryyung/SourceCode/news-subscriber/etl-worker && docker-compose up -d
```

### `down` - Stop All Services
```bash
cd /Users/terryyung/SourceCode/news-subscriber/backend && docker-compose down
cd /Users/terryyung/SourceCode/news-subscriber/etl-worker && docker-compose down
cd /Users/terryyung/SourceCode/news-subscriber/news-fetcher && docker-compose down
```

### `status` - Show Running Containers
```bash
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
```

### `logs [service]` - Tail Service Logs
For a specific container:
```bash
docker logs -f <container_name>
```

## Service Locations

| Directory | Services |
|-----------|----------|
| `backend/` | postgres, redis, ollama, chroma |
| `etl-worker/` | kafka, zookeeper |
| `news-fetcher/` | kafka, news-fetcher |

## Key Container Names
- `news_subscriber_postgres` - PostgreSQL (port 5433)
- `news_subscriber_redis` - Redis (port 6379)
- `news_subscriber_ollama` - Ollama LLM (port 11434)
- `news_subscriber_chroma` - ChromaDB (port 8000)

## Usage
- `/services up` - Start all infrastructure
- `/services down` - Stop all services
- `/services status` - Show what's running
- `/services logs news_subscriber_postgres` - Tail postgres logs

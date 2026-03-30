---
name: db-query
description: Quick database query against PostgreSQL, ChromaDB, or Redis. Use for ad-hoc data inspection without launching the full db-inspector agent.
argument-hint: <postgres|chroma|redis> [query]
allowed-tools: Bash(docker exec *), Bash(curl *)
---

# Database Query

Run a quick query against the specified database based on `$ARGUMENTS`.

The first word of `$ARGUMENTS` specifies the database target. Everything after is the query. If no query is provided, show default useful information.

## Targets

### `postgres [SQL]` — PostgreSQL Query
With query:
```bash
docker exec news_subscriber_postgres psql -U user -d news_subscriber -c "$QUERY"
```

Without query (default — list tables with row counts):
```bash
docker exec news_subscriber_postgres psql -U user -d news_subscriber -c "
SELECT schemaname, tablename,
  (xpath('/row/cnt/text()', xml_count))[1]::text::int as row_count
FROM (
  SELECT schemaname, tablename,
    query_to_xml('SELECT count(*) as cnt FROM ' || schemaname || '.' || tablename, false, true, '') as xml_count
  FROM pg_tables
  WHERE schemaname = 'public'
) t
ORDER BY tablename;
"
```

### `chroma [query]` — ChromaDB Query
With query (search text):
```bash
# First get collections
COLLECTIONS=$(curl -s http://localhost:8000/api/v2/collections)
# Then query the first collection with the search text
COLLECTION_ID=$(echo $COLLECTIONS | python3 -c "import sys,json; print(json.loads(sys.stdin.read())[0]['id'])")
curl -s -X POST "http://localhost:8000/api/v2/collections/$COLLECTION_ID/query" \
  -H "Content-Type: application/json" \
  -d "{\"query_texts\": [\"$QUERY\"], \"n_results\": 5}" | python3 -m json.tool
```

Without query (default — list collections with counts):
```bash
curl -s http://localhost:8000/api/v2/collections | python3 -m json.tool
```

### `redis [command]` — Redis Command
With command:
```bash
docker exec news_subscriber_redis redis-cli $QUERY
```

Without command (default — keyspace info and sample keys):
```bash
docker exec news_subscriber_redis redis-cli INFO keyspace
docker exec news_subscriber_redis redis-cli KEYS "*" | head -20
```

## Connection Info
| Database | Container | Port |
|----------|-----------|------|
| PostgreSQL | news_subscriber_postgres | 5433 |
| ChromaDB | news_subscriber_chroma | 8000 |
| Redis | news_subscriber_redis | 6379 |

## Prerequisites
Docker containers must be running (`/services up`).

## Usage
- `/db-query postgres` — List all tables with row counts
- `/db-query postgres SELECT * FROM users LIMIT 5;` — Run SQL query
- `/db-query chroma` — List ChromaDB collections
- `/db-query chroma artificial intelligence` — Search embeddings
- `/db-query redis` — Show keyspace info and sample keys
- `/db-query redis GET session:abc123` — Run Redis command

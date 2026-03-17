# Fetch News Feature Reference

## Quick Commands

```bash
# Start infrastructure
cd news-subscriber && docker compose up -d

# Run news-fetcher (Go)
cd news-fetcher && go run src/main.go

# Run etl-worker (Python/Spark)
cd etl-worker && python src/main.py

# Check Kafka topic
kafka-console-consumer --bootstrap-server localhost:9092 --topic news --from-beginning
```

## Key Files

| Service | File | Purpose |
|---------|------|---------|
| news-fetcher | `src/fetcher/newsapi.go` | NewsAPI client with quota tracking |
| news-fetcher | `src/main.go` | Entry point, Kafka producer |
| etl-worker | `src/main.py` | Spark streaming entry point |
| etl-worker | `src/spark/stream_processors/chunk_processor.py` | Sentence-based text chunking |

## Configuration

### NewsAPI (news-fetcher)
- **API Key**: Set via `NEWSAPI_API_KEY` env var
- **Query**: `business OR finance OR stocks`
- **Rate Limit**: 100 requests/day (free tier)
- **Kafka Topic**: `news`

### ETL Worker
- **Kafka Bootstrap**: `KAFKA_BOOTSTRAP_SERVERS`
- **Kafka Topic**: `news`
- **ChromaDB**: `CHROMA_HOST`, `CHROMA_PORT`
- **Chunk Size**: 500 chars max (sentence-based)

## Data Flow

```
NewsAPI -> news-fetcher -> Kafka (topic: news) -> etl-worker -> ChromaDB
```

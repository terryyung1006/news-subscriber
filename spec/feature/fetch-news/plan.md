# Fetch News Implementation Plan

## Objective

Make news-fetcher query NewsAPI and ensure articles flow through to ChromaDB via etl-worker.

## Issues Fixed

| # | Issue | Fix |
|---|-------|-----|
| 1 | Kafka topic mismatch (fetcher: `news`, worker: `text-data`) | Changed `job.start()` to `job.start("news")` |
| 2 | Missing required `q` parameter for /everything endpoint | Added `q=business OR finance OR stocks` |
| 3 | Article ID uses `len(url)` (not unique) | Changed to SHA256 hash of URL |
| 4 | Character-based chunking breaks mid-sentence | Implemented sentence-based chunking |
| 5 | No daily request quota tracking | Added counter with daily reset and warnings |

## Files Modified

### news-fetcher/src/fetcher/newsapi.go
- Added `q` parameter to NewsAPI query
- Fixed `generateArticleID()` to use SHA256 hash
- Added request quota tracking with daily reset
- Added warning when approaching limit (< 20 remaining)

### etl-worker/src/main.py
- Changed `job.start()` to `job.start("news")`
- Updated ChunkRowProcessor initialization

### etl-worker/src/spark/stream_processors/chunk_processor.py
- Replaced character-based chunking with sentence-based
- Chunks split on sentence boundaries (`.`, `!`, `?`)
- Max chunk size: 500 chars
- Min chunk size: 100 chars (merges small trailing chunks)

## Verification Steps

1. Start infrastructure: `docker compose up -d` (Kafka, ChromaDB)
2. Start etl-worker: `cd etl-worker && python src/main.py`
3. Run news-fetcher: `cd news-fetcher && go run src/main.go`
4. Check Kafka: `kafka-console-consumer --topic news`
5. Query ChromaDB to verify documents

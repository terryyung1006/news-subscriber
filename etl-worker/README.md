# RAG Financial News Worker

A scalable consumer service that uses Apache Spark Structured Streaming to process financial news from Kafka and store vectors in Chroma. Built with SOLID/OOP design and modern tooling (pyproject + uv).

## 🏗️ Architecture

```
Kafka (multiple topics) -> Spark Structured Streaming (multi-worker) -> Chroma (vector DB)
```

- This project is a **consumer** only. It does not produce data.
- Data is ingested from Kafka by Spark, processed (cleaning, sentiment, embeddings), and persisted in Chroma.
- **NEW**: Multi-worker architecture processes multiple Kafka topics simultaneously with different processing pipelines.

## 🚀 Quick Start

### With Docker

```bash
# Start infra and worker
./start.sh
# or
docker-compose up -d --build
```

### Local Dev (uv)

```bash
# Set up dev env with uv
./scripts/setup-dev.sh

# Run the consumer locally
uv run python src/main.py
```

## 🔧 Configuration

`config/config.yaml`
- `kafka.bootstrap_servers`: e.g. `localhost:9092`
- `kafka.topics`: multiple topics for different data types
- `chroma`: connection and collection settings
- `spark`: app name, master, and tuning options

## 📦 Tooling

- Dependency management: `uv` with `pyproject.toml`
- Testing: `pytest`
- Linting/formatting: `ruff`, `black`, `isort`

Key commands:
```bash
make setup-dev    # prepare dev env with uv
make test         # run tests
make format       # format code
make lint         # lint code
uv run python src/main.py  # run consumer
```

## 🧩 Spark Streaming Job

The `SparkStreamingJob` class now supports multiple workers:

- **Financial News Worker**: Full processing pipeline (news + sentiment + embedding)
- **User Queries Worker**: Raw storage for user query analytics

Each worker:
- Processes different Kafka topics independently
- Applies different processing pipelines
- Uses separate checkpoint locations for fault tolerance
- Writes to different Chroma collections

### Key Features

- **Kafka DataFrame Factory**: Abstracted Kafka DataFrame creation for multiple topics
- **Configurable Processors**: Apply different processors per worker (news, sentiment, embedding)
- **Independent Processing**: Each worker runs its own streaming query
- **Scalable Architecture**: Easy to add more workers and topics

## 🧪 Tests

Run unit tests:
```bash
make test
```

## 🐳 Services
- Kafka: `9092`
- Chroma: `8000`
- Worker container runs the consumer (`python src/main.py`)

## ❗ Notes
- If you need a producer, use a different service/repo. This repository is the consumer/processing worker.
- For production Spark clusters, submit with the Kafka package (`spark-sql-kafka-0-10_2.12:3.5.0`). See the commented `spark-submit` in `docker-compose.yml`.
- The multi-worker architecture enables better resource utilization and horizontal scaling for production deployments.

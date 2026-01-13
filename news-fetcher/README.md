# News Fetcher Worker

A stateless Go worker that fetches news from various online sources via API calls and pushes aggregated data to Kafka. Built with object-oriented design principles, containerization, and modern Go practices.

## Features

- **Object-Oriented Design**: Base fetcher interface with concrete implementations for different news sources
- **Multiple News Sources**: Support for NewsAPI, The Guardian, and easily extensible for more sources
- **Kafka Integration**: Unified message structure pushed to Kafka topics
- **Rate Limiting**: Configurable rate limiting for API requests
- **Retry Logic**: Automatic retry with exponential backoff
- **Graceful Shutdown**: Proper cleanup and graceful shutdown handling
- **Docker Support**: Multi-stage Docker build for production deployment
- **Configuration Management**: Environment-based configuration with sensible defaults
- **Structured Logging**: JSON logging with configurable levels

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   News Sources  │───▶│  News Fetcher    │───▶│     Kafka       │
│                 │    │     Worker       │    │                 │
│ • NewsAPI       │    │                  │    │ • news topic    │
│ • Guardian      │    │ • Rate Limiting  │    │ • Unified msg   │
│ • (Extensible)  │    │ • Retry Logic    │    │ • Partitioning  │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

## Quick Start

### Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- API keys for news sources (NewsAPI, Guardian)

### 1. Clone and Setup

```bash
git clone <repository-url>
cd news-fetcher
go mod download
```

### 2. Environment Configuration

Create a `.env` file:

```bash
# News API Keys
export NEWS_API_KEY="your_newsapi_key_here"
export GUARDIAN_API_KEY="your_guardian_key_here"

# Kafka Configuration
export NEWS_FETCHER_KAFKA_BOOTSTRAP_SERVERS="localhost:9092"
export NEWS_FETCHER_KAFKA_TOPIC="news"
```

### 3. Start Kafka (Local Development)

```bash
# Start Kafka and Zookeeper
docker-compose up -d zookeeper kafka

# Wait for Kafka to be ready (about 30 seconds)
docker-compose logs -f kafka
```

### 4. Run Locally (for debugging)

```bash
# Source environment variables
source .env

# Run the worker
go run cmd/worker/main.go
```

### 5. Run with Docker

```bash
# Build and run with Docker Compose
docker-compose up --build
```

## Configuration

The application supports configuration through:

1. **Environment Variables** (recommended for production)
2. **Configuration Files** (YAML format)
3. **Command Line Arguments** (future enhancement)

### Environment Variables

All configuration can be set via environment variables with the `NEWS_FETCHER_` prefix:

```bash
# Kafka Configuration
NEWS_FETCHER_KAFKA_BOOTSTRAP_SERVERS=localhost:9092
NEWS_FETCHER_KAFKA_TOPIC=news
NEWS_FETCHER_KAFKA_COMPRESSION_TYPE=snappy

# Worker Configuration
NEWS_FETCHER_WORKER_FETCH_INTERVAL=1h
NEWS_FETCHER_WORKER_MAX_WORKERS=5

# Logging Configuration
NEWS_FETCHER_LOGGING_LEVEL=info
NEWS_FETCHER_LOGGING_FORMAT=json

# Fetcher Configuration (array format)
NEWS_FETCHER_FETCHERS_0_NAME=newsapi
NEWS_FETCHER_FETCHERS_0_TYPE=newsapi
NEWS_FETCHER_FETCHERS_0_URL=https://newsapi.org/v2/everything
NEWS_FETCHER_FETCHERS_0_API_KEY=your_key
NEWS_FETCHER_FETCHERS_0_ENABLED=true
```

### Configuration File

Create `config/config.yaml`:

```yaml
kafka:
  bootstrap_servers: "localhost:9092"
  topic: "news"
  compression_type: "snappy"
  acks: "all"
  retries: 3

worker:
  fetch_interval: "1h"
  max_workers: 5
  graceful_shutdown_timeout: "30s"

logging:
  level: "info"
  format: "json"

fetchers:
  - name: "newsapi"
    type: "newsapi"
    url: "https://newsapi.org/v2/everything"
    api_key: "${NEWS_API_KEY}"
    rate_limit: "1s"
    timeout: "30s"
    max_retries: 3
    enabled: true

  - name: "guardian"
    type: "guardian"
    url: "https://content.guardianapis.com/search"
    api_key: "${GUARDIAN_API_KEY}"
    rate_limit: "1s"
    timeout: "30s"
    max_retries: 3
    enabled: true
```

## Development and Debugging

### Running Locally (Recommended for Development)

1. **Start Kafka**:
   ```bash
   docker-compose up -d zookeeper kafka
   ```

2. **Set Environment Variables**:
   ```bash
   export NEWS_API_KEY="your_key"
   export GUARDIAN_API_KEY="your_key"
   ```

3. **Run with Debugger**:
   ```bash
   # Using Delve debugger
   dlv debug cmd/worker/main.go

   # Or run directly
   go run cmd/worker/main.go
   ```

### Adding New News Sources

1. **Create Fetcher Implementation**:
   ```go
   // srcfetcher/yournews.go
   type YourNewsFetcher struct {
       *BaseFetcher
       httpClient *HTTPClient
   }

   func (y *YourNewsFetcher) FetchNews(ctx context.Context) ([]NewsArticle, error) {
       // Implementation
   }
   ```

2. **Register Fetcher**:
   ```go
   // srcfetcher/registry.go
   registry.Register("yournews", NewYourNewsFetcher)
   ```

3. **Add Configuration**:
   ```yaml
   fetchers:
     - name: "yournews"
       type: "yournews"
       url: "https://api.yournews.com/articles"
       api_key: "${YOUR_NEWS_API_KEY}"
       enabled: true
   ```

### Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./srcfetcher -v
```

## Monitoring and Observability

### Logs

The application uses structured JSON logging by default:

```json
{
  "level": "info",
  "msg": "Successfully processed 25 articles from newsapi in 2.3s",
  "time": "2024-01-15T10:30:00Z"
}
```

### Health Checks

The Docker container includes a health check that verifies the process is running.

### Kafka Monitoring

Use Kafka UI (included in docker-compose) to monitor messages:
- URL: http://localhost:8080
- View topics, messages, and consumer groups

## Production Deployment

### Docker Deployment

```bash
# Build production image
docker build -t news-fetcher:latest .

# Run with environment variables
docker run -d \
  --name news-fetcher \
  -e NEWS_API_KEY="your_key" \
  -e GUARDIAN_API_KEY="your_key" \
  -e NEWS_FETCHER_KAFKA_BOOTSTRAP_SERVERS="kafka:9092" \
  news-fetcher:latest
```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: news-fetcher
spec:
  replicas: 3
  selector:
    matchLabels:
      app: news-fetcher
  template:
    metadata:
      labels:
        app: news-fetcher
    spec:
      containers:
      - name: news-fetcher
        image: news-fetcher:latest
        env:
        - name: NEWS_API_KEY
          valueFrom:
            secretKeyRef:
              name: news-fetcher-secrets
              key: news-api-key
        - name: GUARDIAN_API_KEY
          valueFrom:
            secretKeyRef:
              name: news-fetcher-secrets
              key: guardian-api-key
```

## Troubleshooting

### Common Issues

1. **Kafka Connection Failed**:
   - Check if Kafka is running: `docker-compose ps`
   - Verify bootstrap servers configuration
   - Check network connectivity

2. **API Rate Limiting**:
   - Increase rate limit intervals in configuration
   - Check API key limits and quotas

3. **No Articles Fetched**:
   - Verify API keys are correct
   - Check API endpoints are accessible
   - Review logs for specific error messages

### Debug Mode

Enable debug logging:

```bash
export NEWS_FETCHER_LOGGING_LEVEL=debug
go run cmd/worker/main.go
```

## API Keys

### NewsAPI
- Sign up at: https://newsapi.org/
- Free tier: 1000 requests/day
- Get API key from dashboard

### The Guardian
- Sign up at: https://open-platform.theguardian.com/
- Free tier: 5000 requests/day
- Get API key from developer section

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.


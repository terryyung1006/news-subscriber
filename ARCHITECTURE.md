# News Subscriber - System Architecture

## Overview

A personalized news subscription service that uses AI to curate and deliver news based on user interests.

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Users (Web)                              │
└────────────────────────┬────────────────────────────────────────┘
                         │ HTTPS
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                   Frontend (Next.js)                             │
│              news-subscriber-frontend                            │
│                                                                   │
│  • User authentication & session management                      │
│  • News preference management UI                                 │
│  • News report display & chatbot interface                       │
└────────────────────────┬────────────────────────────────────────┘
                         │ gRPC
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                 Backend Core (Go gRPC)                           │
│                news-subscriber-core                              │
│                                                                   │
│  • Business logic orchestration                                  │
│  • User & preference management                                  │
│  • News report generation                                        │
│  • Chat interface for report interaction                         │
└──────┬──────────────────┬────────────────┬─────────────────────┘
       │                  │                │
       │ SQL              │ Vector         │ Queue
       ▼                  ▼                ▼
┌─────────────┐    ┌─────────────┐   ┌──────────┐
│ PostgreSQL  │    │  ChromaDB   │   │  Kafka   │
│             │    │             │   │          │
│ • Users     │    │ • News      │   │ • Events │
│ • Prefs     │    │   Embeddings│   │          │
│ • Metadata  │    │             │   │          │
└─────────────┘    └─────────────┘   └────┬─────┘
                                           │
                    ┌──────────────────────┴──────────────────────┐
                    │                                             │
                    ▼                                             ▼
            ┌───────────────┐                            ┌───────────────┐
            │  ETL Worker   │                            │ News Fetcher  │
            │               │                            │               │
            │ • Transform   │                            │ • APIs        │
            │ • Enrich      │                            │ • RSS         │
            │ • Embed       │                            │ • Crawlers    │
            │ • Store       │                            │               │
            └───────────────┘                            └───────────────┘
```

## Components

### 1. news-subscriber-frontend
- **Tech**: Next.js, React, TypeScript
- **Purpose**: Web UI for user interaction
- **Responsibilities**:
  - User authentication (Google OAuth)
  - Interest/preference management (editable cards)
  - Display personalized news reports
  - Chat interface with LLM about reports

### 2. news-subscriber-core  
- **Tech**: Go, gRPC
- **Purpose**: Central business logic engine
- **Responsibilities**:
  - Process all frontend requests
  - User account & preference management
  - Generate personalized news reports
  - Coordinate with databases and services
  - Handle authentication & authorization

### 3. news-fetcher
- **Tech**: Go/Python (to be determined)
- **Purpose**: Fetch news from various sources
- **Responsibilities**:
  - Run cronjobs to query news APIs
  - Aggregate from multiple sources (RSS, APIs, crawlers)
  - Standardize news format
  - Produce to Kafka topic

### 4. etl-worker
- **Tech**: Python
- **Purpose**: Process and store news data
- **Responsibilities**:
  - Consume from Kafka
  - Transform news content
  - Generate embeddings (for semantic search)
  - Store in ChromaDB (vectors) and PostgreSQL (metadata)

### 5. news-subscriber-grpc-proto
- **Tech**: Protocol Buffers
- **Purpose**: API contract definitions
- **Responsibilities**:
  - Central source of truth for all gRPC APIs
  - Shared across all services

## Data Stores

### PostgreSQL
- **Purpose**: Primary relational database
- **Stores**:
  - User accounts
  - User preferences/interests
  - News metadata
  - Invitation codes (during beta)

### ChromaDB
- **Purpose**: Vector database for semantic search
- **Stores**:
  - News article embeddings
  - User interest embeddings
  - Enables similarity-based news matching

### Kafka
- **Purpose**: Event streaming platform
- **Use Cases**:
  - Decouple news fetching from processing
  - Handle high-volume news ingestion
  - Enable event-driven architecture

## Communication Patterns

### Frontend ↔ Backend
- **Protocol**: gRPC (via HTTP proxy in Next.js API routes)
- **Authentication**: Session tokens (JWT in future)
- **Services**: AuthService, UserService, ChatService, ReportService

### Backend ↔ Databases
- **PostgreSQL**: SQL queries via `database/sql`
- **ChromaDB**: HTTP API via Go client library

### News Pipeline
- **news-fetcher** → Kafka → **etl-worker** → ChromaDB/PostgreSQL

## Data Flow Example: Personalized Report

```
1. User logs in, sets interests: "AI, Climate Change, Tech Startups"
   Frontend → Core (UserService)
   
2. Core stores preferences in PostgreSQL
   Core → PostgreSQL
   
3. Daily job triggers report generation
   Core queries ChromaDB for matching news
   Core → ChromaDB (semantic search by interest embeddings)
   
4. Core aggregates and formats report
   Core generates summary (may use LLM)
   
5. User views report on dashboard
   Frontend → Core (ReportService)
   
6. User chats: "Tell me more about the AI breakthrough"
   Frontend → Core (ChatService)
   Core retrieves context, queries LLM, returns response
```

## Security Architecture

- **Authentication**: Google OAuth (beta: with invitation codes)
- **Authorization**: Session tokens, user-scoped data access
- **Data Protection**: 
  - Passwords never stored (OAuth only)
  - User data isolated by user ID
  - Database credentials in environment configs

## Deployment Architecture (Future)

```
Production:
- Frontend: Vercel / AWS CloudFront + S3
- Backend: Kubernetes cluster / AWS ECS
- PostgreSQL: AWS RDS / Cloud SQL
- ChromaDB: Self-hosted on EC2 / GCP Compute
- Kafka: AWS MSK / Confluent Cloud
```

## Design Principles

1. **Microservices**: Each component has a single responsibility
2. **gRPC First**: All inter-service communication via proto-defined APIs
3. **Event-Driven**: News pipeline uses Kafka for decoupling
4. **Database per Concern**: SQL for structured data, Vector DB for embeddings
5. **API Gateway Pattern**: Frontend proxies to backend, shielding gRPC complexity

## Technology Choices Rationale

| Component | Technology | Why? |
|-----------|-----------|------|
| Frontend | Next.js | SSR, API routes, React ecosystem |
| Backend | Go + gRPC | Performance, type safety, efficient RPC |
| News Processing | Python | Rich ML/NLP libraries, easy scripting |
| Vector DB | ChromaDB | Open-source, simple API, good for embeddings |
| Message Queue | Kafka | Industry standard, high throughput |
| Proto Definitions | Centralized Repo | Single source of truth, version control |

## Future Considerations

- **Caching**: Redis for session data and frequent queries
- **CDN**: For static assets and news images
- **Monitoring**: Prometheus + Grafana for metrics
- **Logging**: Centralized logging (ELK stack or similar)
- **Rate Limiting**: API gateway with throttling
- **Horizontal Scaling**: Load balancer + multiple backend instances

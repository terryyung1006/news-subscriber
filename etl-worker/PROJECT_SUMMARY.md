# RAG Financial News Worker - Project Summary

## 🎯 Project Overview

This is a scalable RAG (Retrieval-Augmented Generation) system designed to process financial news and data using Kafka, Apache Spark, and Chroma vector database. The system follows SOLID principles and is built with extensibility in mind.

## 🏗️ Architecture

```
┌─────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Kafka     │───▶│ Apache Spark     │───▶│     Chroma      │
│ (Data Queue)│    │   (Processing)   │    │ (Vector Store)  │
└─────────────┘    └──────────────────┘    └─────────────────┘
```

### Data Flow
1. **Data Ingestion**: Mock financial news data (extensible to real APIs)
2. **Processing**: Text cleaning, sentiment analysis, embedding generation
3. **Storage**: Vector embeddings stored in ChromaDB
4. **Streaming**: Kafka for reliable data flow
5. **Scalability**: Apache Spark for distributed processing

## 📁 Project Structure

```
RAG-worker/
├── src/
│   ├── data_sources/          # Abstract data source interfaces
│   │   ├── base.py           # DataSource ABC and NewsItem model
│   │   ├── mock_source.py    # Mock financial news implementation
│   │   └── factory.py        # Factory pattern for data sources
│   ├── processors/           # Data processing logic
│   │   ├── base.py          # DataProcessor ABC
│   │   ├── news_processor.py # Text cleaning and preprocessing
│   │   ├── embedding_processor.py # Vector embedding generation
│   │   └── sentiment_processor.py # Sentiment analysis
│   ├── vector_storage/       # Vector database operations
│   │   ├── base.py          # VectorStorage ABC
│   │   ├── chroma_storage.py # ChromaDB implementation
│   │   └── factory.py       # Factory for vector storage
│   ├── kafka/               # Kafka integration
│   │   └── producer.py      # Kafka producer for data streaming
│   └── main.py              # Main application orchestrator
├── tests/                   # Unit tests
├── config/                  # Configuration files
├── docker/                  # Docker configurations
├── data/                    # Sample data
└── logs/                    # Application logs
```

## 🧩 SOLID Principles Implementation

### 1. Single Responsibility Principle (SRP)
- Each class has one reason to change
- `DataSource` handles only data retrieval
- `DataProcessor` handles only data processing
- `VectorStorage` handles only vector operations

### 2. Open/Closed Principle (OCP)
- Open for extension, closed for modification
- New data sources can be added without changing existing code
- New processors can be implemented by extending base classes
- New vector storage backends can be added via factory pattern

### 3. Liskov Substitution Principle (LSP)
- Derived classes can substitute base classes
- All data sources implement the same `DataSource` interface
- All processors implement the same `DataProcessor` interface
- All vector storage implementations follow the same contract

### 4. Interface Segregation Principle (ISP)
- Clients depend only on interfaces they use
- Separate interfaces for different concerns
- No forced dependencies on unused methods

### 5. Dependency Inversion Principle (DIP)
- High-level modules don't depend on low-level modules
- Both depend on abstractions
- Factory patterns used for dependency injection

## 🚀 Key Features

### Data Sources
- **Mock Financial News**: Realistic test data for development
- **Extensible**: Easy to add new data sources (APIs, RSS, databases)
- **Factory Pattern**: Centralized creation of data sources

### Data Processing
- **Text Cleaning**: Normalization, keyword extraction, summarization
- **Sentiment Analysis**: Rule-based financial sentiment analysis
- **Embedding Generation**: Sentence transformers for vector embeddings
- **Modular**: Each processor can be used independently

### Vector Storage
- **ChromaDB Integration**: Efficient similarity search
- **Metadata Support**: Rich metadata for filtering and retrieval
- **Scalable**: Can be extended to other vector databases

### Streaming
- **Kafka Integration**: Reliable message queuing
- **Producer/Consumer Pattern**: Decoupled data flow
- **Fault Tolerant**: Handles connection failures gracefully

## 🐳 Docker Support

### Services
- **Zookeeper**: Kafka coordination
- **Kafka**: Message queuing
- **ChromaDB**: Vector database
- **Spark Master/Worker**: Distributed processing
- **Kafka UI**: Monitoring interface
- **RAG Worker**: Main application

### Quick Start
```bash
# Start all services
./start.sh

# Or manually
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f rag-worker
```

## 🔧 Configuration

### Key Configuration Options
- **Kafka**: Topics, brokers, consumer groups
- **ChromaDB**: Host, port, collection settings
- **Processing**: Batch sizes, embedding models, thresholds
- **Data Sources**: Update intervals, news counts
- **Logging**: Levels, formats, file paths

### Environment Variables
- All sensitive configuration can be externalized
- Docker environment variables supported
- Local development vs production configurations

## 🧪 Testing

### Test Coverage
- **Unit Tests**: All major components tested
- **Integration Tests**: End-to-end pipeline testing
- **Mock Data**: Realistic financial news for testing
- **Factory Tests**: Factory pattern validation

### Running Tests
```bash
# Run all tests
python run_tests.py

# Run specific tests
pytest tests/test_data_sources.py -v

# Code quality checks
flake8 src/
mypy src/
```

## 📈 Extensibility

### Adding New Data Sources
1. Implement `DataSource` interface
2. Register in `DataSourceFactory`
3. Add configuration
4. No changes to existing code required

### Adding New Processors
1. Extend `DataProcessor` base class
2. Implement required methods
3. Register in processing pipeline
4. Configuration-driven activation

### Adding New Vector Storage
1. Implement `VectorStorage` interface
2. Register in `VectorStorageFactory`
3. Update configuration
4. Switch storage backends seamlessly

## 🔍 Monitoring & Observability

### Logging
- Structured logging with loguru
- Multiple output formats
- Log rotation and retention
- Different log levels for different environments

### Health Checks
- Docker health checks
- Service connectivity monitoring
- Application health endpoints
- Graceful degradation

### Metrics
- Processing throughput
- Error rates
- Response times
- Resource utilization

## 🛡️ Error Handling

### Robust Error Handling
- Connection retries with exponential backoff
- Graceful degradation on service failures
- Comprehensive error logging
- Circuit breaker patterns for external services

### Data Validation
- Pydantic models for data validation
- Type checking with mypy
- Input sanitization
- Output validation

## 📚 Usage Examples

### Basic Usage
```python
from src.main import RAGWorker

# Start the RAG worker
worker = RAGWorker()
worker.start()
```

### Custom Configuration
```python
# Load custom configuration
worker = RAGWorker(config_path="config/custom.yaml")
worker.start()
```

### Adding Custom Data Source
```python
from src.data_sources.base import DataSource
from src.data_sources.factory import DataSourceFactory

class CustomDataSource(DataSource):
    # Implementation here
    pass

# Register new data source
factory = DataSourceFactory()
factory.register_source("custom", CustomDataSource)
```

## 🚀 Deployment

### Production Deployment
1. **Environment Setup**: Configure production environment variables
2. **Resource Allocation**: Adjust Docker resource limits
3. **Monitoring**: Set up monitoring and alerting
4. **Scaling**: Configure Spark worker scaling
5. **Backup**: Set up data backup strategies

### Development Setup
1. **Clone Repository**: `git clone <repo-url>`
2. **Install Dependencies**: `pip install -r requirements.txt`
3. **Start Services**: `./start.sh`
4. **Run Tests**: `python run_tests.py`
5. **Monitor**: Check service URLs for monitoring

## 🔮 Future Enhancements

### Planned Features
- **Real-time APIs**: Integration with financial news APIs
- **Advanced ML**: More sophisticated sentiment analysis
- **Multi-modal**: Support for images and charts
- **Caching**: Redis integration for performance
- **API Gateway**: REST API for external access
- **Dashboard**: Web-based monitoring dashboard

### Scalability Improvements
- **Kubernetes**: Container orchestration
- **Auto-scaling**: Dynamic resource allocation
- **Multi-region**: Geographic distribution
- **Load Balancing**: Traffic distribution
- **Caching Layers**: Multi-level caching

## 📄 License

MIT License - See LICENSE file for details.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Follow SOLID principles
4. Add tests for new features
5. Submit a pull request

## 📞 Support

For questions and support:
- Create an issue in the repository
- Check the documentation
- Review the test examples
- Examine the configuration options

---

**Built with ❤️ following SOLID principles for scalable, maintainable, and extensible code.**

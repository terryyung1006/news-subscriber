"""
Data sources package for RAG Financial News Worker.

This package contains abstract interfaces and concrete implementations
for different data sources following SOLID principles.
"""

from .base import DataSource, NewsItem
from .factory import DataSourceFactory
from .mock_source import MockFinancialNewsSource

__all__ = ["DataSource", "NewsItem", "MockFinancialNewsSource", "DataSourceFactory"]

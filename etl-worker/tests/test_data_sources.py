"""
Unit tests for data sources.

This module contains tests for the data source implementations
to ensure they follow the expected interfaces and behavior.
"""

from datetime import datetime

import pytest

from src.data_sources import DataSourceFactory, MockFinancialNewsSource
from src.data_sources.base import NewsItem


class TestMockFinancialNewsSource:
    """Test cases for MockFinancialNewsSource."""

    def test_initialization(self):
        """Test that MockFinancialNewsSource initializes correctly."""
        source = MockFinancialNewsSource(news_count=5)
        assert source.news_count == 5
        assert not source.is_connected()

    def test_connection(self):
        """Test connection functionality."""
        source = MockFinancialNewsSource()
        assert source.connect()
        assert source.is_connected()

        source.disconnect()
        assert not source.is_connected()

    def test_fetch_news(self):
        """Test news fetching functionality."""
        source = MockFinancialNewsSource(news_count=3)
        source.connect()

        news_items = source.fetch_news()
        assert len(news_items) == 3
        assert all(isinstance(item, NewsItem) for item in news_items)

        # Test with limit
        limited_items = source.fetch_news(limit=2)
        assert len(limited_items) == 2

    def test_news_item_structure(self):
        """Test that news items have the correct structure."""
        source = MockFinancialNewsSource(news_count=1)
        source.connect()

        news_items = source.fetch_news()
        assert len(news_items) == 1

        item = news_items[0]
        assert hasattr(item, "id")
        assert hasattr(item, "title")
        assert hasattr(item, "content")
        assert hasattr(item, "source")
        assert hasattr(item, "published_date")
        assert hasattr(item, "category")
        assert hasattr(item, "metadata")

        assert isinstance(item.id, str)
        assert isinstance(item.title, str)
        assert isinstance(item.content, str)
        assert isinstance(item.published_date, datetime)
        assert isinstance(item.metadata, dict)

    def test_source_name(self):
        """Test that source name is correct."""
        source = MockFinancialNewsSource()
        assert source.source_name == "Mock Financial News"


class TestDataSourceFactory:
    """Test cases for DataSourceFactory."""

    def test_factory_initialization(self):
        """Test that factory initializes with default sources."""
        factory = DataSourceFactory()
        available_sources = factory.get_available_sources()
        assert "mock" in available_sources

    def test_create_mock_source(self):
        """Test creating mock data source."""
        factory = DataSourceFactory()
        source = factory.create_source("mock", news_count=5)

        assert isinstance(source, MockFinancialNewsSource)
        assert source.news_count == 5

    def test_create_invalid_source(self):
        """Test that factory raises error for invalid source type."""
        factory = DataSourceFactory()

        with pytest.raises(ValueError, match="Unknown data source type"):
            factory.create_source("invalid_source")

    def test_register_new_source(self):
        """Test registering a new data source type."""
        factory = DataSourceFactory()

        # Create a mock source class for testing
        class TestSource(MockFinancialNewsSource):
            pass

        factory.register_source("test", TestSource)
        assert factory.has_source("test")

        source = factory.create_source("test", news_count=3)
        assert isinstance(source, TestSource)

    def test_get_available_sources(self):
        """Test getting list of available sources."""
        factory = DataSourceFactory()
        sources = factory.get_available_sources()

        assert isinstance(sources, list)
        assert "mock" in sources

    def test_has_source(self):
        """Test checking if source type exists."""
        factory = DataSourceFactory()

        assert factory.has_source("mock")
        assert not factory.has_source("nonexistent")


if __name__ == "__main__":
    pytest.main([__file__])

"""
Base interfaces for data sources following SOLID principles.

This module defines the abstract interfaces that all data sources must implement,
ensuring consistency and allowing for easy extension.
"""

from abc import ABC, abstractmethod
from datetime import datetime
from typing import Any, Dict, List, Optional

from pydantic import BaseModel, Field


class NewsItem(BaseModel):
    """Data model for financial news items."""

    id: str = Field(..., description="Unique identifier for the news item")
    title: str = Field(..., description="News headline")
    content: str = Field(..., description="Full news content")
    source: str = Field(..., description="News source (e.g., Reuters, Bloomberg)")
    published_date: datetime = Field(..., description="Publication date and time")
    category: str = Field(..., description="News category (e.g., stocks, bonds, forex)")
    sentiment: Optional[str] = Field(None, description="Sentiment analysis result")
    url: Optional[str] = Field(None, description="Original news URL")
    metadata: Dict[str, Any] = Field(
        default_factory=dict, description="Additional metadata"
    )


class DataSource(ABC):
    """
    Abstract base class for data sources.

    This interface follows the Single Responsibility Principle by focusing
    solely on data retrieval. It also follows the Open/Closed Principle
    by allowing extension through inheritance without modification.
    """

    @abstractmethod
    def connect(self) -> bool:
        """
        Establish connection to the data source.

        Returns:
            bool: True if connection successful, False otherwise
        """
        pass

    @abstractmethod
    def disconnect(self) -> None:
        """Close connection to the data source."""
        pass

    @abstractmethod
    def fetch_news(self, limit: Optional[int] = None) -> List[NewsItem]:
        """
        Fetch news items from the data source.

        Args:
            limit: Maximum number of news items to fetch

        Returns:
            List[NewsItem]: List of news items
        """
        pass

    @abstractmethod
    def is_connected(self) -> bool:
        """
        Check if the data source is connected.

        Returns:
            bool: True if connected, False otherwise
        """
        pass

    @property
    @abstractmethod
    def source_name(self) -> str:
        """
        Get the name of the data source.

        Returns:
            str: Name of the data source
        """
        pass

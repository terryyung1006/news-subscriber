"""
Factory for creating data sources.

This module implements the Factory pattern to create different types of
data sources, following the Dependency Inversion Principle.
"""

from typing import Dict, Type

from loguru import logger

from .base import DataSource
from .mock_source import MockFinancialNewsSource


class DataSourceFactory:
    """
    Factory for creating data source instances.

    This class follows the Open/Closed Principle by allowing new data sources
    to be added without modifying existing code.
    """

    def __init__(self):
        """Initialize the factory with available data source types."""
        self._sources: Dict[str, Type[DataSource]] = {
            "mock": MockFinancialNewsSource,
            # Future data sources can be added here:
            # "api": APIFinancialNewsSource,
            # "rss": RSSFinancialNewsSource,
            # "database": DatabaseFinancialNewsSource,
        }

    def register_source(self, name: str, source_class: Type[DataSource]) -> None:
        """
        Register a new data source type.

        Args:
            name: Name identifier for the data source
            source_class: Class implementing the DataSource interface
        """
        self._sources[name] = source_class
        logger.info(f"Registered data source: {name}")

    def create_source(self, source_type: str, **kwargs) -> DataSource:
        """
        Create a data source instance.

        Args:
            source_type: Type of data source to create
            **kwargs: Additional arguments for the data source constructor

        Returns:
            DataSource: Instance of the requested data source

        Raises:
            ValueError: If the source type is not registered
        """
        if source_type == "mock":
            news_count = kwargs.get("news_count", 10)
            use_llm = kwargs.get("use_llm", True)  # Default to using LLM
            model_name = kwargs.get("model_name", "llama2")  # Default to llama2
            return MockFinancialNewsSource(
                news_count=news_count, use_llm=use_llm, model_name=model_name
            )
        if source_type not in self._sources:
            available_sources = ", ".join(self._sources.keys())
            raise ValueError(
                f"Unknown data source type: {source_type}. "
                f"Available types: {available_sources}"
            )

        source_class = self._sources[source_type]
        logger.info(f"Creating data source: {source_type}")

        return source_class(**kwargs)

    def get_available_sources(self) -> list[str]:
        """
        Get list of available data source types.

        Returns:
            list[str]: List of registered data source names
        """
        return list(self._sources.keys())

    def has_source(self, source_type: str) -> bool:
        """
        Check if a data source type is registered.

        Args:
            source_type: Name of the data source type

        Returns:
            bool: True if the source type is registered, False otherwise
        """
        return source_type in self._sources

"""
Base interfaces for data processors following SOLID principles.

This module defines the abstract interfaces that all data processors must implement,
ensuring consistency and allowing for easy extension.
"""

from typing import Protocol

from pyspark.sql import DataFrame


class StreamProcessor(Protocol):
    """
    Abstract base class for data processors.

    This interface follows the Single Responsibility Principle by focusing
    solely on data processing. It also follows the Open/Closed Principle
    by allowing extension through inheritance without modification.
    """

    def process(self, batch_df: DataFrame) -> DataFrame:
        """
        Process a batch of data.

        Args:
            batch_df: DataFrame to process

        Returns:
            DataFrame: Processed DataFrame
        """
        pass

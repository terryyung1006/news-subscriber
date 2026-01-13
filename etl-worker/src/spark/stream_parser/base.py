from typing import Protocol

from pyspark.sql import DataFrame


class StreamParser(Protocol):
    """Abstract base class for stream parsing implementations."""

    def start_streaming(self, topic: str) -> DataFrame:
        """
        Start streaming from the specified topic and return parsed DataFrame.

        Args:
            topic: Name of the topic to stream from

        Returns:
            DataFrame: Parsed Spark DataFrame with structured data
        """
        pass

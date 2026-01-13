#!/usr/bin/env python3
"""
Spark streaming application that consumes Kafka data,
processes it with multiple stream processors, and writes to storage using write_streamer.
"""

from typing import Any, Dict, List

from pyspark.sql import DataFrame, SparkSession
from pyspark.sql.streaming import StreamingQuery

from spark.stream_parser.base import StreamParser
from spark.stream_processors.base import StreamProcessor
from spark.write_streamer.base import SparkWriteStreamer


class SparkStreamingJob:
    """
    Main orchestrator for Spark streaming job.

    This class properly orchestrates the streaming pipeline:
    1. Parse data from Kafka
    2. Apply multiple stream processors in sequence (text processing, embeddings, etc.)
    3. Use write_streamer to write processed data to storage
    """

    def __init__(
        self,
        spark: SparkSession,
        stream_parser: StreamParser,
        stream_processors: List[StreamProcessor],
        write_streamers: List[SparkWriteStreamer],
    ):
        """
        Initialize the Spark streaming job.

        Args:
            spark: SparkSession instance
            stream_parser: Parser for incoming data streams
            stream_processors: List of processors for transforming the data in sequence
            write_streamers: Streamers for writing processed data to storage
        """
        self.spark = spark
        self.stream_parser = stream_parser
        self.stream_processors = stream_processors
        self.write_streamers = write_streamers
        self.queries: List[StreamingQuery] = []

    def start(self, topic: str = "text-data") -> None:
        """Start the streaming application."""
        print(f"Starting Spark streaming from Kafka topic: {topic}")

        # 1. Parse data from Kafka
        stream_df = self.stream_parser.start_streaming(topic)

        # 2. Apply multiple stream processors in sequence
        for i, processor in enumerate(self.stream_processors):
            print(f"Applying processor {i}: {processor.__class__.__name__}")
            stream_df = processor.process(stream_df)

        # # 3. Use write_streamer to write processed data to storage
        for write_streamer in self.write_streamers:
            self.queries.append(write_streamer.write(stream_df))

        print("Streaming started. Waiting for data...")

        try:
            # Keep the streaming query alive
            for query in self.queries:
                query.awaitTermination()
        except KeyboardInterrupt:
            print("Stopping streaming...")
            self.stop()
        finally:
            self.spark.stop()

    def stop(self) -> None:
        """Stop the streaming query."""
        if self.queries:
            for query in self.queries:
                if query.isActive:
                    query.stop()
            print("Streaming query stopped")

    def get_query(self) -> List[StreamingQuery]:
        """Get the current streaming query."""
        return self.queries


# Legacy class name for backward compatibility
class SparkStreamingApp(SparkStreamingJob):
    """Legacy class name - use SparkStreamingJob instead."""

    pass


def main():
    """Main function to run the Spark streaming application."""
    # This would need proper initialization with all dependencies
    print("Please use the main.py entry point instead")
    print("This file is meant to be imported, not run directly")


if __name__ == "__main__":
    main()

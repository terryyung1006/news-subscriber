"""
Main entry point for the consumer Spark job.

This service:
- Consumes messages from Kafka via Spark Structured Streaming
- Processes with multiple processors (cleaning, sentiment, embeddings)
- Writes to Chroma vector database using write_streamer
"""

import os
import sys

from loguru import logger
from pyspark.sql import SparkSession
from pyspark.sql.types import (
    ArrayType,
    DoubleType,
    LongType,
    StringType,
    StructField,
    StructType,
)

from load_config import get_config
from spark import SparkStreamingJob
from spark.stream_parser.kafka_stream_parser import KafkaStreamParser
from spark.stream_processors.chunk_processor import ChunkRowProcessor
from spark.stream_processors.embedding_processor import Embedder
from spark.stream_processors.length_processor import LengthProcessor
from spark.write_streamer.chroma_writer_streamer import ChromaWriterStreamer
from spark.write_streamer.console_writer import ConsoleWriter
from storage.chroma_storage import ChromaVectorStorage

# Define schema for incoming Kafka messages (Unified Schema)
schema = StructType(
    [
        StructField("id", StringType(), True),
        StructField("title", StringType(), True),
        StructField("content", StringType(), True),
        StructField("summary", StringType(), True),
        StructField("url", StringType(), True),
        StructField("image_url", StringType(), True),
        StructField("source", StringType(), True),
        StructField("author", StringType(), True),
        StructField("published_at", StringType(), True),  # Timestamp as string or TimestampType
        StructField("category", StringType(), True),
        StructField("tags", ArrayType(StringType()), True),
        StructField("fetched_at", StringType(), True),
        StructField("fetcher_name", StringType(), True),
    ]
)


class App:
    def __init__(self) -> None:
        self.config = get_config()

    def run(self) -> None:
        import os
        import sys

        # Set environment variables for CPU-only operations
        os.environ["CUDA_VISIBLE_DEVICES"] = ""
        os.environ["TORCH_USE_CUDA_DSA"] = "0"

        print(f"Driver Python version: {sys.version}")
        print(f"Driver Python executable: {sys.executable}")
        spark = (
            SparkSession.builder.appName("SparkStreamingExample")
            .config("spark.sql.streaming.checkpointLocation", "/tmp/checkpoint")
            .config("spark.driver.bindAddress", "127.0.0.1")
            .config("spark.driver.host", "127.0.0.1")
            .config("spark.sql.adaptive.enabled", "false")
            .config("spark.sql.adaptive.coalescePartitions.enabled", "false")
            .config(
                "spark.jars.packages",
                "org.apache.spark:spark-sql-kafka-0-10_2.13:4.0.0",
            )
            .config("spark.sql.streaming.forceDeleteTempCheckpointLocation", "true")
            .config("spark.sql.execution.arrow.pyspark.enabled", "false")
            .config("spark.sql.execution.arrow.pyspark.fallback.enabled", "false")
            .config("spark.driver.extraJavaOptions", "-Dfile.encoding=UTF-8")
            .config("spark.executor.extraJavaOptions", "-Dfile.encoding=UTF-8")
            .config("spark.sql.execution.pyspark.udf.faulthandler.enabled", "true")
            .config("spark.python.worker.faulthandler.enabled", "true")
            .config("spark.python.worker.reuse", "false")
            .config("spark.serializer", "org.apache.spark.serializer.KryoSerializer")
            .config("spark.executorEnv.CUDA_VISIBLE_DEVICES", "")
            .config("spark.executorEnv.TORCH_USE_CUDA_DSA", "0")
            .master("local[*]")
            .getOrCreate()
        )

        kafka_stream_parser = KafkaStreamParser(
            spark, self.config.get_string("KAFKA_BOOTSTRAP_SERVERS"), schema
        )

        # Initialize storage and write_streamer
        chroma_storage = ChromaVectorStorage(
            self.config.get_string("CHROMA_HOST"), self.config.get_int("CHROMA_PORT")
        )
        write_streamers = [ConsoleWriter(), ChromaWriterStreamer(chroma_storage)]

        # Create array of stream processors in processing order
        stream_processors = [
            LengthProcessor(),  # First: compute length from content
            ChunkRowProcessor(self.config.get_int("PROCESSING_BATCH_SIZE", 100)),  # Second: chunk the text
            Embedder(
                self.config.get_string("CHROMA_EMBEDDING_MODEL", "all-MiniLM-L6-v2"), "chunks"
            ),  # Third: generate embeddings from chunks
        ]

        # Create and run the streaming job with multiple processors
        job = SparkStreamingJob(
            spark, kafka_stream_parser, stream_processors, write_streamers
        )

        try:
            job.start()
        except KeyboardInterrupt:
            logger.info("Shutting down...")
        finally:
            job.stop()


def main() -> None:
    # Set Python executable for PySpark driver and workers
    os.environ["PYSPARK_DRIVER_PYTHON"] = sys.executable
    os.environ["PYSPARK_PYTHON"] = sys.executable

    # Set environment variables for CPU-only operations
    os.environ["CUDA_VISIBLE_DEVICES"] = ""
    os.environ["TORCH_USE_CUDA_DSA"] = "0"

    App().run()


if __name__ == "__main__":
    main()

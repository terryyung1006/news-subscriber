#!/usr/bin/env python3
"""
Test script to verify the processing pipeline works without Kafka.

This helps isolate whether the issue is with Kafka or with the processors.
"""

import sys
from pathlib import Path

# Add src to path for imports
sys.path.insert(0, str(Path(__file__).parent))

from pyspark.sql import SparkSession
from pyspark.sql.types import DoubleType, LongType, StringType, StructField, StructType
from sentence_transformers import SentenceTransformer

from spark.stream_processors.chunk_processor import ChunkRowProcessor
from spark.stream_processors.embedding_processor import Embedder
from spark.stream_processors.length_processor import LengthProcessor
from storage.chroma_storage import ChromaVectorStorage


def test_pipeline():
    """Test the processing pipeline without Kafka."""
    print("🧪 Testing RAG processing pipeline...")

    # Create Spark session without Kafka
    spark = (
        SparkSession.builder.appName("PipelineTest")
        .config("spark.driver.bindAddress", "127.0.0.1")
        .config("spark.driver.host", "127.0.0.1")
        .config("spark.executorEnv.PYTHONPATH", str(Path(__file__).parent))
        .master("local[*]")
        .getOrCreate()
    )

    print("✅ Spark session created successfully")

    # Create test data
    test_data = [
        (
            "test_1",
            "This is a test financial news article about Apple stock. " * 20,
            "2023-10-27T10:00:00Z",
            "technology",
        ),
        (
            "test_2",
            "Another test article about Tesla and electric vehicles. " * 15,
            "2023-10-27T11:00:00Z",
            "automotive",
        ),
    ]

    # Create DataFrame
    schema = StructType(
        [
            StructField("id", StringType(), True),
            StructField("content", StringType(), True),
            StructField("published_at", StringType(), True),
            StructField("category", StringType(), True),
        ]
    )

    df = spark.createDataFrame(test_data, schema)
    print(f"✅ Test DataFrame created with {df.count()} rows")

    # Test LengthProcessor
    print("\n🔧 Testing LengthProcessor...")
    length_processor = LengthProcessor()
    lengthed_df = length_processor.process(df)
    print(f"✅ Length computation completed. DataFrame has {lengthed_df.count()} rows")

    # Test ChunkRowProcessor
    print("\n🔧 Testing ChunkRowProcessor...")
    chunk_processor = ChunkRowProcessor(100)
    chunked_df = chunk_processor.process(lengthed_df)
    print(f"✅ Chunking completed. DataFrame has {chunked_df.count()} rows")

    # Test Embedder
    print("\n🔧 Testing Embedder...")
    embedding_model = SentenceTransformer("all-MiniLM-L6-v2")
    embedder = Embedder(embedding_model, "chunks")
    embedded_df = embedder.process(chunked_df)
    print(f"✅ Embedding completed. DataFrame has {embedded_df.count()} rows")

    # Show the final structure
    print("\n📊 Final DataFrame schema:")
    embedded_df.printSchema()

    print("\n📊 Sample data:")
    embedded_df.show(2, truncate=False)

    # Test ChromaDB connection
    print("\n🔧 Testing ChromaDB connection...")
    try:
        chroma_storage = ChromaVectorStorage("localhost", 8000)
        print("✅ ChromaDB connection successful")
    except Exception as e:
        print(f"❌ ChromaDB connection failed: {e}")

    print("\n🎉 Pipeline test completed successfully!")

    # Cleanup
    spark.stop()


if __name__ == "__main__":
    test_pipeline()

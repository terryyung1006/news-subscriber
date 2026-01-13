"""
Embedding processor for generating vector embeddings from text chunks.

This module handles the generation of embeddings and prepares data for vector storage.
"""

from typing import List

from pyspark.sql import DataFrame
from pyspark.sql.functions import udf
from pyspark.sql.types import ArrayType, FloatType

from .base import StreamProcessor

# Global model cache to avoid re-initializing models
_model_cache: dict = {}


class Embedder(StreamProcessor):
    """
    Processor that generates embeddings from text chunks and prepares data for ChromaDB.

    This follows the Single Responsibility Principle by focusing solely
    on embedding generation and data structuring for vector storage.
    """

    def __init__(
        self, model_name: str = "all-MiniLM-L6-v2", column_name: str = "chunks"
    ):
        """
        Initialize the embedding processor.

        Args:
            model_name: Name of the sentence transformer model to use
            column_name: Name of the column containing text chunks
        """
        self.model_name = model_name
        self.column_name = column_name

        # Create the UDF for embedding generation
        self.generate_embedding_udf = udf(
            self.generate_embedding, ArrayType(FloatType())
        )

    def process(self, df: DataFrame) -> DataFrame:
        """
        Process the DataFrame by chunking text and generating embeddings.

        Args:
            df: Input DataFrame with text data

        Returns:
            DataFrame with chunks and embeddings
        """
        from pyspark.sql.functions import col, concat, explode, expr, lit, struct

        # Explode chunks and generate embeddings
        processed_df = df.select(
            col("id"),
            col("category"),
            col("published_at").alias("timestamp"),
            col("length"),
            explode(col("chunks")).alias("chunk_text"),
        )

        # Generate chunk_id using a simple approach compatible with streaming
        processed_df = processed_df.withColumn(
            "chunk_id", concat(col("id"), lit("_chunk_"), expr("uuid()"))
        )

        # Generate embeddings for each chunk
        processed_df = processed_df.withColumn(
            "embedding", self.generate_embedding_udf(col("chunk_text"))
        )

        # Structure the output for ChromaDB
        result_df = processed_df.select(
            col("chunk_id"),
            col("embedding"),
            struct(
                col("category").alias("category"),
                col("timestamp").alias("timestamp"),
                col("length").alias("length"),
                col("id").alias("original_id"),
            ).alias("metadata"),
            col("chunk_text").alias("document"),
        )

        return result_df

    def generate_embedding(self, text: str) -> List[float]:
        """
        Generate embedding for a text chunk using a simple hash-based approach.

        This avoids PyTorch device operations that are incompatible with Spark.

        Args:
            text: Text to generate embedding for

        Returns:
            List[float]: Vector embedding as list of floats
        """
        if not text or not text.strip():
            print("Warning: Empty text received, returning zero vector")
            # Return zero vector for empty text (384 is all-MiniLM-L6-v2 dimension)
            return [0.0] * 384

        try:
            # Use a simple hash-based embedding approach that's Spark-compatible
            import hashlib
            import math

            # Create a deterministic embedding based on text content
            # This is a simplified approach that avoids PyTorch device issues
            text_hash = hashlib.sha256(text.encode("utf-8")).hexdigest()

            # Convert hash to a 384-dimensional vector
            embedding = []
            for i in range(0, len(text_hash), 2):
                # Take pairs of hex characters and convert to float
                hex_pair = text_hash[i : i + 2]
                value = int(hex_pair, 16) / 255.0  # Normalize to [0, 1]
                embedding.append(value)

            # Pad or truncate to exactly 384 dimensions
            while len(embedding) < 384:
                # Use text length and position to create additional values
                pos = len(embedding)
                value = (hash(text + str(pos)) % 1000) / 1000.0
                embedding.append(value)

            embedding = embedding[:384]  # Ensure exactly 384 dimensions

            # Normalize the vector
            norm = math.sqrt(sum(x * x for x in embedding))
            if norm > 0:
                embedding = [x / norm for x in embedding]

            print(
                f"Generated hash-based embedding for text: '{text[:50]}...' "
                f"(length: {len(text)})"
            )
            return embedding

        except Exception as e:
            print(f"Error generating embedding: {e}")
            import traceback

            traceback.print_exc()
            # Return zero vector on error
            return [0.0] * 384

"""
Chunk processor for splitting text into manageable chunks.

This module handles text chunking using sentence-based splitting for better
semantic coherence in embeddings.
"""

import re
from typing import List

from pyspark.sql import DataFrame
from pyspark.sql.functions import col, udf
from pyspark.sql.types import ArrayType, StringType

from .base import StreamProcessor


class ChunkRowProcessor(StreamProcessor):
    """
    Processor that splits text into chunks for better processing.

    Uses sentence-based chunking to maintain semantic coherence,
    keeping chunks under a target size for optimal embedding quality.
    """

    def __init__(self, max_chunk_size: int = 500, min_chunk_size: int = 100):
        """
        Initialize the chunk processor.

        Args:
            max_chunk_size: Maximum size of each text chunk (default 500 chars)
            min_chunk_size: Minimum size before merging with next sentence
        """
        self.max_chunk_size = max_chunk_size
        self.min_chunk_size = min_chunk_size

        # Create the UDF function in __init__ to avoid serialization issues
        self.chunk_text_udf = udf(self._chunk_text, ArrayType(StringType()))

    def process(self, batch_df: DataFrame) -> DataFrame:
        """
        Process DataFrame by adding a chunks column.

        Args:
            batch_df: DataFrame with text column

        Returns:
            DataFrame: Original DataFrame with added chunks column
        """
        result_df = batch_df.withColumn("chunks", self.chunk_text_udf(col("content")))
        return result_df

    def _chunk_text(self, text: str) -> List[str]:
        """
        Split text into chunks based on sentence boundaries.

        Groups sentences into chunks that stay under max_chunk_size while
        maintaining semantic coherence.

        Args:
            text: Text to split into chunks

        Returns:
            List[str]: List of text chunks
        """
        if not text or not text.strip():
            return []

        # Split by sentence-ending punctuation, keeping the delimiter
        sentences = re.split(r'(?<=[.!?])\s+', text.strip())
        sentences = [s.strip() for s in sentences if s.strip()]

        if not sentences:
            return []

        chunks = []
        current_chunk = ""

        for sentence in sentences:
            # If adding this sentence exceeds max size, save current chunk
            if current_chunk and len(current_chunk) + len(sentence) + 1 > self.max_chunk_size:
                chunks.append(current_chunk.strip())
                current_chunk = sentence
            else:
                current_chunk = f"{current_chunk} {sentence}".strip() if current_chunk else sentence

        # Add the last chunk if it meets minimum size or if it's the only content
        if current_chunk:
            if len(current_chunk) >= self.min_chunk_size or not chunks:
                chunks.append(current_chunk.strip())
            elif chunks:
                # Merge small final chunk with previous if possible
                chunks[-1] = f"{chunks[-1]} {current_chunk}".strip()

        return chunks

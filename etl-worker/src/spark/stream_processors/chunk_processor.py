"""
Chunk processor for splitting text into manageable chunks.

This module handles text chunking with configurable size and overlap.
"""

from typing import List

from pyspark.sql import DataFrame
from pyspark.sql.functions import col, udf
from pyspark.sql.types import ArrayType, StringType

from .base import StreamProcessor


class ChunkRowProcessor(StreamProcessor):
    """
    Processor that splits text into chunks for better processing.

    This follows the Single Responsibility Principle by focusing solely
    on text chunking operations.
    """

    def __init__(self, chunk_size: int, chunk_overlap: int = 50):
        """
        Initialize the chunk processor.

        Args:
            chunk_size: Size of each text chunk
            chunk_overlap: Overlap between consecutive chunks
        """
        self.chunk_size = chunk_size
        self.chunk_overlap = chunk_overlap

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
        # Apply the UDF - this is safe for streaming DataFrames
        result_df = batch_df.withColumn("chunks", self.chunk_text_udf(col("content")))

        return result_df

    def _chunk_text(self, text: str) -> List[str]:
        """
        Split text into chunks with specified size and overlap.

        Args:
            text: Text to split into chunks

        Returns:
            List[str]: List of text chunks
        """

        if not text or not text.strip():
            return []

        chunks = []
        text_length = len(text)
        current_position = 0

        while current_position < text_length:
            end_position = min(current_position + self.chunk_size, text_length)
            chunk = text[current_position:end_position]
            chunks.append(chunk)

            # Move to the next chunk's starting position
            if end_position < text_length:
                current_position += self.chunk_size - self.chunk_overlap
            else:
                break

        print(f" _chunk_text: Returning {len(chunks)} chunks")
        return chunks

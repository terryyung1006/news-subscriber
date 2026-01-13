"""
Length processor for computing text length.
"""

from pyspark.sql import DataFrame
from pyspark.sql.functions import col, length

from .base import StreamProcessor


class LengthProcessor(StreamProcessor):
    """
    Processor that computes the length of the content.
    """

    def __init__(self, input_col: str = "content", output_col: str = "length"):
        """
        Initialize the length processor.

        Args:
            input_col: Name of the input column (default: "content")
            output_col: Name of the output column (default: "length")
        """
        self.input_col = input_col
        self.output_col = output_col

    def process(self, df: DataFrame) -> DataFrame:
        """
        Process DataFrame by adding a length column.

        Args:
            df: Input DataFrame

        Returns:
            DataFrame: DataFrame with added length column
        """
        return df.withColumn(self.output_col, length(col(self.input_col)))

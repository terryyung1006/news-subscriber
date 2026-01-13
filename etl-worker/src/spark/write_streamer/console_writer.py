from typing import List

from loguru import logger
from pyspark.sql import DataFrame, Row
from pyspark.sql.streaming import StreamingQuery

from .base import SparkWriteStreamer


class ConsoleWriter(SparkWriteStreamer):
    def write(self, df: DataFrame) -> StreamingQuery:
        query = (
            df.writeStream.foreachBatch(self.process_batch)
            .format("console")
            .outputMode("append")
            .trigger(processingTime="5 seconds")
            .option("checkpointLocation", "/tmp/spark-checkpoint")
            .start()
        )

        return query

    def process_batch(self, batch_df: DataFrame, batch_id: int) -> None:
        print(f"Processing batch ID: {batch_id}")
        batch_df.show(truncate=False)

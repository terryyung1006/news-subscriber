from typing import List

from loguru import logger
from pyspark.sql import DataFrame, Row
from pyspark.sql.streaming import StreamingQuery

from .base import SparkWriteToStorageStreamer


class ChromaWriterStreamer(SparkWriteToStorageStreamer):
    def __init__(self, storer):
        super().__init__(storer)

    def write(self, df: DataFrame) -> StreamingQuery:
        return df.writeStream.foreachBatch(self.process).outputMode("append").start()

    def process(self, batch_df: DataFrame, batch_id: int) -> None:
        logger.info(f"Processing batch ID: {batch_id}")
        rows = batch_df.collect()
        logger.info(f"Collected {len(rows)} rows from batch {batch_id}")

        processed_data = self.embed_rows(rows)
        logger.info(f"Processed {len(processed_data)} valid rows from batch {batch_id}")

        if processed_data:  # Only store if we have valid data
            logger.info(
                f"Storing {len(processed_data)} documents to ChromaDB for batch {batch_id}"
            )
            self.storer.store(processed_data)
            logger.info(f"Successfully completed ChromaDB write for batch {batch_id}")
        else:
            logger.warning(f"No valid data to store in batch {batch_id}")

    def embed_rows(self, rows: List[Row]) -> List[any]:
        processed_data = []
        skipped_count = 0

        for row in rows:
            try:
                # Validate that we have all required fields and embedding is not empty
                if (
                    row.document
                    and row.document.strip()
                    and hasattr(row, "embedding")
                    and row.embedding
                    and len(row.embedding) > 0
                    and hasattr(row, "chunk_id")
                    and row.chunk_id
                ):
                    # Ensure embedding is a list of numbers
                    embedding = row.embedding
                    if isinstance(embedding, list) and all(
                        isinstance(x, (int, float)) for x in embedding
                    ):
                        processed_data.append(
                            {
                                "id": row.chunk_id,
                                "embedding": embedding,
                                "metadata": {
                                    "category": row.metadata.category,
                                    "timestamp": str(row.metadata.timestamp),
                                    "length": str(row.metadata.length),
                                    "original_id": row.metadata.original_id,
                                },
                                "document": row.document,
                            }
                        )
                    else:
                        logger.warning(
                            f"Skipping row {row.chunk_id}: Invalid embedding format"
                        )
                        skipped_count += 1
                else:
                    logger.warning(
                        f"Skipping row: Missing required fields or empty embedding"
                    )
                    skipped_count += 1
            except Exception as e:
                logger.error(f"Error processing row: {e}")
                skipped_count += 1
                continue

        if skipped_count > 0:
            logger.warning(
                f"Skipped {skipped_count} invalid rows out of {len(rows)} total rows"
            )

        return processed_data

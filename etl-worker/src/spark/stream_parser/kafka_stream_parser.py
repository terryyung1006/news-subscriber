from pyspark.sql import DataFrame, SparkSession
from pyspark.sql.functions import col, from_json
from pyspark.sql.types import StructType

from .base import StreamParser


class KafkaStreamParser(StreamParser):
    """Kafka stream parser implementation using Spark Structured Streaming."""

    def __init__(
        self, spark: SparkSession, kafka_bootstrap_servers: str, schema: StructType
    ):
        """
        Initialize the Kafka stream parser.

        Args:
            spark: SparkSession instance
            kafka_bootstrap_servers: Comma-separated list of Kafka bootstrap servers
            schema: JSON schema for parsing the Kafka message values
        """
        self.spark = spark
        self.kafka_bootstrap_servers = kafka_bootstrap_servers
        self.schema = schema

    def process_batch(self, batch_df, batch_id):
        """
        This function processes each micro-batch of the streaming data.
        It prints a custom message for each pull attempt.
        """
        # Print a message for each batch processing attempt
        print(f"--- Pulling new messages from Kafka for batch ID: {batch_id} ---")

        # Show the content of the batch if there are messages
        if batch_df.count() > 0:
            print("Messages received in this batch:")
            batch_df.show(truncate=False)
        else:
            print("No new messages found in this batch.")

    def start_streaming(self, topic: str = "text-data") -> DataFrame:
        """Start the streaming application."""
        print(f"Starting Spark streaming from Kafka topic: {topic}")

        # Read from Kafka with better offset management
        streaming_df = (
            self.spark.readStream.format("kafka")
            .option("kafka.bootstrap.servers", self.kafka_bootstrap_servers)
            .option("subscribe", topic)
            .option("startingOffsets", "latest")
            .option("failOnDataLoss", "false")
            .load()
        )

        # Parse JSON from Kafka value
        parsed_df = (
            streaming_df.selectExpr("CAST(key AS STRING)", "CAST(value AS STRING)")
            .select(
                col("key").alias("kafka_key"),
                from_json(col("value"), self.schema).alias("data"),
            )
            .select("kafka_key", "data.*")
        )

        # Start streaming query with better configuration
        # parsed_df \
        #     .writeStream \
        #     .foreachBatch(self.process_batch) \
        #     .format("console") \
        #     .outputMode("append") \
        #     .trigger(processingTime='5 seconds') \
        #     .option("checkpointLocation", "/tmp/spark-checkpoint") \
        #     .start()

        # # Debug: Show the parsed DataFrame content
        # print("=== PARSED DATAFRAME DEBUG INFO ===")
        # print(f"Schema: {parsed_df.schema}")

        # # Try to show some data (be careful with streaming)
        # try:
        #     print("=== FIRST FEW ROWS ===")
        #     # For streaming, we need to be careful about collecting data
        #     # This will only work if there's data available
        #     sample_df = parsed_df.limit(3)
        #     sample_df.show(truncate=False)
        # except Exception as e:
        #     print(f"Could not show sample data (this is normal for empty streams): {e}")

        # print("=== END DEBUG INFO ===")

        return parsed_df

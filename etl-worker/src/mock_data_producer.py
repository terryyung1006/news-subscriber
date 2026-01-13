#!/usr/bin/env python3
"""
Mock data producer that uses the existing MockFinancialNewsSource.

This script fetches mock news from the mock source and sends it to Kafka
in the format expected by the streaming pipeline.
"""

import os
import time

from loguru import logger

from data_sources.factory import DataSourceFactory
from kafka_client.producer import KafkaProducer
from load_config import get_config


def main():
    """Main function to produce mock data using the existing mock source."""
    try:
        config = get_config()
        # Kafka configuration from config
        bootstrap_servers = config.get_string("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092")
        topic = "text-data"  # This could be made configurable too if needed

        # Initialize Kafka producer
        producer = KafkaProducer(bootstrap_servers=bootstrap_servers, topic=topic)

        logger.info(f"Connected to Kafka at {bootstrap_servers}")
        logger.info(f"Producing mock data to topic: {topic}")

        # Connect to Kafka
        if not producer.connect():
            logger.error("Failed to connect to Kafka")
            return 1

        # Produce mock data
        count = 0
        try:
            while True:
                # Create a new mock source instance each time to get fresh random data
                # Enable LLM generation for dynamic content
                factory = DataSourceFactory()
                model_name = config.get_string("OLLAMA_DEFAULT_MODEL", "llama2")
                print(f"Using model name: {model_name}")
                mock_source = factory.create_source(
                    "mock", news_count=5, use_llm=True, model_name=model_name
                )

                logger.info(f"Using mock source: {mock_source.source_name}")

                # Connect to mock source
                if not mock_source.connect():
                    logger.error("Failed to connect to mock source")
                    continue

                # Fetch news from mock source
                news_items = mock_source.fetch_news(limit=3)  # Fetch 3 items at a time

                for news_item in news_items:
                    # Convert to the format expected by Spark schema
                    streaming_message = {
                        "id": news_item.id,
                        "text": f"{news_item.title}\n\n{news_item.content}",
                        "timestamp": news_item.published_date.timestamp(),
                        "category": news_item.category,
                        "length": len(news_item.content),
                    }

                    # Send to Kafka using the raw producer to send custom format
                    success = producer._producer.send(
                        topic=producer.topic,
                        key=news_item.id,
                        value=streaming_message,
                    ).get(timeout=10)

                    if success:
                        count += 1
                        logger.info(
                            f"Sent news item {count}: {news_item.id} - "
                            f"{news_item.title[:50]}..."
                        )
                    else:
                        logger.error(f"Failed to send news item: {news_item.id}")

                # Disconnect from mock source to clean up
                mock_source.disconnect()

                # Wait before sending next batch
                logger.info("Waiting 10 seconds before next batch...")
                time.sleep(10)

        except KeyboardInterrupt:
            logger.info("Stopping mock data producer...")

    except Exception as e:
        logger.error(f"Error: {e}")
        logger.error("Make sure Kafka is running: make docker-up")
        return 1

    finally:
        if "producer" in locals():
            producer.disconnect()
        logger.info("Cleanup completed")

    return 0


if __name__ == "__main__":
    exit(main())

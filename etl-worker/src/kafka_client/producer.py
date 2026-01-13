"""
Kafka producer for sending financial news data.

This module provides a Kafka producer implementation for streaming
financial news data to Kafka topics.
"""

import json
from typing import Any, Dict, Optional

from kafka import KafkaProducer as KafkaClient
from loguru import logger

from data_sources.base import NewsItem


class KafkaProducer:
    """
    Kafka producer for sending financial news data.

    This class follows the Single Responsibility Principle by focusing
    solely on Kafka message production.
    """

    def __init__(self, bootstrap_servers: str, topic: str):
        """
        Initialize Kafka producer.

        Args:
            bootstrap_servers: Kafka bootstrap servers
            topic: Kafka topic to produce messages to
        """
        self.bootstrap_servers = bootstrap_servers
        self.topic = topic
        self._producer: Optional[KafkaClient] = None
        self._connected = False

    def connect(self) -> bool:
        """Establish connection to Kafka."""
        try:
            logger.info(f"Connecting to Kafka at {self.bootstrap_servers}")
            self._producer = KafkaClient(
                bootstrap_servers=self.bootstrap_servers,
                value_serializer=lambda v: json.dumps(v, default=str).encode("utf-8"),
                key_serializer=lambda k: k.encode("utf-8") if k else None,
            )
            self._connected = True
            logger.info("Successfully connected to Kafka")
            return True
        except Exception as e:
            logger.error(f"Failed to connect to Kafka: {e}")
            self._connected = False
            return False

    def disconnect(self) -> None:
        """Close connection to Kafka."""
        if self._producer:
            self._producer.close()
            self._producer = None
        self._connected = False
        logger.info("Disconnected from Kafka")

    def send_news_item(self, news_item: NewsItem, key: Optional[str] = None) -> bool:
        """
        Send a news item to Kafka.

        Args:
            news_item: News item to send
            key: Optional message key

        Returns:
            bool: True if message sent successfully, False otherwise
        """
        if not self._connected:
            logger.error("Not connected to Kafka")
            return False

        try:
            # Convert news item to dictionary
            message = {
                "id": news_item.id,
                "title": news_item.title,
                "content": news_item.content,
                "source": news_item.source,
                "published_date": news_item.published_date.isoformat(),
                "category": news_item.category,
                "sentiment": news_item.sentiment,
                "url": news_item.url,
                "metadata": news_item.metadata,
            }

            # Send message
            future = self._producer.send(
                topic=self.topic, key=key or news_item.id, value=message
            )

            # Wait for send to complete
            record_metadata = future.get(timeout=10)

            logger.info(
                f"Sent news item {news_item.id} to topic {self.topic} "
                f"partition {record_metadata.partition} offset {record_metadata.offset}"
            )
            return True
        except Exception as e:
            logger.error(f"Failed to send news item {news_item.id}: {e}")
            return False

    def send_processed_news_item(
        self, processed_item: Dict[str, Any], key: Optional[str] = None
    ) -> bool:
        """
        Send a processed news item to Kafka.

        Args:
            processed_item: Processed news item to send
            key: Optional message key

        Returns:
            bool: True if message sent successfully, False otherwise
        """
        if not self._connected:
            logger.error("Not connected to Kafka")
            return False

        try:
            # Convert processed item to dictionary
            message = {
                "original_id": processed_item["original_id"],
                "title": processed_item["title"],
                "content": processed_item["content"],
                "embedding": processed_item["embedding"],
                "sentiment_score": processed_item["sentiment_score"],
                "sentiment_label": processed_item["sentiment_label"],
                "keywords": processed_item["keywords"],
                "summary": processed_item["summary"],
                "metadata": processed_item["metadata"],
            }

            # Send message
            future = self._producer.send(
                topic=self.topic,
                key=key or processed_item["original_id"],
                value=message,
            )

            # Wait for send to complete
            record_metadata = future.get(timeout=10)

            logger.info(
                f"Sent processed news item {processed_item['original_id']} to topic {self.topic} "
                f"partition {record_metadata.partition} offset {record_metadata.offset}"
            )
            return True
        except Exception as e:
            logger.error(
                f"Failed to send processed news item {processed_item['original_id']}: {e}"
            )
            return False

    def send_batch(self, news_items: list[NewsItem]) -> int:
        """
        Send a batch of news items to Kafka.

        Args:
            news_items: List of news items to send

        Returns:
            int: Number of successfully sent messages
        """
        if not self._connected:
            logger.error("Not connected to Kafka")
            return 0

        success_count = 0
        for news_item in news_items:
            if self.send_news_item(news_item):
                success_count += 1

        logger.info(f"Sent {success_count}/{len(news_items)} news items to Kafka")
        return success_count

    def flush(self) -> None:
        """Flush any pending messages."""
        if self._producer:
            self._producer.flush()
            logger.info("Flushed pending Kafka messages")

    def is_connected(self) -> bool:
        """Check if connected to Kafka."""
        return self._connected

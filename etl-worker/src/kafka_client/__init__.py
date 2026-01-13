"""
Kafka integration package for RAG Financial News Worker.

This package contains Kafka producers and consumers for handling
financial news data streaming.
"""

from .producer import KafkaProducer

__all__ = ["KafkaProducer"]

"""
Vector storage package for RAG Financial News Worker.

This package contains interfaces and implementations for vector database
operations, primarily using Chroma for similarity search.
"""

from .base import SparkProcessedDataStorer
from .chroma_storage import ChromaVectorStorage

__all__ = ["SparkProcessedDataStorer", "ChromaVectorStorage"]

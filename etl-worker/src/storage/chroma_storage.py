"""
Chroma vector storage implementation.

This module provides a concrete implementation of the SparkProcessedDataStorer interface
using ChromaDB for vector similarity search.
"""

from typing import Any, Dict, List, Optional

import chromadb
from loguru import logger

from .base import SparkProcessedDataStorer

instance = None


class ChromaVectorStorage(SparkProcessedDataStorer):
    """
    ChromaDB implementation of vector storage.

    This class follows the Liskov Substitution Principle by properly
    implementing the SparkProcessedDataStorer interface.
    """

    def __init__(
        self,
        host: str = "localhost",
        port: int = 8000,
        collection_name: str = "financial_news",
    ):
        # Initialize ChromaDB client with v2 API configuration
        try:
            self.chroma_client = chromadb.HttpClient(host=host, port=port)

            # Test connection
            self.chroma_client.heartbeat()
            print("✅ Connected to ChromaDB successfully")

        except Exception as e:
            print(f"❌ Failed to connect to ChromaDB: {e}")
            print("Make sure ChromaDB is running: docker-compose up -d")
            raise

        # Create or get collection
        try:
            self.collection = self.chroma_client.get_collection("spark_documents")
            print("Using existing ChromaDB collection: spark_documents")
        except:
            try:
                self.collection = self.chroma_client.create_collection(
                    "spark_documents"
                )
                print("Created new ChromaDB collection: spark_documents")
            except Exception as e:
                print(f"Error creating collection: {e}")
                raise

    def store(self, data: List[Any]) -> None:
        try:
            if not data:
                print("No data to store")
                return

            # Filter out any items with empty embeddings
            valid_data = []
            for item in data:
                if (
                    item.get("embedding")
                    and isinstance(item["embedding"], list)
                    and len(item["embedding"]) > 0
                    and all(isinstance(x, (int, float)) for x in item["embedding"])
                ):
                    valid_data.append(item)
                else:
                    print(
                        f"Skipping item {item.get('id', 'unknown')}: Invalid or empty embedding"
                    )

            if not valid_data:
                print("No valid data to store after filtering")
                return

            # Prepare data for ChromaDB
            ids = [item["id"] for item in valid_data]
            embeddings = [item["embedding"] for item in valid_data]
            metadatas = [item["metadata"] for item in valid_data]
            documents = [item["document"] for item in valid_data]

            # Add to ChromaDB
            logger.info(
                f"Adding {len(ids)} documents to ChromaDB collection 'spark_documents'"
            )
            self.collection.add(
                ids=ids, embeddings=embeddings, metadatas=metadatas, documents=documents
            )
            logger.info(f"✅ Successfully added {len(ids)} documents to ChromaDB")
            logger.info(
                f"Document IDs: {ids[:5]}{'...' if len(ids) > 5 else ''}"
            )  # Show first 5 IDs
        except Exception as e:
            logger.error(f"❌ Error adding to ChromaDB: {e}")
            import traceback

            traceback.print_exc()

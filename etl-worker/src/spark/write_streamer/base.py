from abc import ABC, abstractmethod
from typing import Protocol

from pyspark.sql import DataFrame
from pyspark.sql.streaming import StreamingQuery

from storage.base import SparkProcessedDataStorer


class SparkWriteStreamer(Protocol):
    def write(self, df: DataFrame) -> StreamingQuery:
        pass


class SparkWriteToStorageStreamer(ABC, SparkWriteStreamer):
    def __init__(self, storer: SparkProcessedDataStorer):
        self.storer = storer

    @abstractmethod
    def write(self, df: DataFrame) -> StreamingQuery:
        pass

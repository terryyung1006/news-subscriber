from typing import List, Protocol


class SparkProcessedDataStorer(Protocol):
    def store(self, data: List[any]) -> None:
        pass

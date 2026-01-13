# Use Python 3.11 slim image
FROM python:3.11-slim

# Set working directory
WORKDIR /app

# System deps
RUN apt-get update && apt-get install -y \
    curl \
    openjdk-17-jre-headless \
    && rm -rf /var/lib/apt/lists/*

ENV JAVA_HOME=/usr/lib/jvm/java-17-openjdk-amd64
ENV PATH="$JAVA_HOME/bin:$PATH"

# Install uv (ultra-fast Python package manager)
RUN curl -LsSf https://astral.sh/uv/install.sh | sh && \
    ln -s /root/.local/bin/uv /usr/local/bin/uv

# Copy pyproject and source
COPY pyproject.toml README.md ./
COPY src/ ./src/
COPY config/ ./config/

# Install project with uv (prod dependencies only)
RUN uv venv /app/.venv && \
    . /app/.venv/bin/activate && \
    uv pip install .

ENV PATH="/app/.venv/bin:$PATH"
ENV PYTHONUNBUFFERED=1
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONPATH=/app

# Create non-root user
RUN useradd --create-home --shell /bin/bash app && chown -R app:app /app
USER app

# Default command runs the consumer
CMD ["python", "src/main.py"]

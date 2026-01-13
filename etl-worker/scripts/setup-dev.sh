#!/bin/bash

# Modern Development Setup using uv
set -e

echo "🚀 Setting up dev environment with uv"

# Load config as environment variables
echo "Loading configuration..."
eval "$(python3 src/load_config.py)"

# Ensure uv is installed
if ! command -v uv &> /dev/null; then
  echo "⬇️ Installing uv..."
  curl -LsSf https://astral.sh/uv/install.sh | sh
  export PATH="$HOME/.local/bin:$PATH"
fi

# Create and activate project env
uv venv .venv
source .venv/bin/activate

# Install project with dev extras
uv pip install -e ".[dev,test,docs]"

# Install pre-commit hooks if available
if command -v pre-commit &> /dev/null; then
  pre-commit install
fi

echo "✅ Dev environment ready. Activate with: source .venv/bin/activate"

#!/bin/bash

# Setup script for Ollama LLM integration

echo "Setting up Ollama for RAG Worker..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "Error: Docker is not running. Please start Docker first."
    exit 1
fi

# Start Ollama service
echo "Starting Ollama service..."
docker-compose -f docker-compose.ollama.yml up -d

# Wait for Ollama to be ready
echo "Waiting for Ollama to be ready..."
sleep 10

# Pull the default model from config
echo "Pulling ${OLLAMA_DEFAULT_MODEL:-llama2} model (this may take a while)..."
curl -X POST http://localhost:11434/api/pull -d "{\"name\": \"${OLLAMA_DEFAULT_MODEL:-llama2}\"}"

# Verify the model is available
echo "Verifying ${OLLAMA_DEFAULT_MODEL:-llama2} model availability..."
curl -X POST http://localhost:11434/api/generate -d "{\"model\": \"${OLLAMA_DEFAULT_MODEL:-llama2}\", \"prompt\": \"Hello\", \"stream\": false}" > /dev/null 2>&1

if [ $? -eq 0 ]; then
    echo "✅ Ollama setup completed successfully!"
    echo "You can now run the mock data producer with LLM generation."
else
    echo "❌ Ollama setup failed. Please check the logs:"
    echo "docker-compose -f docker-compose.ollama.yml logs"
fi

docker-compose up -d

export PYSPARK_PYTHON=$(which python3.11)
export PYSPARK_DRIVER_PYTHON=$(which python3.11)

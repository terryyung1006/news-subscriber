# Modern Makefile for RAG Financial News Worker (uv-based)

.PHONY: help install install-dev install-test install-docs clean test test-cov lint format type-check build docs serve-docs docker-build docker-up docker-down docker-logs setup-dev run mock-data

help: ## Show this help message
	@echo "RAG Financial News Worker - Commands"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

setup-dev: ## Set up dev env with uv
	./scripts/setup-dev.sh

install: ## Install package (prod)
	uv pip install .

install-dev: ## Install with dev extras
	uv pip install -e ".[dev,test,docs]"

install-test: ## Install with test extras only
	uv pip install -e ".[test]"

install-docs: ## Install with docs extras only
	uv pip install -e ".[docs]"

run: ## Run the main RAG worker application
	cd src && uv run python main.py

mock-data: ## Start mock data producer using existing mock source
	cd src && uv run python mock_data_producer.py

lint: ## Lint code
	uv run ruff check src/ tests/
	uv run black --check src/ tests/
	uv run isort --check-only src/ tests/

format: ## Format code
	uv run black src/ tests/
	uv run isort src/ tests/
	uv run ruff check --fix src/ tests/

type-check: ## Type checking
	uv run mypy src/

test: ## Run tests
	uv run pytest tests/ -v

test-cov: ## Tests with coverage
	uv run pytest tests/ --cov=src --cov-report=term-missing --cov-report=html

build: ## Build wheel/sdist
	uv run hatch build

docs: ## Build docs
	uv run mkdocs build

serve-docs: ## Serve docs locally
	uv run mkdocs serve

docker-build: ## Build Docker image
	docker build -f docker/Dockerfile -t rag-worker .

docker-up: ## Start services
	docker-compose up -d

docker-down: ## Stop services
	docker-compose down && docker-compose -f docker-compose.ollama.yml down

docker-down-flush: ## Stop services and flush data
	docker-compose down -v && docker-compose -f docker-compose.ollama.yml down -v

docker-logs: ## Tail logs
	docker-compose logs -f rag-worker

clean: ## Clean artifacts
	rm -rf dist/ build/ *.egg-info/ .pytest_cache/ .mypy_cache/ htmlcov/ __pycache__
	find . -type d -name __pycache__ -delete
	find . -type f -name "*.pyc" -delete

.DEFAULT_GOAL := help

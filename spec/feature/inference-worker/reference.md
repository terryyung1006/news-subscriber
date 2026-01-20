# Quick Reference: Inference Worker

## Commands
-   **Start Infrastructure**: `cd backend && docker-compose up -d`
-   **Start Worker**: `cd inference-worker && pip3 install -r requirements.txt && export PYTHONPATH=$PYTHONPATH:. && python3 src/worker.py`
-   **Start Backend**: `cd backend && go run src/main.go`
-   **Start Frontend**: `cd frontend && npm run dev`

## Key Files
-   `inference-worker/src/worker.py`: Main entry point, listens to Redis.
-   `inference-worker/src/tasks.py`: LLM Logic (LangChain).
-   `backend/src/service/chat.go`: Enqueues jobs to Redis.
-   `frontend/src/app/api/chat/send/route.ts`: Next.js API Proxy.

## Configuration
-   **Redis**: `localhost:6379`
-   **Ollama**: `localhost:11434` (Model: `llama3`)
-   **Queue Name**: `inference_queue`

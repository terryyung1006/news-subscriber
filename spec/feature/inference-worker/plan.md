# Inference Worker Implementation Plan

## Objective
Create a scalable inference worker system to handle LLM tasks (chat, summarization) asynchronously, decoupling heavy processing from the main backend.

## Architecture
- **Frontend**: Next.js (Chat Interface)
- **Backend**: Go (API Gateway, Enqueue Jobs)
- **Broker**: Redis (Message Queue)
- **Worker**: Python (LangChain + Ollama/OpenAI)

## Tasks
1.  [x] **Infrastructure**: Add Redis and Ollama to `docker-compose.yaml`.
2.  [x] **Worker Service**:
    -   Create `inference-worker` directory.
    -   Implement `worker.py` to listen to Redis `inference_queue`.
    -   Implement `tasks.py` to handle LLM logic using LangChain.
    -   Dockerize the worker.
3.  [x] **Backend Integration**:
    -   Add Redis client to Go backend.
    -   Implement `POST /api/chat` (HTTP) to enqueue jobs.
    -   Implement polling logic to wait for results.
4.  [x] **Frontend Integration**:
    -   Create API Route `api/chat/send` to proxy requests.
    -   Update `chat-interface.tsx` to use real API and show loading state.

## Future Improvements
-   Implement WebSockets for real-time streaming tokens.
-   Add persistent storage for chat history (Postgres).
-   Scale workers using Kubernetes/KEDA.

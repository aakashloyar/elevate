# problem_generation

HTTP, persistence, and worker service for AI generation jobs.

Endpoints:

- `POST /generation-jobs`
- `GET /generation-jobs/{jobId}`

Environment:

- `POSTGRES_HOST`
- `POSTGRES_PORT`
- `POSTGRES_USER`
- `POSTGRES_PASSWORD`
- `POSTGRES_DB`
- `POSTGRES_SSLMODE`
- `HTTP_PORT`
- `KAFKA_BROKERS`
- `KAFKA_CLIENT_ID`
- `KAFKA_GROUP_ID`
- `KAFKA_API_KEY`
- `KAFKA_API_SECRET`
- `GENERATION_WORKER_ENABLED`
- `AI_PROVIDER`
- `AI_MODEL`
- `AI_BASE_URL`
- `GEMINI_API_KEY`

`Elevate_Gemini_API_Key` is also supported for the Gemini API key to match existing local config.

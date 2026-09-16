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
- `GROQ_API_KEY` (or `GROK_API_KEY` for existing local configuration)
- `AI_PROVIDER` (`gemini`/`google` or `groq`/`grok`)

`Elevate_Gemini_API_Key` is also supported for the Gemini API key to match existing local config.

To use Groq, set `AI_PROVIDER=groq`, provide `GROQ_API_KEY` (or the existing `GROK_API_KEY` name), and optionally set `AI_MODEL` (for example, `openai/gpt-oss-120b`). Leave `AI_BASE_URL` empty to use Groq's default endpoint, or set it to a compatible OpenAI-style endpoint.

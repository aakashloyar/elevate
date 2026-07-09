# problem_generation

HTTP and persistence service for generation jobs.

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
- `KAFKA_GENERATION_REQUESTS_TOPIC`

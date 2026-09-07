# Problem generation architecture

The HTTP side creates a generation job, stores it as `pending`, and publishes a generation-requested event to Kafka.

The worker runs inside the same service when `GENERATION_WORKER_ENABLED=true`.

Worker flow:

- consume the generation-requested Kafka event
- mark the job as `processing`
- call the configured AI provider through the generic `AITextGenerator` port
- parse and validate the generated problem JSON
- publish the generated problem batch to the problem service Kafka topic
- mark the job as `completed`
- mark the job as `failed` if AI generation, parsing, validation, or publishing fails

Gemini is currently implemented as an outbound AI adapter. Other providers can be added by implementing the same `AITextGenerator` port and selecting them from config.

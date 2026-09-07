1. POST /generation-jobs
-> creates a problem generation job.

Request body:

```json
{
  "user_id": "U123",
  "single_correct_count": 5,
  "multi_correct_count": 3,
  "numerical_count": 2,
  "document_id": "D123",
  "assessment_id": "A123",
  "level": "medium",
  "description": "Generate questions from algebra basics",
  "topic_ids": ["T1", "T2"]
}
```

Allowed `level` values:

```text
easy
medium
hard
```

2. GET /generation-jobs/{jobId}
-> get status of specific jobId

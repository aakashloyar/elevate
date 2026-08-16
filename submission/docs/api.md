1. Post /submission
-> on this endpoint send assessment_id, user_id, and duration_seconds
-> we will get submissionid and starttime in reponse

2. Post /submission/id
-> here we are posting our responses to submission service 

3. Get /submissions/{submissionId}/status
-> returns the submission status and expiry timestamp for frontend polling.

Response

```json
{
  "submission_id": "S123",
  "status": "IN_PROGRESS",
  "expires_at": "2026-08-16T22:30:00Z"
}
```

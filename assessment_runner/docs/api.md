1. Start Exam
POST /sessions

Request

{
    "assessmentId": "A123"
}

Response

{
    "sessionId": "...",
    "submissionId": "...",
    "remainingTime": 7200,
    "totalQuestions": 100
}

Responsibilities:

Verify assessment exists
Verify user has access
Verify exam can be started
Create submission
Return session information


2. Get Session
GET /sessions/{sessionId}

Returns

Remaining time
Total questions
Current status

3. Get Questions

Instead of fetching one question at a time, I'd support pagination.

GET /sessions/{sessionId}/questions

Query parameters

offset=0
limit=10

Response

{
    "questions": [...]
}

Internally, the Exam Runner:

Gets the problem IDs from the Assessment Service.
Fetches the corresponding problems from the Problem Service.
Returns them to the frontend.

4. Submit Exam
POST /sessions/{sessionId}/submit

Responsibilities:

Mark submission as complete.
Publish the event for evaluation.
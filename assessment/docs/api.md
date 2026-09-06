1. Post /assessment
Description -> This api is just to create assessment

2. Get /assessment/id
Description -> This api is just sending the meta data of the specific assessment 

3. Delete /assesment/id
Description -> This is just to delete a specific assesment 

4. GET /assessments/{assessmentId}/marking-scheme
Description -> Returns the marking scheme configured for an assessment.

5. POST /assessments/{assessmentId}/marking-scheme
Description -> Creates a marking scheme. Returns `409 Conflict` if the assessment already has one.

6. PUT /assessments/{assessmentId}/marking-scheme
Description -> Creates or replaces the complete marking scheme for an assessment.

Request body:
```json
{
  "single": {
    "correct": 4,
    "incorrect": -1,
    "skipped": 0
  },
  "multiple": {
    "correct": 4,
    "incorrect": -1,
    "skipped": 0
  },
  "numerical": {
    "correct": 4,
    "incorrect": 0,
    "skipped": 0
  }
}
```

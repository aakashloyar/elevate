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
  "single_correct_marks": 4,
  "single_incorrect_marks": -1,
  "single_skipped_marks": 0,
  "multiple_correct_marks": 4,
  "multiple_incorrect_marks": -1,
  "multiple_skipped_marks": 0,
  "numerical_correct_marks": 4,
  "numerical_incorrect_marks": 0,
  "numerical_skipped_marks": 0
}
```


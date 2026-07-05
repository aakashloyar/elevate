Database Design

1. Problem
Problem
-------
id (UUID)

createdBy

title (optional)

statement

type
// SINGLE_CORRECT
// MULTIPLE_CORRECT
// NUMERICAL
// CODING (future)
// ESSAY (future)

difficulty
// EASY
// MEDIUM
// HARD

sourceType
// MANUAL
// AI

status
// DRAFT
// ACTIVE
// ARCHIVED

createdAt
updatedAt


2. ProblemOption

ProblemOption
-------------
id

problemId

text

isCorrect


3. ProblemTag

ProblemTag
----------
problemId

tag
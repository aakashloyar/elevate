# database 

Table1:- Assessment
1.  id
2.  admin -> array of userid
3.  acess -> array of userid
    topics -> array of string
    maximumtime -> number in seconds
4.  totalquestion -> number
5.  totalsingle correct -> number
6.  totalmultiple correct -> number
7.  total numerical -> number
8.  singlecorrectpositivemarks -> number
9.  singlecorrectnegativemarks -> number
10. singlecorrectskipmarks -> number
11. multiplecorrectpositivemarks -> number
12. multiplecorrectnegativemarks -> number
13. multiplecorrectskipmarks -> number
14. numericalcorrectpositivemarks -> number
15. numericalcorrectnegativeemarks -> number
16. numericalcorrectskipmarks -> number
17. createdAt
18. updatedAt

1. Assessment

```
Assessment
----------
id

title
description

status
// DRAFT
// PUBLISHED
// RUNNING
// COMPLETED
// ARCHIVED

durationSeconds

startTime
endTime

problemCount

createdBy

createdAt
updatedAt
```

2. AssessmentAdmin

```
AssessmentAdmin
---------------
assessmentId

userId

createdAt
```

3. AssessmentAccess

```
AssessmentAccess
----------------
assessmentId

userId

createdAt

```

4. AssessmentTopic

```
AssessmentTopic
---------------
assessmentId

topic
```

5. AssessmentMarking

```
AssessmentMarking
-----------------

assessmentId

singleCorrectPositive
singleCorrectNegative
singleCorrectSkip

multipleCorrectPositive
multipleCorrectNegative
multipleCorrectSkip

numericalPositive
numericalNegative
numericalSkip
```

├── AssessmentAdmin
├── AssessmentAccess
├── AssessmentTopic
├── AssessmentMarking
why we are not creating id for each new entry is this fine

-> here if you will see that id is required for unique entry 
-> and everywhere there is not requirement of id 
-> as let us take example of Assessment Admin 
-> here assessmentid and userid are Primary key 
-> now these both must not match for any entry 


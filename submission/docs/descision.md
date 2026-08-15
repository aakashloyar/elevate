1. Draft
-> now working on submission service now currently 
-> i want to know what strategy to choose saving 
-> as draft in backend or frontend only 
-> backend saving is good

-> should we use kafka or http for draft saving
-> use http
-> as for frontend it does not matter 
-> for load
-> Kafka is for service-to-service asynchronous events
-> not for making the frontend communicate with your backend.

2. Unsumbitted exam time expired
ok now i want to understand ki suppose user submit then it is fine answer gets submitted but i want to understand suppose user exits the application and user not saved and then exam time finsihed now hot handle it ki exam is submitted 

2.1. Scheduled/timeout worker — primary mechanism

```
Submission
     │
     │ expires_at = 11:00
     ▼
Expiration Worker
     │
     │ 11:00
     ▼
Check submission
     │
     ▼
IN_PROGRESS?
     │
     ▼
Mark SUBMITTED
```

-> The worker updates:

```
status = SUBMITTED
submitted_at = 11:00
```

-> Then publishes:

-> SubmissionSubmitted to Kafka.

```
Submission Service
       │
       ▼
Kafka
       │
       ▼
Evaluation Service
```

-> So even if the user has completely closed the application, the submission is finalized.

2.2. Lazy expiration — important backup

-> suppose user donot submit
-> but open application after some time
-> it will check if expired if yes then 
-> proceed


-> what we will select we will proceed with option 1
-> as it is more better and instead of migrating 
-> now only we can proceed with it 
-> background worker you can start during start of main program
-> it can check status of all the submission 
-> will pick which one are in progress and expired
-> to make them submitted


Problem1:-
-> but suppose if background worker take some inprogress and suppose now those are already submitted

Solution:- 
-> The solution: don't blindly update what you fetched
-> Use a conditional update in the database.

```
UPDATE submissions
SET
    status = 'SUBMITTED',
    submitted_at = NOW()
WHERE id = $1
  AND status = 'IN_PROGRESS'
  AND expires_at <= NOW();
```
Worker
  │
  ├── Submission A → already SUBMITTED → 0 rows updated
  │
  ├── Submission B → IN_PROGRESS + expired → 1 row updated
  │
  └── Submission C → IN_PROGRESS + expired → 1 row updated

Problem2:- 
Multiple workers

-> Even better: claim rows atomically
-> If you have multiple worker instances, you can use PostgreSQL's FOR UPDATE SKIP LOCKED:

```
SELECT id
FROM submissions
WHERE status = 'IN_PROGRESS'
  AND expires_at <= NOW()
FOR UPDATE SKIP LOCKED
LIMIT 100;
```

-> Then the worker processes those rows inside a transaction.
-> This allows:
```
Worker 1 → Submission 1-100
Worker 2 → Submission 101-200
Worker 3 → Submission 201-300
```
-> without workers fighting over the same rows.

Problem3:-
-> But there is another important race

Imagine:
10:59:59
Worker finds expired submission

10:59:59
Student clicks Submit

11:00:00
Worker updates submission

-> You need to define your business rule around the deadline.
```
if current_time <= expires_at
    manual submission is allowed

if current_time > expires_at
    manual submission is rejected
```

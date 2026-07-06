1. Backend cache -> future


2. Heartbeat (Optional) -> future

If you later add proctoring or live monitoring:

POST /sessions/{sessionId}/heartbeat

This tells the backend the student is still connected.

3. Resume Session (Optional) -> future

If your platform allows resuming:

GET /sessions/{sessionId}/resume

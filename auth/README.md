# Auth service

The auth service owns one-time codes and issues short-lived HS256 access tokens. It does not store passwords, verification-state columns, or refresh tokens; every login requires a fresh OTP.

Required environment:

- `JWT_SECRET`: the same high-entropy secret must be configured in auth and every service that verifies tokens.
- `ACCESS_TOKEN_TTL_HOURS`: defaults to `3`.
- `USER_SERVICE_URL`: defaults to `http://localhost:8081`.
- Standard `POSTGRES_*` settings and `HTTP_PORT`.

Notification delivery uses Kafka with JSON event payloads. Configure `KAFKA_NOTIFICATION_EMAIL_TOPIC`, `KAFKA_BROKERS`, `KAFKA_CLIENT_ID`, `KAFKA_API_KEY`, and `KAFKA_API_SECRET`.

Routes:

- `POST /auth/register` with `username` and `email`; sends an OTP
- `POST /auth/login` with `identifier` (username or email); sends an OTP
- `POST /auth/verify-otp` with `user_id` and `otp`; verifies the email and returns an access token
- `GET /health`

The reusable verifier is in `auth/token`. Services can import it and call `token.Parse(JWT_SECRET, bearerToken, time.Now())` to validate the signature and expiry.

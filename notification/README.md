# Notification service

The notification service sends email through SMTP. It accepts the existing `POST /notifications/email` endpoint and also consumes asynchronous email events from Kafka.

Kafka configuration:

- `KAFKA_BROKERS`: comma-separated broker addresses
- `KAFKA_CLIENT_ID`: Kafka client ID
- `KAFKA_API_KEY` and `KAFKA_API_SECRET`: SASL credentials
- `KAFKA_NOTIFICATION_EMAIL_TOPIC`: defaults to `notification.email`
- `KAFKA_NOTIFICATION_GROUP_ID`: defaults to `notification-service`

The producer and consumer exchange JSON events. The consumer commits an event only after SMTP succeeds. Invalid events are committed to avoid blocking the partition, while temporary SMTP failures are left uncommitted and retried.

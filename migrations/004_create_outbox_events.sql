CREATE TABLE outbox_events (
    id CHAR(36) PRIMARY KEY,
    event_type VARCHAR(255) NOT NULL,
    aggregate_id CHAR(36) NOT NULL,
    payload JSON NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    published_at TIMESTAMP NULL DEFAULT NULL
);
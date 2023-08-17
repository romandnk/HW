CREATE TABLE events (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    date TIMESTAMP NOT NULL,
    duration INTERVAL HOUR TO SECOND NOT NULL,
    description TEXT,
    user_id INTEGER NOT NULL,
    notification_interval INTERVAL,
    scheduled boolean DEFAULT FALSE NOT NULL
);

CREATE INDEX idx_events_id ON events (id);
CREATE INDEX idx_events_date ON events (date);

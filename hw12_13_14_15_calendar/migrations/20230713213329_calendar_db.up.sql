CREATE TABLE events (
    id VARCHAR(32) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    date TIMESTAMP NOT NULL,
    duration INTERVAL NOT NULL,
    description TEXT,
    user_id UUID NOT NULL,
    notification_interval INTERVAL
);

CREATE INDEX idx_events_id ON events (id);
CREATE INDEX idx_events_date ON events (date);


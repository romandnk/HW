CREATE TABLE events (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    date TIMESTAMP NOT NULL,
    duration INTERVAL NOT NULL,
    description TEXT,
    user_id UUID NOT NULL,
    notification_interval INTERVAL
);

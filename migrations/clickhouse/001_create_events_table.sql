CREATE TABLE IF NOT EXISTS events (
    event_id UUID,
    event_type String,
    timestamp DateTime DEFAULT now()
    event_payload
) ENGINE = MergeTree()
ORDER BY (timestamp, event_type);

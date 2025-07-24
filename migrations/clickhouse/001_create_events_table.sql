CREATE TABLE IF NOT EXISTS events (
    event_id UUID,
    event_type String,
    timestamp DateTime DEFAULT now(),
    event_payload String
) ENGINE = MergeTree()
ORDER BY (timestamp, event_type);

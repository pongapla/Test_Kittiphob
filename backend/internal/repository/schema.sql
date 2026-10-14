CREATE TABLE IF NOT EXISTS queue_state (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    current_index INTEGER NOT NULL DEFAULT -1
        CHECK (current_index BETWEEN -1 AND 259),
    cycle INTEGER NOT NULL DEFAULT 1,
    updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS queue_tickets (
    id TEXT PRIMARY KEY,
    cycle INTEGER NOT NULL,
    queue_number TEXT NOT NULL,
    issued_at TEXT NOT NULL,
    UNIQUE(cycle, queue_number)
);

INSERT OR IGNORE INTO queue_state
    (id, current_index, cycle, updated_at)
VALUES
    (1, -1, 1, CURRENT_TIMESTAMP);

CREATE INDEX IF NOT EXISTS idx_queue_tickets_cycle
ON queue_tickets(cycle);

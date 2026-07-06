-- Core tables. Later milestones add their own migrations; this one is just
-- enough for accounts, course enrollment, and progress bookkeeping.

CREATE TABLE users (
    id            INTEGER PRIMARY KEY,
    username      TEXT NOT NULL UNIQUE,
    email         TEXT UNIQUE,
    password_hash TEXT,
    display_name  TEXT NOT NULL DEFAULT '',
    avatar        TEXT NOT NULL DEFAULT 'tama-1',
    daily_goal_xp INTEGER NOT NULL DEFAULT 20,
    timezone      TEXT NOT NULL DEFAULT 'UTC',
    is_admin      INTEGER NOT NULL DEFAULT 0,
    created_at    TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE sessions (
    token      TEXT PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    expires_at TEXT NOT NULL
);

-- One row per (user, course) enrollment. The course id is a pair like
-- "ja-from-en"; pack contents live on disk, not in the database.
CREATE TABLE enrollments (
    user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_id  TEXT NOT NULL,
    active     INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    PRIMARY KEY (user_id, course_id)
);

-- Path position and per-node completion.
CREATE TABLE progress (
    user_id      INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_id    TEXT NOT NULL,
    node_id      TEXT NOT NULL,
    lessons_done INTEGER NOT NULL DEFAULT 0,
    legendary    INTEGER NOT NULL DEFAULT 0,
    updated_at   TEXT NOT NULL DEFAULT (datetime('now')),
    PRIMARY KEY (user_id, course_id, node_id)
);

-- Append-only XP ledger; balances are sums, never stored counters.
CREATE TABLE xp_events (
    id         INTEGER PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount     INTEGER NOT NULL,
    source     TEXT NOT NULL,
    source_ref TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX xp_events_user_day ON xp_events (user_id, created_at);

-- Same idea for gems.
CREATE TABLE gem_events (
    id         INTEGER PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount     INTEGER NOT NULL,
    source     TEXT NOT NULL,
    source_ref TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

-- Hearts and streaks are small per-user state rows, updated in place.
CREATE TABLE hearts (
    user_id      INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    hearts       INTEGER NOT NULL DEFAULT 5,
    refill_start TEXT
);

CREATE TABLE streaks (
    user_id       INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    current       INTEGER NOT NULL DEFAULT 0,
    longest       INTEGER NOT NULL DEFAULT 0,
    last_day      TEXT,
    freezes       INTEGER NOT NULL DEFAULT 0
);

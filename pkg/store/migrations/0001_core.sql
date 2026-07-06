-- Core tables for accounts, courses, packs, progress, and gamification.
-- Every timestamp column is UTC unix milliseconds stored as INTEGER.

CREATE TABLE users (
    id            INTEGER PRIMARY KEY,
    username      TEXT NOT NULL UNIQUE COLLATE NOCASE,
    password_hash TEXT NOT NULL,
    display_name  TEXT NOT NULL DEFAULT '',
    avatar        TEXT NOT NULL DEFAULT '',
    is_admin      INTEGER NOT NULL DEFAULT 0,
    created_at    INTEGER NOT NULL,
    settings_json TEXT NOT NULL DEFAULT '{}'
);

-- Token is 32 random bytes base64url, generated in pkg/store.
CREATE TABLE sessions (
    token      TEXT PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at INTEGER NOT NULL,
    expires_at INTEGER NOT NULL,
    user_agent TEXT NOT NULL DEFAULT ''
);

CREATE INDEX sessions_user ON sessions (user_id);
CREATE INDEX sessions_expiry ON sessions (expires_at);

-- id is "<target>-from-<base>" like "es-from-en".
CREATE TABLE courses (
    id          TEXT PRIMARY KEY,
    base_lang   TEXT NOT NULL,
    target_lang TEXT NOT NULL,
    title       TEXT NOT NULL,
    pack_id     INTEGER REFERENCES course_packs(id),
    status      TEXT NOT NULL DEFAULT 'empty' CHECK (status IN ('empty', 'generating', 'ready', 'error')),
    created_at  INTEGER NOT NULL,
    UNIQUE (base_lang, target_lang)
);

-- content is the zstd-compressed pack JSON; generated_by records the model
-- name string only.
CREATE TABLE course_packs (
    id           INTEGER PRIMARY KEY,
    course_id    TEXT NOT NULL REFERENCES courses(id),
    version      INTEGER NOT NULL,
    format       INTEGER NOT NULL,
    content      BLOB NOT NULL,
    sha256       TEXT NOT NULL,
    generated_by TEXT NOT NULL DEFAULT '',
    created_at   INTEGER NOT NULL,
    UNIQUE (course_id, version)
);

CREATE INDEX course_packs_course ON course_packs (course_id, version);

CREATE TABLE progress (
    user_id      INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_id    TEXT NOT NULL REFERENCES courses(id),
    node_id      TEXT NOT NULL,
    crowns       INTEGER NOT NULL DEFAULT 0,
    legendary    INTEGER NOT NULL DEFAULT 0,
    completed_at INTEGER,
    updated_at   INTEGER NOT NULL,
    PRIMARY KEY (user_id, course_id, node_id)
);

-- Append-only; XP totals are always SUM over this ledger, never a counter.
CREATE TABLE xp_events (
    id         INTEGER PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_id  TEXT,
    amount     INTEGER NOT NULL,
    reason     TEXT NOT NULL CHECK (reason IN ('lesson', 'practice', 'story', 'quest', 'bonus', 'legendary')),
    created_at INTEGER NOT NULL
);

CREATE INDEX xp_events_user_time ON xp_events (user_id, created_at);

-- last_active_day is the user-local YYYY-MM-DD.
CREATE TABLE streaks (
    user_id         INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    current         INTEGER NOT NULL DEFAULT 0,
    longest         INTEGER NOT NULL DEFAULT 0,
    last_active_day TEXT NOT NULL DEFAULT '',
    freezes         INTEGER NOT NULL DEFAULT 0,
    updated_at      INTEGER NOT NULL
);

-- Refill math (one heart per 5 hours) lives in pkg/engine; this row stores
-- only the anchor timestamps.
CREATE TABLE hearts (
    user_id           INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    count             INTEGER NOT NULL DEFAULT 5,
    max               INTEGER NOT NULL DEFAULT 5,
    refill_started_at INTEGER,
    unlimited_until   INTEGER
);

-- Signed amounts, balance is SUM; overdrafts are rejected in the engine
-- inside the write transaction, not by a CHECK trigger.
CREATE TABLE gems_ledger (
    id         INTEGER PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount     INTEGER NOT NULL,
    reason     TEXT NOT NULL,
    ref        TEXT NOT NULL DEFAULT '',
    created_at INTEGER NOT NULL
);

CREATE INDEX gems_ledger_user ON gems_ledger (user_id);

CREATE TABLE review_items (
    user_id     INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_id   TEXT NOT NULL,
    item_id     TEXT NOT NULL,
    state       INTEGER NOT NULL DEFAULT 0,
    stability   REAL NOT NULL DEFAULT 0,
    difficulty  REAL NOT NULL DEFAULT 0,
    due         INTEGER NOT NULL,
    last_review INTEGER,
    lapses      INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, course_id, item_id)
);

CREATE INDEX review_items_due ON review_items (user_id, course_id, due);

CREATE TABLE players (
    id TEXT PRIMARY KEY,
    first_name TEXT NOT NULL CHECK(first_name != ''),
    last_name TEXT NOT NULL CHECK(last_name != ''),
    gender TEXT NOT NULL CHECK(gender != ''),
    deleted_at TEXT
);

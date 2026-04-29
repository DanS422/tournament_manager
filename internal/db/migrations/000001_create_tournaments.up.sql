CREATE TABLE tournaments (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL CHECK(name != ''),
    location TEXT NOT NULL CHECK(location != '')
);

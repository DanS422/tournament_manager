CREATE TABLE disciplines (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL CHECK(name != ''),
    no_of_team_players INTEGER NOT NULL DEFAULT 1, 
    tournament_id TEXT NOT NULL,

    FOREIGN KEY (tournament_id) REFERENCES tournaments(id) ON DELETE CASCADE
);

CREATE INDEX idx_disciplines_tournament_id ON disciplines(tournament_id);


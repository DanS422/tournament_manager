CREATE TABLE attendants (
    id TEXT PRIMARY KEY,
    tournament_id TEXT NOT NULL,
    player_id TEXT NOT NULL,

    FOREIGN KEY (tournament_id) REFERENCES tournaments(id) ON DELETE CASCADE,
    FOREIGN KEY (player_id) REFERENCES players(id),

    UNIQUE(tournament_id, player_id)
);

CREATE INDEX idx_attendants_tournament_id ON attendants(tournament_id);
CREATE INDEX idx_attendants_player_id ON attendants(player_id);
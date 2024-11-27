CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE songs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    group_id INTEGER NOT NULL,
    release_date VARCHAR(255),
    text TEXT,
    link VARCHAR(255),
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
);

CREATE INDEX idx_song_name ON songs (name);
CREATE INDEX idx_group_id ON songs (group_id);
CREATE INDEX idx_release_date ON songs (release_date);
CREATE INDEX idx_link ON songs (link);
CREATE INDEX idx_group_name ON groups (name)
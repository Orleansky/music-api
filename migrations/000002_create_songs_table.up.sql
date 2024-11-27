CREATE TABLE songs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    group_id INTEGER NOT NULL,
    release_date VARCHAR(255),
    text TEXT,
    link VARCHAR(255),
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
);
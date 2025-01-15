-- +goose Up
-- +goose StatementBegin
CREATE TABLE words (
    id INTEGER PRIMARY KEY,
    value TEXT NOT NULL UNIQUE,
    word TEXT NOT NULL
);

CREATE TABLE subjects (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    word1 INTEGER NOT NULL REFERENCES words(id),
    word2 INTEGER NOT NULL REFERENCES words(id),
    word3 INTEGER NOT NULL REFERENCES words(id),
    word4 INTEGER NOT NULL REFERENCES words(id),
    UNIQUE(word1, word2, word3, word4)
);

CREATE TABLE games (
    id INTEGER PRIMARY KEY,
    active BOOLEAN DEFAULT FALSE,
    turns INTEGER DEFAULT 0,
    wrong INTEGER DEFAULT 0,
    win BOOLEAN DEFAULT 0,
    start INTEGER,
    end INTEGER,
    bid INTEGER
);

CREATE TABLE users_games (
    uid INTEGER NOT NULL,
    gid INTEGER NOT NULL,
    PRIMARY KEY(uid, gid)
);

CREATE TABLE boards (
    id INTEGER PRIMARY KEY,
    subject1 INTEGER REFERENCES subjects(id),
    subject2 INTEGER REFERENCES subjects(id),
    subject3 INTEGER REFERENCES subjects(id),
    subject4 INTEGER REFERENCES subjects(id),
    played INTEGER DEFAULT 0,
    wins INTEGER DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE words;
DROP TABLE subjects;
DROP TABLE games;
DROP TABLE users_games;
DROP TABLE boards;
-- +goose StatementEnd

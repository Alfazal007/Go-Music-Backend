-- +goose Up
CREATE TABLE users (
    id uuid PRIMARY KEY,
    username text unique not null,
    password text not null,
    refresh_token text,
    profile_picture text,
    CHECK (LENGTH(username) BETWEEN 5 AND 30),
    CHECK (LENGTH(password) >= 6)
);
CREATE TABLE songs (
    id uuid PRIMARY KEY,
    name text unique not null,
    song_link text not null,
    user_id uuid not null references users(id) on delete cascade
    CHECK (LENGTH(name) BETWEEN 3 AND 20)
);
CREATE TABLE liked_songs (
    id uuid PRIMARY KEY,
    user_id uuid not null references users(id) on delete cascade,
    song_id uuid not null references songs(id) on delete cascade,
    unique(user_id, song_id)
);

-- +goose Down
DROP TABLE liked_songs;
DROP TABLE songs;
DROP TABLE users;


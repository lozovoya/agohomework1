CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    login TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    token TEXT,
    roles TEXT[] NOT NULL DEFAULT '{}',
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE cards (
                       id BIGSERIAL PRIMARY KEY,
                       number TEXT NOT NULL UNIQUE,
                       owner BIGINT NOT NULL REFERENCES users,
                       balance BIGINT NOT NULL DEFAULT 0
);

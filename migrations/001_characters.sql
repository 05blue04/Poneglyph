-- +goose Up
CREATE TABLE IF NOT EXISTS characters (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text UNIQUE NOT NULL,
    age integer NOT NULL,
    description TEXT NOT NULL,
    origin text NOT NULL,
    race text NOT NULL DEFAULT 'human',
    bounty bigint , -- change to current bounty at some point
    -- add previous bounties a [] bigint
    episode int NOT NULL
);

-- +goose Down
DROP TABLE characters;
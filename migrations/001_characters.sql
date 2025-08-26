-- +goose Up
CREATE TABLE IF NOT EXISTS characters (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    age integer NOT NULL,
    description TEXT NOT NULL,
    origin text NOT NULL,
    race text NOT NULL DEFAULT 'human',
    bounty bigint ,
    episode int NOT NULL,
    time_skip text NOT NULL,
    CONSTRAINT unique_name_timeskip UNIQUE (name, time_skip)
);

-- +goose Down
DROP TABLE characters;
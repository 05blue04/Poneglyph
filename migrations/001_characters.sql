-- +goose Up
CREATE TABLE IF NOT EXISTS characters (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    description TEXT NOT NULL,
    age integer NOT NULL,
    bounty bigint ,
    origin text NOT NULL,
    Episode int NOT NULL,
    race text NOT NULL DEFAULT 'human',
    organizations text[],
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
    time_skip text NOT NULL;
    CONSTRAINT unique_name_timeskip UNIQUE (name, time_skip)
);

-- +goose Down
DROP TABLE characters;
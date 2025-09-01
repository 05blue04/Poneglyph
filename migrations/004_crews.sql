-- +goose Up
CREATE TABLE IF NOT EXISTS crews (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text UNIQUE NOT NULL,
    description text NOT NULL,
    ship_name text,
    captain_id bigserial REFERENCES characters(id) ON DELETE SET NULL
); 

CREATE TABLE IF NOT EXISTS crew_members (
    character_id bigserial REFERENCES characters(id),
    crew_id bigserial REFERENCES crews(id),
    role text,
    PRIMARY KEY (character_id, crew_id)
);

-- +goose Down
DROP TABLE IF EXISTS crew_members;
DROP TABLE IF EXISTS crews;

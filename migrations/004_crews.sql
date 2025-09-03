-- +goose Up
CREATE TABLE IF NOT EXISTS crews (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text UNIQUE NOT NULL,
    description text NOT NULL,
    ship_name text,
    captain_id bigint REFERENCES characters(id) ON DELETE SET NULL,
    captain_name text,
    total_bounty bigint
); 

CREATE TABLE IF NOT EXISTS crew_members (
    character_id bigint REFERENCES characters(id) ON DELETE CASCADE,
    crew_id bigint REFERENCES crews(id) ON DELETE CASCADE,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    PRIMARY KEY (character_id, crew_id)
);

-- +goose Down
DROP TABLE IF EXISTS crew_members;
DROP TABLE IF EXISTS crews;

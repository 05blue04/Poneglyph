CREATE TABLE IF NOT EXISTS crews (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    description text NOT NULL,
    ship_name text,
    captain_id bigserial INT REFERENCES characters(id) ON DELETE SET NULL,
    total_bounty bigint NOT NULL,
    episode int NOT NULL,
    time_skip text NOT NULL
    CONSTRAINT unique_crew UNIQUE (name, time_skip)
); 
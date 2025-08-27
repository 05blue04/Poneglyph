CREATE TABLE IF NOT EXISTS devilfruits(
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text UNIQUE NOT NULL,
    description TEXT NOT NULL,
    type TEXT NOT NULL,
    character_id bigserial REFERENCES characters(id) ON DELETE SET NULL,
    previousOwners text[],
    episode int NOT NULL
);
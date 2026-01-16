CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE incidents (
    id              BIGSERIAL PRIMARY KEY,
    type            VARCHAR(50) NOT NULL,
    description     TEXT,
    location        GEOGRAPHY(POINT, 4326) NOT NULL,
    radius_meters   INTEGER NOT NULL CHECK (radius_meters > 0),
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- CREATE INDEX idx_incidents_location ON incidents USING GIST (location);
-- CREATE INDEX idx_incidents_active_time ON incidents (is_active, starts_at, ends_at);

CREATE TABLE location_checks (
    id          BIGSERIAL PRIMARY KEY,
    user_id     INT NOT NULL,
    location    GEOGRAPHY(POINT, 4326) NOT NULL,
    has_danger  BOOLEAN NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
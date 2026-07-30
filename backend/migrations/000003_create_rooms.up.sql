-- Migration: 000003_create_rooms
-- Domain:    room (bounded context: property)
-- Table:     rooms
-- Created:   2026-07-30
-- Depends on: 000001_create_properties

CREATE TABLE IF NOT EXISTS rooms (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    property_id  UUID        NOT NULL REFERENCES properties (id) ON DELETE CASCADE,
    name         VARCHAR(255) NOT NULL,
    floor        INTEGER     NOT NULL DEFAULT 1,
    capacity     INTEGER     NOT NULL DEFAULT 1,
    price        NUMERIC(15, 2) NOT NULL DEFAULT 0,
    is_available BOOLEAN     NOT NULL DEFAULT true,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at   TIMESTAMPTZ
);

-- indexes
CREATE INDEX IF NOT EXISTS idx_rooms_property_id  ON rooms (property_id);
CREATE INDEX IF NOT EXISTS idx_rooms_is_available ON rooms (is_available);
CREATE INDEX IF NOT EXISTS idx_rooms_deleted_at   ON rooms (deleted_at);
CREATE INDEX IF NOT EXISTS idx_rooms_created_at   ON rooms (created_at DESC);

-- auto-update updated_at — set_updated_at() already defined in 000001
CREATE TRIGGER trg_rooms_updated_at
    BEFORE UPDATE ON rooms
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Migration: 000001_create_properties
-- Domain:    property
-- Table:     properties
-- Created:   2026-07-30

CREATE TABLE IF NOT EXISTS properties (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    name          VARCHAR(255) NOT NULL,
    description   TEXT,
    address       TEXT,
    property_type VARCHAR(100) NOT NULL,
    is_active     BOOLEAN     NOT NULL DEFAULT true,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at    TIMESTAMPTZ
);

-- indexes
CREATE INDEX IF NOT EXISTS idx_properties_is_active  ON properties (is_active);
CREATE INDEX IF NOT EXISTS idx_properties_deleted_at ON properties (deleted_at);
CREATE INDEX IF NOT EXISTS idx_properties_created_at ON properties (created_at DESC);

-- auto-update updated_at on every row change
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_properties_updated_at
    BEFORE UPDATE ON properties
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

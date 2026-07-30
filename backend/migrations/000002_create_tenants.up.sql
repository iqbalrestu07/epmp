-- Migration: 000002_create_tenants
-- Domain:    tenant
-- Table:     tenants
-- Created:   2026-07-30

CREATE TABLE IF NOT EXISTS tenants (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name       VARCHAR(255) NOT NULL,
    email           VARCHAR(255) NOT NULL UNIQUE,
    phone           VARCHAR(50),
    identity_number VARCHAR(100) UNIQUE,
    is_active       BOOLEAN     NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ
);

-- indexes
CREATE INDEX IF NOT EXISTS idx_tenants_email      ON tenants (email);
CREATE INDEX IF NOT EXISTS idx_tenants_is_active  ON tenants (is_active);
CREATE INDEX IF NOT EXISTS idx_tenants_deleted_at ON tenants (deleted_at);
CREATE INDEX IF NOT EXISTS idx_tenants_created_at ON tenants (created_at DESC);

-- set_updated_at() is already defined in 000001 — reuse it
CREATE TRIGGER trg_tenants_updated_at
    BEFORE UPDATE ON tenants
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

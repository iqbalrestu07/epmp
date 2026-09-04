-- Migration: 000032_org_members_and_room_floor
-- Purpose  : 1) Add organization_members pivot table (user <-> org M:N)
--            2) Add created_by to organizations
--            3) Add floor_id to rooms (proper hierarchy: room -> floor -> building -> property)
--            4) Drop deprecated rooms.floor integer column
-- Created  : 2026-09-04

-- ─────────────────────────────────────────────────────────────────────────────
-- 1. organization_members  (M:N pivot: user <-> organization)
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS organization_members (
    id              TEXT         PRIMARY KEY,  -- ULID
    organization_id TEXT         NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id         TEXT         NOT NULL REFERENCES users(id)         ON DELETE CASCADE,
    role            VARCHAR(50)  NOT NULL DEFAULT 'member',
    -- role values: 'owner' | 'admin' | 'member'
    invited_by      TEXT         REFERENCES users(id),
    joined_at       TIMESTAMPTZ  NOT NULL DEFAULT now(),
    is_active       BOOLEAN      NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at      TIMESTAMPTZ,
    CONSTRAINT uq_org_members UNIQUE (organization_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_org_members_org_id  ON organization_members (organization_id);
CREATE INDEX IF NOT EXISTS idx_org_members_user_id ON organization_members (user_id);
CREATE INDEX IF NOT EXISTS idx_org_members_role    ON organization_members (role);

CREATE TRIGGER trg_org_members_updated_at
    BEFORE UPDATE ON organization_members
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ─────────────────────────────────────────────────────────────────────────────
-- 2. created_by already added in migration 000032
-- ─────────────────────────────────────────────────────────────────────────────

-- ─────────────────────────────────────────────────────────────────────────────
-- 3. floor_id already added in migration 000032; just add the index
-- ─────────────────────────────────────────────────────────────────────────────

CREATE INDEX IF NOT EXISTS idx_rooms_floor_id ON rooms (floor_id);

-- ─────────────────────────────────────────────────────────────────────────────
-- 4. Drop deprecated rooms.floor integer column
-- ─────────────────────────────────────────────────────────────────────────────
ALTER TABLE rooms DROP COLUMN IF EXISTS floor;

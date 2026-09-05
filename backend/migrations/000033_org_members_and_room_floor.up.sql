-- Migration: 000033_org_members_and_room_floor
-- Purpose  : Add organization_members pivot table (user <-> org M:N) with native UUID

CREATE TABLE IF NOT EXISTS organization_members (
    id              UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID         NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    user_id         UUID         NOT NULL REFERENCES users(id)         ON DELETE CASCADE,
    role            VARCHAR(50)  NOT NULL DEFAULT 'member',
    invited_by      UUID         REFERENCES users(id),
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

ALTER TABLE rooms DROP COLUMN IF EXISTS floor;

-- Migration: 000035_property_user_roles
-- Purpose  : Property-level staff assignments and roles with native UUID

CREATE TABLE IF NOT EXISTS property_user_roles (
    id              UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID         NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    property_id     UUID         NOT NULL REFERENCES properties(id)    ON DELETE CASCADE,
    user_id         UUID         NOT NULL REFERENCES users(id)         ON DELETE CASCADE,
    role_id         UUID         NOT NULL REFERENCES roles(id)         ON DELETE CASCADE,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
    CONSTRAINT uq_property_user_role UNIQUE (property_id, user_id, role_id)
);

CREATE INDEX IF NOT EXISTS idx_pur_property_id ON property_user_roles (property_id);
CREATE INDEX IF NOT EXISTS idx_pur_user_id     ON property_user_roles (user_id);
CREATE INDEX IF NOT EXISTS idx_pur_role_id     ON property_user_roles (role_id);
CREATE INDEX IF NOT EXISTS idx_pur_org_id      ON property_user_roles (organization_id);

CREATE TRIGGER trg_pur_updated_at
    BEFORE UPDATE ON property_user_roles
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Seed front_desk role if not exists
INSERT INTO roles (name, description, is_system) VALUES
    ('front_desk', 'Manage check-ins, reservations, and room bookings', true)
ON CONFLICT (name) DO NOTHING;

-- Seed front_desk permissions (reservations, occupancies, rooms, tenants)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM   roles r
JOIN   permissions p ON (p.resource IN ('room', 'tenant', 'reservation', 'occupancy'))
WHERE  r.name = 'front_desk'
ON CONFLICT DO NOTHING;

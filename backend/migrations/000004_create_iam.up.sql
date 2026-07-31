-- Migration: 000004_create_iam
-- Domain:    Identity & Access Management
-- Created:   2026-07-31

-- ─────────────────────────────────────────────────────────────────────────────
-- users
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS users (
    id            UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    email         VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name          VARCHAR(255) NOT NULL,
    is_active     BOOLEAN      NOT NULL DEFAULT true,
    last_login_at TIMESTAMPTZ,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at    TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS uidx_users_email        ON users (email) WHERE deleted_at IS NULL;
CREATE        INDEX IF NOT EXISTS idx_users_is_active     ON users (is_active);
CREATE        INDEX IF NOT EXISTS idx_users_deleted_at    ON users (deleted_at);

CREATE TRIGGER trg_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ─────────────────────────────────────────────────────────────────────────────
-- roles
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS roles (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    is_system   BOOLEAN      NOT NULL DEFAULT false,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS uidx_roles_name ON roles (name);

CREATE TRIGGER trg_roles_updated_at
    BEFORE UPDATE ON roles
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ─────────────────────────────────────────────────────────────────────────────
-- permissions  (resource:action pairs)
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS permissions (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    resource    VARCHAR(100) NOT NULL,
    action      VARCHAR(50)  NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS uidx_permissions_resource_action ON permissions (resource, action);

-- ─────────────────────────────────────────────────────────────────────────────
-- user_roles  (M:N)
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS user_roles (
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id    UUID        NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, role_id)
);

CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles (user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON user_roles (role_id);

-- ─────────────────────────────────────────────────────────────────────────────
-- role_permissions  (M:N)
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id       UUID        NOT NULL REFERENCES roles(id)       ON DELETE CASCADE,
    permission_id UUID        NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (role_id, permission_id)
);

CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id ON role_permissions (role_id);

-- ─────────────────────────────────────────────────────────────────────────────
-- refresh_tokens
-- ─────────────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(64) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS uidx_refresh_tokens_hash  ON refresh_tokens (token_hash);
CREATE        INDEX IF NOT EXISTS idx_refresh_tokens_user   ON refresh_tokens (user_id);

-- ─────────────────────────────────────────────────────────────────────────────
-- Seed: system permissions
-- ─────────────────────────────────────────────────────────────────────────────
INSERT INTO permissions (resource, action, description) VALUES
    -- properties
    ('property', 'read',   'View properties'),
    ('property', 'write',  'Create/update properties'),
    ('property', 'delete', 'Delete properties'),
    -- rooms
    ('room', 'read',   'View rooms'),
    ('room', 'write',  'Create/update rooms'),
    ('room', 'delete', 'Delete rooms'),
    -- tenants
    ('tenant', 'read',   'View tenants'),
    ('tenant', 'write',  'Create/update tenants'),
    ('tenant', 'delete', 'Delete tenants'),
    -- users
    ('user', 'read',   'View users'),
    ('user', 'write',  'Create/update users'),
    ('user', 'delete', 'Delete/deactivate users'),
    -- roles
    ('role', 'read',  'View roles'),
    ('role', 'write', 'Create/update roles'),
    ('role', 'delete','Delete roles'),
    -- reports
    ('report', 'read', 'View reports'),
    -- invoices
    ('invoice', 'read',  'View invoices'),
    ('invoice', 'write', 'Create/update invoices'),
    -- payments
    ('payment', 'read',  'View payments'),
    ('payment', 'write', 'Record payments'),
    -- assets
    ('asset', 'read',  'View assets'),
    ('asset', 'write', 'Create/update assets'),
    -- work orders
    ('workorder', 'read',  'View work orders'),
    ('workorder', 'write', 'Create/update work orders')
ON CONFLICT (resource, action) DO NOTHING;

-- ─────────────────────────────────────────────────────────────────────────────
-- Seed: system roles
-- ─────────────────────────────────────────────────────────────────────────────
INSERT INTO roles (name, description, is_system) VALUES
    ('super_admin',       'Full system access',                      true),
    ('property_manager',  'Manage properties, rooms, and tenants',   true),
    ('finance_staff',     'Manage invoices, payments, and reports',  true),
    ('maintenance_staff', 'Manage work orders and assets',           true)
ON CONFLICT (name) DO NOTHING;

-- super_admin gets ALL permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM   roles r
CROSS JOIN permissions p
WHERE  r.name = 'super_admin'
ON CONFLICT DO NOTHING;

-- property_manager permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM   roles r
JOIN   permissions p ON (p.resource IN ('property','room','tenant'))
WHERE  r.name = 'property_manager'
ON CONFLICT DO NOTHING;

-- finance_staff permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM   roles r
JOIN   permissions p
       ON (p.resource IN ('invoice','payment','report'))
WHERE  r.name = 'finance_staff'
ON CONFLICT DO NOTHING;

-- maintenance_staff permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM   roles r
JOIN   permissions p
       ON (p.resource = 'workorder' OR (p.resource = 'asset' AND p.action = 'read'))
WHERE  r.name = 'maintenance_staff'
ON CONFLICT DO NOTHING;

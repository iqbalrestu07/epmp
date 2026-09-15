-- Migration: 000040_create_audit_logs
-- Description: Append-only audit trail for all mutating API operations.

CREATE TABLE IF NOT EXISTS audit_logs (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID        REFERENCES organizations(id),
    user_id         UUID,
    user_email      VARCHAR(255),
    action          VARCHAR(20)  NOT NULL,
    module          VARCHAR(100) NOT NULL,
    entity_id       TEXT,
    method          VARCHAR(10)  NOT NULL,
    path            TEXT         NOT NULL,
    status_code     INT          NOT NULL,
    ip_address      VARCHAR(64),
    user_agent      TEXT,
    request_body    JSONB,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_audit_logs_org_created ON audit_logs (organization_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_logs_module      ON audit_logs (module);
CREATE INDEX IF NOT EXISTS idx_audit_logs_entity      ON audit_logs (entity_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user        ON audit_logs (user_id);

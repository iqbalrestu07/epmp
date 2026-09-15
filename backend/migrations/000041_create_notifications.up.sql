-- Migration: 000041_create_notifications
-- Description: In-app notification store for real-time user alerts.

CREATE TABLE IF NOT EXISTS notifications (
    id              UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID         NOT NULL REFERENCES organizations(id),
    user_id         UUID,
    type            VARCHAR(50)  NOT NULL DEFAULT 'info',
    title           VARCHAR(255) NOT NULL,
    message         TEXT         NOT NULL,
    link            TEXT,
    is_read         BOOLEAN      NOT NULL DEFAULT false,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_notifications_org_user ON notifications (organization_id, user_id, is_read);
CREATE INDEX IF NOT EXISTS idx_notifications_created  ON notifications (created_at DESC);

-- Migration: 000032_add_created_by_and_floor_id
-- Add created_by to organizations and floor_id to rooms with native UUID type

ALTER TABLE organizations ADD COLUMN IF NOT EXISTS created_by UUID REFERENCES users(id);
CREATE INDEX IF NOT EXISTS idx_organizations_created_by ON organizations(created_by);

ALTER TABLE rooms ADD COLUMN IF NOT EXISTS floor_id UUID REFERENCES floors(id);
CREATE INDEX IF NOT EXISTS idx_rooms_floor_id ON rooms(floor_id);

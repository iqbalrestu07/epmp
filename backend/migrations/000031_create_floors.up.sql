CREATE TABLE IF NOT EXISTS floors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id),
    building_id UUID NOT NULL REFERENCES buildings(id),
    name VARCHAR(100) NOT NULL,
    floor_number INT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);
CREATE TRIGGER set_updated_at BEFORE UPDATE ON floors FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE asset_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id UUID NOT NULL REFERENCES assets(id),
    room_id UUID NOT NULL REFERENCES rooms(id),
    assigned_date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TRIGGER set_updated_at
BEFORE UPDATE ON asset_assignments
FOR EACH ROW
EXECUTE FUNCTION trigger_set_updated_at();

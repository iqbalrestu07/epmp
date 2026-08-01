CREATE TABLE reservations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    property_id UUID NOT NULL REFERENCES properties(id),
    room_id UUID NOT NULL REFERENCES rooms(id),
    status VARCHAR(50) NOT NULL DEFAULT 'Draft',
    check_in_date TIMESTAMP WITH TIME ZONE NOT NULL,
    check_out_date TIMESTAMP WITH TIME ZONE,
    booking_fee DECIMAL(19, 4) NOT NULL DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TRIGGER set_updated_at
BEFORE UPDATE ON reservations
FOR EACH ROW
EXECUTE FUNCTION trigger_set_updated_at();

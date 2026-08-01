CREATE TABLE tenant_identities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    identity_type VARCHAR(50) NOT NULL DEFAULT 'KTP',
    identity_number VARCHAR(100) NOT NULL,
    file_url VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TRIGGER set_updated_at
BEFORE UPDATE ON tenant_identities
FOR EACH ROW
EXECUTE FUNCTION trigger_set_updated_at();

CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL REFERENCES invoices(id),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    amount DECIMAL(19, 4) NOT NULL DEFAULT 0,
    payment_date TIMESTAMP WITH TIME ZONE NOT NULL,
    payment_method VARCHAR(50) NOT NULL DEFAULT 'Transfer',
    status VARCHAR(50) NOT NULL DEFAULT 'Pending',
    reference_number VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TRIGGER set_updated_at
BEFORE UPDATE ON payments
FOR EACH ROW
EXECUTE FUNCTION trigger_set_updated_at();

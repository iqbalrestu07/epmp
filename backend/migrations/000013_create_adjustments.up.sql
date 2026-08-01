CREATE TABLE adjustments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL REFERENCES invoices(id),
    adjustment_type VARCHAR(50) NOT NULL DEFAULT 'Credit',
    amount DECIMAL(19, 4) NOT NULL DEFAULT 0,
    adjustment_date TIMESTAMP WITH TIME ZONE NOT NULL,
    reason TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TRIGGER set_updated_at
BEFORE UPDATE ON adjustments
FOR EACH ROW
EXECUTE FUNCTION trigger_set_updated_at();

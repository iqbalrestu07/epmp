CREATE OR REPLACE FUNCTION sync_invoice_status_on_payment()
RETURNS TRIGGER AS $$
DECLARE
    v_invoice_id UUID;
    v_total_paid NUMERIC(19,4);
    v_inv_amount NUMERIC(19,4);
    v_pay_method VARCHAR(255);
    v_pay_date TIMESTAMPTZ;
BEGIN
    IF TG_OP = 'DELETE' THEN
        v_invoice_id := OLD.invoice_id;
    ELSE
        v_invoice_id := NEW.invoice_id;
        v_pay_method := NEW.payment_method;
        v_pay_date := NEW.payment_date;
    END IF;

    IF v_invoice_id IS NOT NULL THEN
        SELECT amount INTO v_inv_amount FROM invoices WHERE id = v_invoice_id AND deleted_at IS NULL;
        
        IF v_inv_amount IS NOT NULL THEN
            SELECT COALESCE(SUM(amount), 0) INTO v_total_paid
            FROM payments
            WHERE invoice_id = v_invoice_id
              AND (status = 'Success' OR status = 'Completed')
              AND deleted_at IS NULL;

            IF v_total_paid >= v_inv_amount AND v_inv_amount > 0 THEN
                UPDATE invoices
                SET status = 'Paid',
                    paid_date = COALESCE(v_pay_date, NOW()),
                    payment_method = COALESCE(v_pay_method, payment_method, 'Transfer'),
                    updated_at = NOW()
                WHERE id = v_invoice_id;
            ELSIF v_total_paid > 0 THEN
                UPDATE invoices
                SET status = 'Partial',
                    updated_at = NOW()
                WHERE id = v_invoice_id;
            ELSE
                UPDATE invoices
                SET status = CASE WHEN due_date < NOW() THEN 'Overdue' ELSE 'Unpaid' END,
                    paid_date = NULL,
                    updated_at = NOW()
                WHERE id = v_invoice_id;
            END IF;
        END IF;
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_sync_invoice_status_on_payment ON payments;
CREATE TRIGGER trg_sync_invoice_status_on_payment
AFTER INSERT OR UPDATE OR DELETE ON payments
FOR EACH ROW
EXECUTE FUNCTION sync_invoice_status_on_payment();

-- Sync existing invoices
UPDATE invoices i
SET status = 'Paid',
    paid_date = p.payment_date,
    payment_method = p.payment_method
FROM (
    SELECT invoice_id, MAX(payment_date) as payment_date, MAX(payment_method) as payment_method, SUM(amount) as total_paid
    FROM payments
    WHERE (status = 'Success' OR status = 'Completed') AND deleted_at IS NULL
    GROUP BY invoice_id
) p
WHERE i.id = p.invoice_id AND p.total_paid >= i.amount;

-- Migration: 000036_add_currency_support down

ALTER TABLE organizations DROP COLUMN IF EXISTS currency;
ALTER TABLE properties DROP COLUMN IF EXISTS currency;
ALTER TABLE room_types DROP COLUMN IF EXISTS currency;
ALTER TABLE rooms DROP COLUMN IF EXISTS currency;
ALTER TABLE contracts DROP COLUMN IF EXISTS currency;
ALTER TABLE invoices DROP COLUMN IF EXISTS currency;
ALTER TABLE payments DROP COLUMN IF EXISTS currency;
ALTER TABLE charges DROP COLUMN IF EXISTS currency;
ALTER TABLE deposits DROP COLUMN IF EXISTS currency;
ALTER TABLE refunds DROP COLUMN IF EXISTS currency;
ALTER TABLE adjustments DROP COLUMN IF EXISTS currency;
ALTER TABLE penalties DROP COLUMN IF EXISTS currency;

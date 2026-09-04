-- Reverse migration 000033
-- NOTE: Casting TEXT back to UUID only works if all stored values are UUID-format strings.
-- If ULID values have been inserted, this DOWN migration will NOT work.
-- In that case, restore from backup.
SELECT 'DOWN migration for 000033 is not fully reversible once ULID values exist. Restore from backup.';

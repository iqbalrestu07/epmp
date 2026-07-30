-- Rollback: 000001_create_properties

DROP TRIGGER IF EXISTS trg_properties_updated_at ON properties;
DROP FUNCTION IF EXISTS set_updated_at();
DROP TABLE IF EXISTS properties;

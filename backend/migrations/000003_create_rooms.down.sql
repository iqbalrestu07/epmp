-- Rollback: 000003_create_rooms

DROP TRIGGER IF EXISTS trg_rooms_updated_at ON rooms;
DROP TABLE IF EXISTS rooms;

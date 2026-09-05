ALTER TABLE rooms DROP CONSTRAINT IF EXISTS fk_room_room_type;
DROP INDEX IF EXISTS idx_rooms_room_type_id;
ALTER TABLE rooms DROP COLUMN IF EXISTS room_type_id;

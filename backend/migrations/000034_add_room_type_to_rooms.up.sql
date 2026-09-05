ALTER TABLE rooms ADD COLUMN IF NOT EXISTS room_type_id text;
DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_room_room_type') THEN
    ALTER TABLE rooms ADD CONSTRAINT fk_room_room_type FOREIGN KEY (room_type_id) REFERENCES room_types(id) ON DELETE SET NULL;
  END IF;
END $$;
CREATE INDEX IF NOT EXISTS idx_rooms_room_type_id ON rooms(room_type_id);

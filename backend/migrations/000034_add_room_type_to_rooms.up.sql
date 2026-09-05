ALTER TABLE rooms ADD COLUMN IF NOT EXISTS room_type_id UUID REFERENCES room_types(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS idx_rooms_room_type_id ON rooms(room_type_id);

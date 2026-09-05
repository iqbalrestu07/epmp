-- Migration: 000037_room_status_and_lifecycle_sync down

DROP TRIGGER IF EXISTS trg_reservations_sync_room_status ON reservations;
DROP TRIGGER IF EXISTS trg_contracts_sync_room_status ON contracts;
DROP TRIGGER IF EXISTS trg_occupancies_sync_room_status ON occupancies;
DROP FUNCTION IF EXISTS sync_room_status_from_occupancy();

ALTER TABLE rooms DROP COLUMN IF EXISTS status;

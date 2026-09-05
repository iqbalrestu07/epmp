-- Migration: 000037_room_status_and_lifecycle_sync
-- Description: Add status column to rooms and auto-sync room status across Reservation -> Contract -> Occupancy

ALTER TABLE rooms ADD COLUMN IF NOT EXISTS status VARCHAR(50) NOT NULL DEFAULT 'Available';

-- Update existing rooms based on current occupancy
UPDATE rooms r
SET status = 'Occupied', is_available = false
WHERE EXISTS (
    SELECT 1 FROM occupancies o 
    WHERE o.room_id = r.id AND o.status = 'CheckedIn' AND o.deleted_at IS NULL
);

UPDATE rooms r
SET status = 'Reserved', is_available = false
WHERE r.status = 'Available' AND (
    EXISTS (SELECT 1 FROM contracts c WHERE c.room_id = r.id AND c.status = 'Active' AND c.deleted_at IS NULL)
    OR EXISTS (SELECT 1 FROM reservations res WHERE res.room_id = r.id AND res.status IN ('Confirmed', 'Pending') AND res.deleted_at IS NULL)
);

-- Trigger function to automatically maintain room availability and status
CREATE OR REPLACE FUNCTION sync_room_status_from_occupancy()
RETURNS TRIGGER AS $$
DECLARE
    target_room_id UUID;
    has_checked_in BOOLEAN;
    has_active_contract BOOLEAN;
    has_confirmed_reservation BOOLEAN;
    current_room_status VARCHAR(50);
BEGIN
    target_room_id := COALESCE(NEW.room_id, OLD.room_id);
    IF target_room_id IS NULL THEN
        RETURN NEW;
    END IF;

    SELECT status INTO current_room_status FROM rooms WHERE id = target_room_id;
    IF current_room_status = 'Maintenance' THEN
        RETURN NEW;
    END IF;

    SELECT EXISTS (
        SELECT 1 FROM occupancies 
        WHERE room_id = target_room_id 
          AND status = 'CheckedIn' 
          AND deleted_at IS NULL
    ) INTO has_checked_in;

    SELECT EXISTS (
        SELECT 1 FROM contracts 
        WHERE room_id = target_room_id 
          AND status = 'Active' 
          AND deleted_at IS NULL
    ) INTO has_active_contract;

    SELECT EXISTS (
        SELECT 1 FROM reservations 
        WHERE room_id = target_room_id 
          AND status IN ('Confirmed', 'Pending') 
          AND deleted_at IS NULL
    ) INTO has_confirmed_reservation;

    IF has_checked_in THEN
        UPDATE rooms 
        SET status = 'Occupied', is_available = false 
        WHERE id = target_room_id;
    ELSIF has_active_contract OR has_confirmed_reservation THEN
        UPDATE rooms 
        SET status = 'Reserved', is_available = false 
        WHERE id = target_room_id;
    ELSE
        UPDATE rooms 
        SET status = 'Available', is_available = true 
        WHERE id = target_room_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_occupancies_sync_room_status ON occupancies;
CREATE TRIGGER trg_occupancies_sync_room_status
AFTER INSERT OR UPDATE OR DELETE ON occupancies
FOR EACH ROW EXECUTE FUNCTION sync_room_status_from_occupancy();

DROP TRIGGER IF EXISTS trg_contracts_sync_room_status ON contracts;
CREATE TRIGGER trg_contracts_sync_room_status
AFTER INSERT OR UPDATE OR DELETE ON contracts
FOR EACH ROW EXECUTE FUNCTION sync_room_status_from_occupancy();

DROP TRIGGER IF EXISTS trg_reservations_sync_room_status ON reservations;
CREATE TRIGGER trg_reservations_sync_room_status
AFTER INSERT OR UPDATE OR DELETE ON reservations
FOR EACH ROW EXECUTE FUNCTION sync_room_status_from_occupancy();

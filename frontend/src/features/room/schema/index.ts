import { z } from "zod";

export const roomSchema = z.object({
  organization_id: z.string().optional(),
  id: z.string(),
  name: z.string().min(1, "Name is required").max(100),
  property_id: z.string().min(1, "Property is required"),
  floor_id: z.string().optional().nullable(),
  room_type_id: z.string().optional().nullable(),
  capacity: z.coerce.number().int().min(1, "Capacity must be at least 1"),
  price: z.coerce.number().min(0, "Price must be >= 0"),
  is_available: z.boolean(),
  created_at: z.string().optional(),
  updated_at: z.string().optional(),
});

export const createRoomSchema = roomSchema.omit({
  id: true,
  organization_id: true,
  created_at: true,
  updated_at: true,
});

export const updateRoomSchema = createRoomSchema.partial();

export type RoomFormData = z.infer<typeof roomSchema>;
export type CreateRoomFormData = z.infer<typeof createRoomSchema>;
export type UpdateRoomFormData = z.infer<typeof updateRoomSchema>;

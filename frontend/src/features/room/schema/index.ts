import { z } from "zod";

export const roomSchema = z.object({
  id: z.string(),
  name: z.string().max(100),
  floor: z.number(),
  capacity: z.number(),
  price: z.number(),
  is_available: z.boolean(),
  property_id: z.string(),
  created_at: z.string().optional(),
  updated_at: z.string().optional(),
});

export const createRoomSchema = roomSchema.omit({
  id: true,
  created_at: true,
  updated_at: true,
});

export const updateRoomSchema = createRoomSchema.partial();

export type RoomFormData = z.infer<typeof roomSchema>;
export type CreateRoomFormData = z.infer<typeof createRoomSchema>;
export type UpdateRoomFormData = z.infer<typeof updateRoomSchema>;

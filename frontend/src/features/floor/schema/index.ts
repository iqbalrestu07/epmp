import { z } from "zod";

export const floorSchema = z.object({
  id: z.string(),
  organization_id: z.string(),
  building_id: z.string().min(1, "Building is required"),
  name: z.string().min(1, "Name is required").max(100),
  floor_number: z.number().int().min(0, "Floor number must be >= 0"),
  is_active: z.boolean(),
  created_at: z.string().optional(),
  updated_at: z.string().optional(),
});

export const createFloorSchema = floorSchema.omit({
  id: true,
  organization_id: true,
  created_at: true,
  updated_at: true,
});

export const updateFloorSchema = createFloorSchema.partial();

export type FloorFormData = z.infer<typeof floorSchema>;
export type CreateFloorFormData = z.infer<typeof createFloorSchema>;
export type UpdateFloorFormData = z.infer<typeof updateFloorSchema>;

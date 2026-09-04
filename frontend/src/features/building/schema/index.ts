import { z } from "zod";

export const buildingSchema = z.object({
  id: z.string(),
  property_id: z.string().min(1, "Property is required"),
  name: z.string().min(1, "Name is required").max(100),
  total_floors: z.coerce.number().min(1, "At least 1 floor required"),
  created_at: z.string().optional(),
  updated_at: z.string().optional(),
});

export const createBuildingSchema = buildingSchema.omit({
  id: true,
  created_at: true,
  updated_at: true,
});

export const updateBuildingSchema = createBuildingSchema.partial();

export type BuildingFormData = z.infer<typeof buildingSchema>;
export type CreateBuildingFormData = z.infer<typeof createBuildingSchema>;
export type UpdateBuildingFormData = z.infer<typeof updateBuildingSchema>;

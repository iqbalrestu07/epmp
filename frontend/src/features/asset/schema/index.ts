import { z } from "zod";

export const assetSchema = z.object({
  organization_id: z.string(),
  id: z.string(),
  property_id: z.string().min(1, "Property is required"),
  name: z.string().min(1, "Name is required").max(100),
  category: z.string().min(1, "Category is required").max(100),
  status: z.enum(["Available", "Assigned", "Maintenance", "Disposed"]).default("Available"),
  purchase_price: z.number().min(0),
  created_at: z.string().optional(),
  updated_at: z.string().optional(),
});

export const createAssetSchema = assetSchema.omit({
  id: true,
  organization_id: true, // Organization ID is injected by backend/context
  created_at: true,
  updated_at: true,
});

export const updateAssetSchema = createAssetSchema.partial();

export type AssetFormData = z.infer<typeof assetSchema>;
export type CreateAssetFormData = z.infer<typeof createAssetSchema>;
export type UpdateAssetFormData = z.infer<typeof updateAssetSchema>;

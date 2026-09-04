import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { usePropertys } from "../../property/hooks";
import {
  createAssetSchema,
  type CreateAssetFormData,
} from "../schema";

interface AssetFormProps {
  onSubmit: (data: CreateAssetFormData) => void;
  defaultValues?: Partial<CreateAssetFormData>;
  isSubmitting?: boolean;
}

export function AssetForm({
  onSubmit,
  defaultValues,
  isSubmitting,
}: AssetFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<CreateAssetFormData>({
    resolver: zodResolver(createAssetSchema),
    defaultValues: {
      status: "Available",
      ...defaultValues,
    },
  });

  const { data: propertiesData } = usePropertys({ per_page: 100 });
  const properties = propertiesData?.data || [];

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="property_id">Property</Label>
        {properties.length > 0 ? (
          <select
            id="property_id"
            {...register("property_id")}
            className="w-full rounded-md border border-gray-300 px-3 py-2 text-sm"
          >
            <option value="">Select a property…</option>
            {properties.map(p => (
              <option key={p.id} value={p.id}>{p.name}</option>
            ))}
          </select>
        ) : (
          <Input
            id="property_id"
            placeholder="Enter property ID"
            {...register("property_id")}
          />
        )}
        {errors.property_id && (
          <p className="text-sm text-red-500">{errors.property_id.message}</p>
        )}
      </div>

      <div className="space-y-2">
        <Label htmlFor="name">Asset Name</Label>
        <Input
          id="name"
          placeholder="e.g. Samsung 43 inch TV"
          {...register("name")}
        />
        {errors.name && (
          <p className="text-sm text-red-500">{errors.name.message}</p>
        )}
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div className="space-y-2">
          <Label htmlFor="category">Category</Label>
          <select
            id="category"
            {...register("category")}
            className="w-full rounded-md border border-gray-300 px-3 py-2 text-sm"
          >
            <option value="">Select category...</option>
            <option value="Furniture">Furniture</option>
            <option value="Electronics">Electronics</option>
            <option value="Appliance">Appliance</option>
            <option value="Decor">Decor</option>
            <option value="Other">Other</option>
          </select>
          {errors.category && (
            <p className="text-sm text-red-500">{errors.category.message}</p>
          )}
        </div>

        <div className="space-y-2">
          <Label htmlFor="status">Status</Label>
          <select
            id="status"
            {...register("status")}
            className="w-full rounded-md border border-gray-300 px-3 py-2 text-sm"
          >
            <option value="Available">Available</option>
            <option value="Assigned">Assigned</option>
            <option value="Maintenance">Maintenance</option>
            <option value="Disposed">Disposed</option>
          </select>
          {errors.status && (
            <p className="text-sm text-red-500">{errors.status.message}</p>
          )}
        </div>
      </div>

      <div className="space-y-2">
        <Label htmlFor="purchase_price">Purchase Price</Label>
        <Input
          id="purchase_price"
          type="number"
          step="0.01"
          placeholder="0.00"
          {...register("purchase_price", { valueAsNumber: true })}
        />
        {errors.purchase_price && (
          <p className="text-sm text-red-500">{errors.purchase_price.message}</p>
        )}
      </div>

      <Button type="submit" disabled={isSubmitting}>
        {isSubmitting ? "Saving..." : "Save Asset"}
      </Button>
    </form>
  );
}

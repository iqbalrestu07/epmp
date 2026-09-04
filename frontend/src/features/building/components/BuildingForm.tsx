import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  createBuildingSchema,
  type CreateBuildingFormData,
} from "../schema";
import { usePropertys } from "../../property/hooks";

interface BuildingFormProps {
  onSubmit: (data: CreateBuildingFormData) => void;
  defaultValues?: Partial<CreateBuildingFormData>;
  isSubmitting?: boolean;
}

export function BuildingForm({
  onSubmit,
  defaultValues,
  isSubmitting,
}: BuildingFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<CreateBuildingFormData>({
    resolver: zodResolver(createBuildingSchema),
    defaultValues: {
      total_floors: 1,
      ...defaultValues,
    },
  });

  const { data: propertiesData } = usePropertys({ per_page: 100 });
  const properties = propertiesData?.data ?? [];

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4 max-w-md">
      <div className="space-y-2">
        <Label htmlFor="property_id">Property</Label>
        <select
          id="property_id"
          {...register("property_id")}
          className="w-full rounded-md border border-gray-300 px-3 py-2 text-sm"
        >
          <option value="">Select a property…</option>
          {properties.map((p) => (
            <option key={p.id} value={p.id}>
              {p.name}
            </option>
          ))}
        </select>
        {errors.property_id && (
          <p className="text-sm text-red-500">{errors.property_id.message}</p>
        )}
      </div>

      <div className="space-y-2">
        <Label htmlFor="name">Building Name</Label>
        <Input
          id="name"
          {...register("name")}
          placeholder="e.g. Tower A"
        />
        {errors.name && (
          <p className="text-sm text-red-500">{errors.name.message}</p>
        )}
      </div>

      <div className="space-y-2">
        <Label htmlFor="total_floors">Total Floors</Label>
        <Input
          id="total_floors"
          type="number"
          min={1}
          {...register("total_floors")}
        />
        {errors.total_floors && (
          <p className="text-sm text-red-500">{errors.total_floors.message}</p>
        )}
      </div>

      <Button type="submit" disabled={isSubmitting}>
        {isSubmitting ? "Saving..." : "Save"}
      </Button>
    </form>
  );
}

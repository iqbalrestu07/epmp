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
  const properties = Array.isArray(propertiesData?.data)
    ? propertiesData.data
    : Array.isArray(propertiesData)
    ? propertiesData
    : [];

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
      <div className="space-y-2">
        <Label htmlFor="property_id">Property</Label>
        <select
          id="property_id"
          {...register("property_id")}
          className="w-full rounded-lg border border-slate-300 bg-white px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-orange/30 text-slate-800"
        >
          <option value="">Select a property…</option>
          {properties.map((p) => (
            <option key={p.id} value={p.id}>
              {p.name}
            </option>
          ))}
        </select>
        {errors.property_id && (
          <p className="text-xs text-red-600 mt-1">{errors.property_id.message}</p>
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
          <p className="text-xs text-red-600 mt-1">{errors.name.message}</p>
        )}
      </div>

      <div className="space-y-2">
        <Label htmlFor="total_floors">Total Floors (Planned Capacity)</Label>
        <Input
          id="total_floors"
          type="number"
          min={1}
          {...register("total_floors")}
        />
        {errors.total_floors && (
          <p className="text-xs text-red-600 mt-1">{errors.total_floors.message}</p>
        )}
        <p className="text-xs text-slate-400">Maximum number of floors allowed. Actual floors are created in the Floors module.</p>
      </div>

      <Button type="submit" disabled={isSubmitting} className="bg-orange hover:bg-orange/90 text-white w-full sm:w-auto">
        {isSubmitting ? "Saving..." : "Save"}
      </Button>
    </form>
  );
}

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  createPropertySchema,
  type CreatePropertyFormData,
} from "../schema";

interface PropertyFormProps {
  onSubmit: (data: CreatePropertyFormData) => void;
  defaultValues?: Partial<CreatePropertyFormData>;
  isSubmitting?: boolean;
}

export function PropertyForm({
  onSubmit,
  defaultValues,
  isSubmitting,
}: PropertyFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<CreatePropertyFormData>({
    resolver: zodResolver(createPropertySchema),
    defaultValues: {
      is_active: true,
      property_type: "boarding_house",
      ...defaultValues,
    },
  });

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
      <div className="space-y-2">
        <Label htmlFor="name">Name</Label>
        <Input
          id="name"
          {...register("name")}
          placeholder="e.g. Tower Alpha"
        />
        {errors.name && (
          <p className="text-xs text-red-600 mt-1">{errors.name.message}</p>
        )}
      </div>
      <div className="space-y-2">
        <Label htmlFor="description">Description</Label>
        <Input
          id="description"
          {...register("description")}
          placeholder="Optional description"
        />
      </div>
      <div className="space-y-2">
        <Label htmlFor="address">Address</Label>
        <Input
          id="address"
          {...register("address")}
          placeholder="Street address"
        />
      </div>
      <div className="space-y-2">
        <Label htmlFor="property_type">Property Type</Label>
        <select
          id="property_type"
          {...register("property_type")}
          className="w-full rounded-lg border border-slate-300 bg-white px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-orange/30 text-slate-800"
        >
          <option value="boarding_house">Boarding House</option>
          <option value="apartment">Apartment</option>
          <option value="villa">Villa</option>
          <option value="warehouse">Warehouse</option>
        </select>
        {errors.property_type && (
          <p className="text-xs text-red-600 mt-1">{errors.property_type.message}</p>
        )}
      </div>
      <div className="flex items-center space-x-2">
        <input
          id="is_active"
          type="checkbox"
          {...register("is_active")}
          className="h-4 w-4 rounded border-slate-300"
        />
        <Label htmlFor="is_active">Active</Label>
      </div>

      <Button type="submit" disabled={isSubmitting} className="bg-orange hover:bg-orange/90 text-white w-full sm:w-auto">
        {isSubmitting ? "Saving..." : "Save"}
      </Button>
    </form>
  );
}

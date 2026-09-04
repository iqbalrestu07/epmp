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
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="name">Name</Label>
        <Input
          id="name"
          {...register("name")}
          placeholder="e.g. Tower Alpha"
        />
        {errors.name && (
          <p className="text-sm text-red-500">{errors.name.message}</p>
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
          className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        >
          <option value="boarding_house">Boarding House</option>
          <option value="apartment">Apartment</option>
          <option value="villa">Villa</option>
          <option value="warehouse">Warehouse</option>
        </select>
        {errors.property_type && (
          <p className="text-sm text-red-500">{errors.property_type.message}</p>
        )}
      </div>
      <div className="flex items-center space-x-2">
        <input
          id="is_active"
          type="checkbox"
          {...register("is_active")}
          className="h-4 w-4 rounded border-gray-300"
        />
        <Label htmlFor="is_active">Active</Label>
      </div>

      <Button type="submit" disabled={isSubmitting}>
        {isSubmitting ? "Saving..." : "Save"}
      </Button>
    </form>
  );
}

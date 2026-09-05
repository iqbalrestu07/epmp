import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  createFloorSchema,
  type CreateFloorFormData,
} from "../schema";
import { useBuildings } from "../../building/hooks";

interface FloorFormProps {
  onSubmit: (data: CreateFloorFormData) => void;
  defaultValues?: Partial<CreateFloorFormData>;
  isSubmitting?: boolean;
  buildings?: { id: string; name: string }[];
}

export function FloorForm({
  onSubmit,
  defaultValues,
  isSubmitting,
  buildings: propBuildings,
}: FloorFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<CreateFloorFormData>({
    resolver: zodResolver(createFloorSchema),
    defaultValues: {
      is_active: true,
      ...defaultValues,
    },
  });

  const { data: buildingsData } = useBuildings({ per_page: 100 });
  const buildings = propBuildings || (Array.isArray(buildingsData?.data) ? buildingsData.data : Array.isArray(buildingsData) ? buildingsData : []);

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
      <div className="space-y-2">
        <Label htmlFor="building_id">Building</Label>
        {buildings.length > 0 ? (
          <select
            id="building_id"
            {...register("building_id")}
            className="w-full rounded-lg border border-slate-300 bg-white px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-orange/30"
          >
            <option value="">Select a building…</option>
            {buildings.map(b => (
              <option key={b.id} value={b.id}>{b.name}</option>
            ))}
          </select>
        ) : (
          <Input
            id="building_id"
            placeholder="Enter building ID"
            {...register("building_id")}
          />
        )}
        {errors.building_id && (
          <p className="text-xs text-red-600 mt-1">{errors.building_id.message}</p>
        )}
      </div>

      <div className="space-y-2">
        <Label htmlFor="name">Floor Name</Label>
        <Input
          id="name"
          placeholder="e.g. Ground Floor, 1st Floor, Basement"
          {...register("name")}
        />
        {errors.name && (
          <p className="text-xs text-red-600 mt-1">{errors.name.message}</p>
        )}
      </div>

      <div className="space-y-2">
        <Label htmlFor="floor_number">Floor Number</Label>
        <Input
          id="floor_number"
          type="number"
          placeholder="0, 1, 2, -1 for basement…"
          {...register("floor_number", { valueAsNumber: true })}
        />
        {errors.floor_number && (
          <p className="text-xs text-red-600 mt-1">{errors.floor_number.message}</p>
        )}
      </div>

      <div className="space-y-2">
        <Label htmlFor="is_active">Active</Label>
        <label className="flex items-center gap-3 cursor-pointer">
          <input
            id="is_active"
            type="checkbox"
            {...register("is_active")}
            className="w-5 h-5 rounded border-black/20 text-orange focus:ring-orange/30"
          />
          <span className="text-sm text-slate-600">Floor is active and available</span>
        </label>
        {errors.is_active && (
          <p className="text-xs text-red-600 mt-1">{errors.is_active.message}</p>
        )}
      </div>

      <Button type="submit" disabled={isSubmitting}>
        {isSubmitting ? "Saving…" : "Save Floor"}
      </Button>
    </form>
  );
}

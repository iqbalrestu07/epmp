import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  createRoomSchema,
  type CreateRoomFormData,
} from "../schema";
import { usePropertys } from "../../property/hooks";
import { useFloors } from "../../floor/hooks";

interface RoomFormProps {
  onSubmit: (data: CreateRoomFormData) => void;
  defaultValues?: Partial<CreateRoomFormData>;
  isSubmitting?: boolean;
  properties?: { id: string; name: string }[];
  floors?: { id: string; name: string }[];
}

export function RoomForm({
  onSubmit,
  defaultValues,
  isSubmitting,
  properties: propProperties,
  floors: propFloors,
}: RoomFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<CreateRoomFormData>({
    resolver: zodResolver(createRoomSchema),
    defaultValues: {
      is_available: true,
      ...defaultValues,
    },
  });

  const { data: propertiesData } = usePropertys({ per_page: 100 });
  const properties = propProperties || propertiesData?.data || [];

  const { data: floorsData } = useFloors({ per_page: 100 });
  const floors = propFloors || floorsData?.data || [];

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
      <div className="space-y-2">
        <Label htmlFor="property_id">Property</Label>
        {properties.length > 0 ? (
          <select
            id="property_id"
            {...register("property_id")}
            className="w-full rounded-lg border border-black/10 bg-white px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-orange/30"
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
        <Label htmlFor="floor_id">Floor</Label>
        {floors.length > 0 ? (
          <select
            id="floor_id"
            {...register("floor_id")}
            className="w-full rounded-lg border border-black/10 bg-white px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-orange/30"
          >
            <option value="">Select a floor…</option>
            {floors.map(f => (
              <option key={f.id} value={f.id}>{f.name}</option>
            ))}
          </select>
        ) : (
          <Input
            id="floor_id"
            placeholder="Enter floor ID"
            {...register("floor_id")}
          />
        )}
        {errors.floor_id && (
          <p className="text-sm text-red-500">{errors.floor_id.message}</p>
        )}
      </div>

      <div className="space-y-2">
        <Label htmlFor="name">Room Name</Label>
        <Input
          id="name"
          placeholder="e.g. Room 101, Suite A, Studio 5"
          {...register("name")}
        />
        {errors.name && (
          <p className="text-sm text-red-500">{errors.name.message}</p>
        )}
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div className="space-y-2">
          <Label htmlFor="capacity">Capacity</Label>
          <Input
            id="capacity"
            type="number"
            placeholder="1, 2, 3…"
            {...register("capacity", { valueAsNumber: true })}
          />
          {errors.capacity && (
            <p className="text-sm text-red-500">{errors.capacity.message}</p>
          )}
        </div>

        <div className="space-y-2">
          <Label htmlFor="price">Price</Label>
          <Input
            id="price"
            type="number"
            step="0.01"
            placeholder="0.00"
            {...register("price", { valueAsNumber: true })}
          />
          {errors.price && (
            <p className="text-sm text-red-500">{errors.price.message}</p>
          )}
        </div>
      </div>

      <div className="space-y-2">
        <Label htmlFor="is_available">Available</Label>
        <label className="flex items-center gap-3 cursor-pointer">
          <input
            id="is_available"
            type="checkbox"
            {...register("is_available")}
            className="w-5 h-5 rounded border-black/20 text-orange focus:ring-orange/30"
          />
          <span className="text-sm text-black/60">Room is available for booking</span>
        </label>
        {errors.is_available && (
          <p className="text-sm text-red-500">{errors.is_available.message}</p>
        )}
      </div>

      <Button type="submit" disabled={isSubmitting}>
        {isSubmitting ? "Saving…" : "Save Room"}
      </Button>
    </form>
  );
}

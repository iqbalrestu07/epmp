import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useAssets } from "../../asset/hooks";
import { useRooms } from "../../room/hooks";
import {
  createAssetAssignmentSchema,
  type CreateAssetAssignmentFormData,
} from "../schema";

interface AssetAssignmentFormProps {
  onSubmit: (data: CreateAssetAssignmentFormData) => void;
  defaultValues?: Partial<CreateAssetAssignmentFormData>;
  isSubmitting?: boolean;
}

export function AssetAssignmentForm({
  onSubmit,
  defaultValues,
  isSubmitting,
}: AssetAssignmentFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<CreateAssetAssignmentFormData>({
    resolver: zodResolver(createAssetAssignmentSchema),
    defaultValues: {
      assigned_date: new Date().toISOString().split('T')[0], // YYYY-MM-DD
      ...defaultValues,
    },
  });

  const { data: assetsData } = useAssets({ per_page: 100 });
  const { data: roomsData } = useRooms({ per_page: 100 });
  
  const assets = assetsData?.data || [];
  const rooms = roomsData?.data || [];

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
      <div className="space-y-2">
        <Label htmlFor="asset_id">Asset</Label>
        <select
          id="asset_id"
          {...register("asset_id")}
          className="w-full rounded-lg border border-slate-300 bg-white px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-orange/30 text-slate-800"
        >
          <option value="">Select an asset...</option>
          {assets.map(a => (
            <option key={a.id} value={a.id}>{a.name} ({a.status})</option>
          ))}
        </select>
        {errors.asset_id && (
          <p className="text-xs text-red-600 mt-1">{errors.asset_id.message}</p>
        )}
      </div>

      <div className="space-y-2">
        <Label htmlFor="room_id">Room</Label>
        <select
          id="room_id"
          {...register("room_id")}
          className="w-full rounded-lg border border-slate-300 bg-white px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-orange/30 text-slate-800"
        >
          <option value="">Select a room...</option>
          {rooms.map(r => (
            <option key={r.id} value={r.id}>{r.name}</option>
          ))}
        </select>
        {errors.room_id && (
          <p className="text-xs text-red-600 mt-1">{errors.room_id.message}</p>
        )}
      </div>

      <div className="space-y-2">
        <Label htmlFor="assigned_date">Assigned Date</Label>
        <Input
          id="assigned_date"
          type="date"
          {...register("assigned_date")}
        />
        {errors.assigned_date && (
          <p className="text-xs text-red-600 mt-1">{errors.assigned_date.message}</p>
        )}
      </div>

      <Button type="submit" disabled={isSubmitting} className="bg-orange hover:bg-orange/90 text-white w-full sm:w-auto">
        {isSubmitting ? "Saving..." : "Save"}
      </Button>
    </form>
  );
}

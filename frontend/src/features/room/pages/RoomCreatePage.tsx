import { useNavigate, Link, useSearchParams } from "react-router-dom";
import { ChevronRight } from "lucide-react";
import { useCreateRoom } from "../hooks";
import { RoomForm } from "../components/RoomForm";
import type { CreateRoomFormData } from "../schema";

export function RoomCreatePage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const createMutation = useCreateRoom();

  const defaultValues: Partial<CreateRoomFormData> = {};
  const prefillPropertyId = searchParams.get("property_id");
  const prefillFloorId = searchParams.get("floor_id");
  if (prefillPropertyId) {
    defaultValues.property_id = prefillPropertyId;
  }
  if (prefillFloorId) {
    defaultValues.floor_id = prefillFloorId;
  }

  const handleSubmit = (data: CreateRoomFormData) => {
    createMutation.mutate(data, {
      onSuccess: () => navigate("/dashboard/rooms"),
    });
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-2 text-sm text-slate-500">
        <Link to="/dashboard" className="hover:text-slate-600">Dashboard</Link>
        <ChevronRight size={14} />
        <Link to="/dashboard/rooms" className="hover:text-slate-600">Rooms</Link>
        <ChevronRight size={14} />
        <span className="text-slate-600">New</span>
      </div>

      <h1 className="text-2xl font-bold">New Room</h1>

      <div className="max-w-2xl bg-white rounded-xl border border-slate-200 p-6">
        <RoomForm
          onSubmit={handleSubmit}
          defaultValues={defaultValues}
          isSubmitting={createMutation.isPending}
        />
      </div>
    </div>
  );
}

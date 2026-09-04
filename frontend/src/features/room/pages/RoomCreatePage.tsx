import { useNavigate, Link } from "react-router-dom";
import { ChevronRight } from "lucide-react";
import { useCreateRoom } from "../hooks";
import { RoomForm } from "../components/RoomForm";
import type { CreateRoomFormData } from "../schema";

export function RoomCreatePage() {
  const navigate = useNavigate();
  const createMutation = useCreateRoom();

  const handleSubmit = (data: CreateRoomFormData) => {
    createMutation.mutate(data, {
      onSuccess: () => navigate("/dashboard/rooms"),
    });
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-2 text-sm text-black/40">
        <Link to="/dashboard" className="hover:text-black/60">Dashboard</Link>
        <ChevronRight size={14} />
        <Link to="/dashboard/rooms" className="hover:text-black/60">Rooms</Link>
        <ChevronRight size={14} />
        <span className="text-black/60">New</span>
      </div>

      <h1 className="text-2xl font-bold">New Room</h1>

      <div className="max-w-2xl bg-white rounded-xl border border-black/5 p-6">
        <RoomForm
          onSubmit={handleSubmit}
          isSubmitting={createMutation.isPending}
        />
      </div>
    </div>
  );
}

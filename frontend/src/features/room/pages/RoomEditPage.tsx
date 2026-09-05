import { useNavigate, useParams, Link } from "react-router-dom";
import { ChevronRight } from "lucide-react";
import { useRoom, useUpdateRoom } from "../hooks";
import { RoomForm } from "../components/RoomForm";
import type { CreateRoomFormData } from "../schema";

export function RoomEditPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { data, isLoading } = useRoom(id);
  const updateMutation = useUpdateRoom();

  const handleSubmit = (formData: CreateRoomFormData) => {
    updateMutation.mutate(
      { id: id!, data: formData },
      { onSuccess: () => navigate(`/dashboard/rooms/${id}`) }
    );
  };

  if (isLoading) return <div className="text-center py-12 text-slate-400">Loading…</div>;
  if (!data) return <div className="text-center py-12 text-slate-400">Not found</div>;

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-2 text-sm text-slate-500">
        <Link to="/dashboard" className="hover:text-slate-600">Dashboard</Link>
        <ChevronRight size={14} />
        <Link to="/dashboard/rooms" className="hover:text-slate-600">Rooms</Link>
        <ChevronRight size={14} />
        <Link to={`/dashboard/rooms/${id}`} className="hover:text-slate-600 truncate max-w-32">{data.name}</Link>
        <ChevronRight size={14} />
        <span className="text-slate-600">Edit</span>
      </div>

      <h1 className="text-2xl font-bold">Edit Room</h1>

      <div className="max-w-2xl bg-white rounded-xl border border-slate-200 p-6">
        <RoomForm
          onSubmit={handleSubmit}
          defaultValues={data}
          isSubmitting={updateMutation.isPending}
        />
      </div>
    </div>
  );
}

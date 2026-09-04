import { useParams, useNavigate, Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { useRoom, useDeleteRoom } from "../hooks";
import { DoorOpen, ChevronRight, Pencil, Trash2, ArrowLeft, Users, DollarSign, Layers } from "lucide-react";

export function RoomDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { data, isLoading } = useRoom(id);
  const deleteMutation = useDeleteRoom();

  const handleDelete = () => {
    if (!confirm(`Delete room "${data?.name}"?`)) return;
    deleteMutation.mutate(id!, {
      onSuccess: () => navigate("/dashboard/rooms"),
    });
  };

  if (isLoading) return <div className="text-center py-12 text-black/30">Loading…</div>;
  if (!data) return <div className="text-center py-12 text-black/30">Not found</div>;

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-2 text-sm text-black/40">
        <Link to="/dashboard" className="hover:text-black/60">Dashboard</Link>
        <ChevronRight size={14} />
        <Link to="/dashboard/rooms" className="hover:text-black/60">Rooms</Link>
        <ChevronRight size={14} />
        <span className="text-black/60 truncate max-w-32">{data.name}</span>
      </div>

      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="w-12 h-12 rounded-xl bg-orange/10 flex items-center justify-center">
            <DoorOpen size={24} className="text-orange" />
          </div>
          <div>
            <h1 className="text-2xl font-bold">{data.name}</h1>
            <p className="text-sm text-black/40">
              <span className="inline-flex items-center gap-1">
                <Users size={14} /> {data.capacity} guests
              </span>
              <span className="mx-2">·</span>
              <span className="inline-flex items-center gap-1">
                <DollarSign size={14} /> {data.price.toLocaleString()}
              </span>
            </p>
          </div>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" onClick={() => navigate("/dashboard/rooms")}>
            <ArrowLeft size={16} className="mr-1" />
            Back
          </Button>
          <Button variant="outline" onClick={() => navigate(`/dashboard/rooms/${id}/edit`)}>
            <Pencil size={16} className="mr-1" />
            Edit
          </Button>
          <Button variant="destructive" onClick={handleDelete}>
            <Trash2 size={16} className="mr-1" />
            Delete
          </Button>
        </div>
      </div>

      <div className="bg-white rounded-xl border border-black/5 p-6">
        <dl className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div className="space-y-1">
            <dt className="text-xs font-semibold text-black/40 uppercase tracking-wider">Capacity</dt>
            <dd className="text-sm font-medium">{data.capacity} guests</dd>
          </div>
          <div className="space-y-1">
            <dt className="text-xs font-semibold text-black/40 uppercase tracking-wider">Price</dt>
            <dd className="text-sm font-medium">{data.price.toLocaleString()}</dd>
          </div>
          <div className="space-y-1">
            <dt className="text-xs font-semibold text-black/40 uppercase tracking-wider">Property ID</dt>
            <dd className="text-sm font-mono text-black/60">{data.property_id}</dd>
          </div>
          <div className="space-y-1">
            <dt className="text-xs font-semibold text-black/40 uppercase tracking-wider">Floor ID</dt>
            <dd className="text-sm font-mono text-black/60">{data.floor_id}</dd>
          </div>
          <div className="space-y-1">
            <dt className="text-xs font-semibold text-black/40 uppercase tracking-wider">Status</dt>
            <dd>
              <span className={`inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium ${
                data.is_available
                  ? "bg-green-100 text-green-700"
                  : "bg-red-100 text-red-700"
              }`}>
                {data.is_available ? "Available" : "Occupied"}
              </span>
            </dd>
          </div>
          <div className="space-y-1">
            <dt className="text-xs font-semibold text-black/40 uppercase tracking-wider">Room ID</dt>
            <dd className="text-sm font-mono text-black/60">{data.id}</dd>
          </div>
          <div className="space-y-1">
            <dt className="text-xs font-semibold text-black/40 uppercase tracking-wider">Created At</dt>
            <dd className="text-sm text-black/60">{new Date(data.created_at).toLocaleString()}</dd>
          </div>
          <div className="space-y-1">
            <dt className="text-xs font-semibold text-black/40 uppercase tracking-wider">Updated At</dt>
            <dd className="text-sm text-black/60">{new Date(data.updated_at).toLocaleString()}</dd>
          </div>
        </dl>
      </div>

      <div className="bg-white rounded-xl border border-black/5 p-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold flex items-center gap-2">
            <Layers size={18} className="text-orange" />
            Related Floor
          </h2>
          {data.floor_id && (
            <Button
              variant="outline"
              size="sm"
              onClick={() => navigate(`/dashboard/floors/${data.floor_id}`)}
            >
              View Floor
            </Button>
          )}
        </div>
        <p className="text-sm text-black/30 font-mono">{data.floor_id || "No floor assigned"}</p>
      </div>
    </div>
  );
}

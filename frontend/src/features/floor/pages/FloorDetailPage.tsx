import { useParams, useNavigate, Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { useFloor, useDeleteFloor } from "../hooks";
import { Layers, ChevronRight, Pencil, Trash2, ArrowLeft, DoorOpen } from "lucide-react";

export function FloorDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { data, isLoading } = useFloor(id);
  const deleteMutation = useDeleteFloor();

  const handleDelete = () => {
    if (!confirm(`Delete floor "${data?.name}"?`)) return;
    deleteMutation.mutate(id!, {
      onSuccess: () => navigate("/dashboard/floors"),
    });
  };

  if (isLoading) return <div className="text-center py-12 text-slate-400">Loading…</div>;
  if (!data) return <div className="text-center py-12 text-slate-400">Not found</div>;

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-2 text-sm text-slate-500">
        <Link to="/dashboard" className="hover:text-slate-600">Dashboard</Link>
        <ChevronRight size={14} />
        <Link to="/dashboard/floors" className="hover:text-slate-600">Floors</Link>
        <ChevronRight size={14} />
        <span className="text-slate-600 truncate max-w-32">{data.name}</span>
      </div>

      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="w-12 h-12 rounded-xl bg-orange/10 flex items-center justify-center">
            <Layers size={24} className="text-orange" />
          </div>
          <div>
            <h1 className="text-2xl font-bold">{data.name}</h1>
            <p className="text-sm text-slate-500">Floor {data.floor_number}</p>
          </div>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" onClick={() => navigate("/dashboard/floors")}>
            <ArrowLeft size={16} className="mr-1" />
            Back
          </Button>
          <Button variant="outline" onClick={() => navigate(`/dashboard/floors/${id}/edit`)}>
            <Pencil size={16} className="mr-1" />
            Edit
          </Button>
          <Button variant="destructive" onClick={handleDelete}>
            <Trash2 size={16} className="mr-1" />
            Delete
          </Button>
        </div>
      </div>

      <div className="bg-white rounded-xl border border-slate-200 p-6">
        <dl className="grid grid-cols-1 gap-5 sm:grid-cols-2 bg-white rounded-xl border border-slate-200 shadow-sm p-6">
          <div className="space-y-1">
            <dt className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Floor Number</dt>
            <dd className="text-sm font-medium">{data.floor_number}</dd>
          </div>
          <div className="space-y-1">
            <dt className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Building ID</dt>
            <dd className="text-sm font-mono text-slate-600">{data.building_id}</dd>
          </div>
          <div className="space-y-1">
            <dt className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Status</dt>
            <dd>
              <span className={`inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium ${
                data.is_active
                  ? "bg-green-50 text-green-700 border border-green-200"
                  : "bg-red-50 text-red-700 border border-red-200"
              }`}>
                {data.is_active ? "Active" : "Inactive"}
              </span>
            </dd>
          </div>
          <div className="space-y-1">
            <dt className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Floor ID</dt>
            <dd className="text-sm font-mono text-slate-600">{data.id}</dd>
          </div>
          <div className="space-y-1">
            <dt className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Created At</dt>
            <dd className="text-sm text-slate-600">{new Date(data.created_at).toLocaleString()}</dd>
          </div>
          <div className="space-y-1">
            <dt className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Updated At</dt>
            <dd className="text-sm text-slate-600">{new Date(data.updated_at).toLocaleString()}</dd>
          </div>
        </dl>
      </div>

      <div className="bg-white rounded-xl border border-slate-200 p-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold flex items-center gap-2">
            <DoorOpen size={18} className="text-orange" />
            Rooms on this Floor
          </h2>
          <Button
            variant="outline"
            size="sm"
            onClick={() => navigate(`/dashboard/rooms?floor_id=${data.id}`)}
          >
            View Rooms
          </Button>
        </div>
        <p className="text-sm text-slate-400">Rooms assigned to this floor will appear here.</p>
      </div>
    </div>
  );
}

import { useParams, useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { useBuilding, useDeleteBuilding } from "../hooks";

export function BuildingDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { data, isLoading } = useBuilding(id);
  const deleteMutation = useDeleteBuilding();

  const handleDelete = () => {
    if (!confirm("Are you sure you want to delete this building?")) return;
    deleteMutation.mutate(id!, {
      onSuccess: () => navigate("/dashboard/buildings"),
    });
  };

  if (isLoading) return <div className="text-center py-12 text-slate-400">Loading...</div>;
  if (!data) return <div className="text-center py-12 text-slate-400">Not found</div>;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">{data.name}</h1>
        <div className="flex gap-2">
          <Button variant="outline" onClick={() => navigate("/dashboard/buildings")}>
            Back
          </Button>
          <Button variant="outline" onClick={() => navigate(`/dashboard/buildings/${id}/edit`)}>
            Edit
          </Button>
          <Button variant="destructive" onClick={handleDelete}>
            Delete
          </Button>
        </div>
      </div>

      <dl className="grid grid-cols-1 gap-5 sm:grid-cols-2 bg-white rounded-xl border border-slate-200 shadow-sm p-6">
        <div className="space-y-1">
          <dt className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Name</dt>
          <dd className="text-sm text-slate-800 mt-1">{data.name}</dd>
        </div>
        <div className="space-y-1">
          <dt className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Total Floors</dt>
          <dd className="text-sm text-slate-800 mt-1">{data.total_floors}</dd>
        </div>
        <div className="space-y-1">
          <dt className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Property ID</dt>
          <dd className="text-sm text-slate-600">{data.property_id}</dd>
        </div>
        <div className="space-y-1">
          <dt className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Created At</dt>
          <dd className="text-sm text-slate-800 mt-1">{new Date(data.created_at).toLocaleString()}</dd>
        </div>
      </dl>

      <div className="border-t pt-4">
        <h2 className="text-lg font-semibold mb-2">Floors in this Building</h2>
        <p className="text-sm text-slate-400">Floors linked to this building will appear here.</p>
        <Button
          variant="outline"
          className="mt-3"
          onClick={() => navigate("/dashboard/floors")}
        >
          View Floors
        </Button>
      </div>
    </div>
  );
}

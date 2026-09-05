import { useParams, useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { useProperty, useDeleteProperty } from "../hooks";

export function PropertyDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { data, isLoading } = useProperty(id);
  const deleteMutation = useDeleteProperty();

  const handleDelete = () => {
    if (!confirm("Are you sure you want to delete this property?")) return;
    deleteMutation.mutate(id!, {
      onSuccess: () => navigate("/dashboard/properties"),
    });
  };

  if (isLoading) return <div className="text-center py-12 text-slate-400">Loading...</div>;
  if (!data) return <div className="text-center py-12 text-slate-400">Not found</div>;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">{data.name}</h1>
        <div className="flex gap-2">
          <Button variant="outline" onClick={() => navigate("/dashboard/properties")}>
            Back
          </Button>
          <Button variant="outline" onClick={() => navigate(`/dashboard/properties/${id}/edit`)}>
            Edit
          </Button>
          <Button variant="destructive" onClick={handleDelete}>
            Delete
          </Button>
        </div>
      </div>

      <dl className="grid grid-cols-1 gap-5 sm:grid-cols-2 bg-white rounded-xl border border-slate-200 shadow-sm p-6">
        <div className="space-y-1">
          <dt className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Type</dt>
          <dd className="text-sm capitalize">{data.property_type.replace(/_/g, " ")}</dd>
        </div>
        <div className="space-y-1">
          <dt className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Status</dt>
          <dd className="text-sm text-slate-800 mt-1">
            <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${
              data.is_active
                ? "bg-green-50 text-green-700 border border-green-200"
                : "bg-red-50 text-red-700 border border-red-200"
            }`}>
              {data.is_active ? "Active" : "Inactive"}
            </span>
          </dd>
        </div>
        <div className="space-y-1">
          <dt className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Address</dt>
          <dd className="text-sm text-slate-800 mt-1">{data.address || "—"}</dd>
        </div>
        <div className="space-y-1">
          <dt className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Description</dt>
          <dd className="text-sm text-slate-800 mt-1">{data.description || "—"}</dd>
        </div>
      </dl>
    </div>
  );
}

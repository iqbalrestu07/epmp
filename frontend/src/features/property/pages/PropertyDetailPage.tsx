import { useParams, useNavigate, Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { useProperty, useDeleteProperty } from "../hooks";
import { useBuildings } from "../../building/hooks";
import { Building2, Plus, ChevronRight, Pencil, Trash2, ArrowLeft } from "lucide-react";

export function PropertyDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { data, isLoading } = useProperty(id);
  const deleteMutation = useDeleteProperty();
  const { data: buildingsData } = useBuildings({ per_page: 100, property_id: id });
  const buildings = Array.isArray(buildingsData?.data)
    ? buildingsData.data
    : Array.isArray(buildingsData)
    ? buildingsData
    : [];

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
      <div className="flex items-center gap-2 text-sm text-slate-500">
        <Link to="/dashboard/properties" className="hover:text-slate-900">Properties</Link>
        <ChevronRight size={14} />
        <span className="text-slate-900 font-medium truncate max-w-32">{data.name}</span>
      </div>

      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="w-12 h-12 rounded-xl bg-orange/10 flex items-center justify-center">
            <Building2 size={24} className="text-orange" />
          </div>
          <div>
            <h1 className="text-2xl font-bold text-slate-900">{data.name}</h1>
            <p className="text-sm text-slate-500 capitalize">{data.property_type.replace(/_/g, " ")}</p>
          </div>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" onClick={() => navigate("/dashboard/properties")}>
            <ArrowLeft size={16} className="mr-1" />
            Back
          </Button>
          <Button variant="outline" onClick={() => navigate(`/dashboard/properties/${id}/edit`)}>
            <Pencil size={16} className="mr-1" />
            Edit
          </Button>
          <Button variant="destructive" onClick={handleDelete}>
            <Trash2 size={16} className="mr-1" />
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

      <div className="bg-white rounded-xl border border-slate-200 shadow-sm p-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold flex items-center gap-2">
            <Building2 size={18} className="text-orange" />
            Buildings in this Property
            <span className="text-sm font-normal text-slate-400">({buildings.length})</span>
          </h2>
          <div className="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => navigate(`/dashboard/buildings?property_id=${data.id}`)}
            >
              View All Buildings
            </Button>
            <Button
              size="sm"
              onClick={() => navigate(`/dashboard/buildings/new?property_id=${data.id}`)}
            >
              <Plus size={14} className="mr-1" />
              New Building
            </Button>
          </div>
        </div>
        {buildings.length === 0 ? (
          <p className="text-sm text-slate-400">No buildings created yet for this property.</p>
        ) : (
          <div className="space-y-2">
            {buildings.map((b) => (
              <div
                key={b.id}
                onClick={() => navigate(`/dashboard/buildings/${b.id}`)}
                className="flex items-center justify-between p-3 rounded-lg border border-slate-200 hover:border-slate-300 hover:bg-slate-50 cursor-pointer transition-all"
              >
                <div className="flex items-center gap-2">
                  <span className="text-lg">🏢</span>
                  <div>
                    <p className="font-semibold text-sm text-slate-900">{b.name}</p>
                    <p className="text-xs text-slate-400">{b.total_floors} floors</p>
                  </div>
                </div>
                <ChevronRight size={16} className="text-slate-400" />
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

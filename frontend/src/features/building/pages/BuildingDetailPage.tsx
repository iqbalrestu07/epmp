import { useParams, useNavigate, Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { useBuilding, useDeleteBuilding } from "../hooks";
import { usePropertys } from "../../property/hooks";
import { useFloors } from "../../floor/hooks";
import { Building2, Plus, ChevronRight, Pencil, Trash2, ArrowLeft, Layers } from "lucide-react";

export function BuildingDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { data, isLoading } = useBuilding(id);
  const deleteMutation = useDeleteBuilding();
  const { data: propertiesData } = usePropertys({ per_page: 100 });
  const properties = Array.isArray(propertiesData?.data)
    ? propertiesData.data
    : Array.isArray(propertiesData)
    ? propertiesData
    : [];
  const parentProperty = properties.find((p) => p.id === data?.property_id);

  const { data: floorsData } = useFloors({ per_page: 100, building_id: id });
  const floors = Array.isArray(floorsData?.data)
    ? floorsData.data
    : Array.isArray(floorsData)
    ? floorsData
    : [];

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
      <div className="flex items-center gap-2 text-sm text-slate-500">
        <Link to="/dashboard/properties" className="hover:text-slate-900">Properties</Link>
        <ChevronRight size={14} />
        {parentProperty ? (
          <>
            <Link to={`/dashboard/properties/${parentProperty.id}`} className="hover:text-slate-900">{parentProperty.name}</Link>
            <ChevronRight size={14} />
          </>
        ) : null}
        <Link to="/dashboard/buildings" className="hover:text-slate-900">Buildings</Link>
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
            <p className="text-sm text-slate-500">{data.total_floors} floors</p>
          </div>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" onClick={() => navigate("/dashboard/buildings")}>
            <ArrowLeft size={16} className="mr-1" />
            Back
          </Button>
          <Button variant="outline" onClick={() => navigate(`/dashboard/buildings/${id}/edit`)}>
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
          <dt className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Name</dt>
          <dd className="text-sm text-slate-800 mt-1">{data.name}</dd>
        </div>
        <div className="space-y-1">
          <dt className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Total Floors</dt>
          <dd className="text-sm text-slate-800 mt-1">{data.total_floors}</dd>
        </div>
        <div className="space-y-1">
          <dt className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Parent Property</dt>
          <dd className="text-sm text-slate-800 mt-1">
            {parentProperty ? (
              <Link to={`/dashboard/properties/${parentProperty.id}`} className="text-orange hover:underline">
                {parentProperty.name}
              </Link>
            ) : (
              <span className="text-slate-400">—</span>
            )}
          </dd>
        </div>
        <div className="space-y-1">
          <dt className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Created At</dt>
          <dd className="text-sm text-slate-800 mt-1">{new Date(data.created_at).toLocaleString()}</dd>
        </div>
      </dl>

      <div className="bg-white rounded-xl border border-slate-200 shadow-sm p-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold flex items-center gap-2">
            <Layers size={18} className="text-orange" />
            Floors in this Building
            <span className="text-sm font-normal text-slate-400">({floors.length})</span>
          </h2>
          <div className="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => navigate(`/dashboard/floors?building_id=${data.id}`)}
            >
              View All Floors
            </Button>
            <Button
              size="sm"
              onClick={() => navigate(`/dashboard/floors/new?building_id=${data.id}`)}
            >
              <Plus size={14} className="mr-1" />
              New Floor
            </Button>
          </div>
        </div>
        {floors.length === 0 ? (
          <p className="text-sm text-slate-400">No floors created yet for this building.</p>
        ) : (
          <div className="space-y-2">
            {[...floors]
              .sort((a, b) => a.floor_number - b.floor_number)
              .map((f) => (
                <div
                  key={f.id}
                  onClick={() => navigate(`/dashboard/floors/${f.id}`)}
                  className="flex items-center justify-between p-3 rounded-lg border border-slate-200 hover:border-slate-300 hover:bg-slate-50 cursor-pointer transition-all"
                >
                  <div className="flex items-center gap-2">
                    <Layers size={16} className="text-slate-400" />
                    <div>
                      <p className="font-semibold text-sm text-slate-900">{f.name}</p>
                      <p className="text-xs text-slate-400">Level {f.floor_number}</p>
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

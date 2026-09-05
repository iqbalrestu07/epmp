import { useParams, useNavigate, Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { useRoom, useDeleteRoom } from "../hooks";
import { useFloors } from "../../floor/hooks";
import { useBuildings } from "../../building/hooks";
import { usePropertys } from "../../property/hooks";
import { DoorOpen, ChevronRight, Pencil, Trash2, ArrowLeft, Users, DollarSign, Layers, Building2 } from "lucide-react";

export function RoomDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { data, isLoading } = useRoom(id);
  const deleteMutation = useDeleteRoom();

  const { data: floorsData } = useFloors({ per_page: 100 });
  const floors = Array.isArray(floorsData?.data)
    ? floorsData.data
    : Array.isArray(floorsData)
    ? floorsData
    : [];
  const parentFloor = floors.find((f) => f.id === data?.floor_id);

  const { data: buildingsData } = useBuildings({ per_page: 100 });
  const buildings = Array.isArray(buildingsData?.data)
    ? buildingsData.data
    : Array.isArray(buildingsData)
    ? buildingsData
    : [];
  const parentBuilding = parentFloor
    ? buildings.find((b) => b.id === parentFloor.building_id)
    : undefined;

  const { data: propertiesData } = usePropertys({ per_page: 100 });
  const properties = Array.isArray(propertiesData?.data)
    ? propertiesData.data
    : Array.isArray(propertiesData)
    ? propertiesData
    : [];
  const parentProperty = properties.find((p) => p.id === data?.property_id);

  const handleDelete = () => {
    if (!confirm(`Delete room "${data?.name}"?`)) return;
    deleteMutation.mutate(id!, {
      onSuccess: () => navigate("/dashboard/rooms"),
    });
  };

  if (isLoading) return <div className="text-center py-12 text-slate-400">Loading…</div>;
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
        {parentBuilding ? (
          <>
            <Link to={`/dashboard/buildings/${parentBuilding.id}`} className="hover:text-slate-900">{parentBuilding.name}</Link>
            <ChevronRight size={14} />
          </>
        ) : null}
        {parentFloor ? (
          <>
            <Link to={`/dashboard/floors/${parentFloor.id}`} className="hover:text-slate-900">{parentFloor.name}</Link>
            <ChevronRight size={14} />
          </>
        ) : null}
        <Link to="/dashboard/rooms" className="hover:text-slate-900">Rooms</Link>
        <ChevronRight size={14} />
        <span className="text-slate-900 font-medium truncate max-w-32">{data.name}</span>
      </div>

      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="w-12 h-12 rounded-xl bg-orange/10 flex items-center justify-center">
            <DoorOpen size={24} className="text-orange" />
          </div>
          <div>
            <h1 className="text-2xl font-bold text-slate-900">{data.name}</h1>
            <p className="text-sm text-slate-500">
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

      <dl className="grid grid-cols-1 gap-5 sm:grid-cols-2 bg-white rounded-xl border border-slate-200 shadow-sm p-6">
        <div className="space-y-1">
          <dt className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Capacity</dt>
          <dd className="text-sm font-medium">{data.capacity} guests</dd>
        </div>
        <div className="space-y-1">
          <dt className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Price</dt>
          <dd className="text-sm font-medium">{data.price.toLocaleString()}</dd>
        </div>
        <div className="space-y-1">
          <dt className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Status</dt>
          <dd>
            <span className={`inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium ${
              data.is_available
                ? "bg-green-50 text-green-700 border border-green-200"
                : "bg-red-50 text-red-700 border border-red-200"
            }`}>
              {data.is_available ? "Available" : "Occupied"}
            </span>
          </dd>
        </div>
        <div className="space-y-1">
          <dt className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Created At</dt>
          <dd className="text-sm text-slate-600">{new Date(data.created_at).toLocaleString()}</dd>
        </div>
      </dl>

      <div className="bg-white rounded-xl border border-slate-200 shadow-sm p-6">
        <h2 className="text-lg font-semibold mb-4">Parent Hierarchy</h2>
        <div className="space-y-3">
          <div className="flex items-center justify-between p-3 rounded-lg border border-slate-200 hover:bg-slate-50 transition-all">
            <div className="flex items-center gap-2">
              <Building2 size={16} className="text-slate-400" />
              <div>
                <p className="text-xs text-slate-500">Property</p>
                <p className="font-semibold text-sm text-slate-900">
                  {parentProperty ? parentProperty.name : "—"}
                </p>
              </div>
            </div>
            {parentProperty && (
              <Button variant="outline" size="sm" onClick={() => navigate(`/dashboard/properties/${parentProperty.id}`)}>
                View Property
              </Button>
            )}
          </div>
          <div className="flex items-center justify-between p-3 rounded-lg border border-slate-200 hover:bg-slate-50 transition-all">
            <div className="flex items-center gap-2">
              <Building2 size={16} className="text-slate-400" />
              <div>
                <p className="text-xs text-slate-500">Building</p>
                <p className="font-semibold text-sm text-slate-900">
                  {parentBuilding ? parentBuilding.name : "—"}
                </p>
              </div>
            </div>
            {parentBuilding && (
              <Button variant="outline" size="sm" onClick={() => navigate(`/dashboard/buildings/${parentBuilding.id}`)}>
                View Building
              </Button>
            )}
          </div>
          <div className="flex items-center justify-between p-3 rounded-lg border border-slate-200 hover:bg-slate-50 transition-all">
            <div className="flex items-center gap-2">
              <Layers size={16} className="text-slate-400" />
              <div>
                <p className="text-xs text-slate-500">Floor</p>
                <p className="font-semibold text-sm text-slate-900">
                  {parentFloor ? `${parentFloor.name} (Level ${parentFloor.floor_number})` : "—"}
                </p>
              </div>
            </div>
            {parentFloor && (
              <Button variant="outline" size="sm" onClick={() => navigate(`/dashboard/floors/${parentFloor.id}`)}>
                View Floor
              </Button>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

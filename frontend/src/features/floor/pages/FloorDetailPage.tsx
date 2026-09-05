import { useParams, useNavigate, Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { useFloor, useDeleteFloor } from "../hooks";
import { useBuildings } from "../../building/hooks";
import { usePropertys } from "../../property/hooks";
import { useRooms } from "../../room/hooks";
import { Layers, ChevronRight, Pencil, Trash2, ArrowLeft, DoorOpen, Plus } from "lucide-react";

export function FloorDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { data, isLoading } = useFloor(id);
  const deleteMutation = useDeleteFloor();

  const { data: buildingsData } = useBuildings({ per_page: 100 });
  const buildings = Array.isArray(buildingsData?.data)
    ? buildingsData.data
    : Array.isArray(buildingsData)
    ? buildingsData
    : [];
  const parentBuilding = buildings.find((b) => b.id === data?.building_id);

  const { data: propertiesData } = usePropertys({ per_page: 100 });
  const properties = Array.isArray(propertiesData?.data)
    ? propertiesData.data
    : Array.isArray(propertiesData)
    ? propertiesData
    : [];
  const parentProperty = parentBuilding
    ? properties.find((p) => p.id === parentBuilding.property_id)
    : undefined;

  const { data: roomsData } = useRooms({ per_page: 100, floor_id: id });
  const rooms = Array.isArray(roomsData?.data)
    ? roomsData.data
    : Array.isArray(roomsData)
    ? roomsData
    : [];

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
        <Link to="/dashboard/floors" className="hover:text-slate-900">Floors</Link>
        <ChevronRight size={14} />
        <span className="text-slate-900 font-medium truncate max-w-32">{data.name}</span>
      </div>

      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="w-12 h-12 rounded-xl bg-orange/10 flex items-center justify-center">
            <Layers size={24} className="text-orange" />
          </div>
          <div>
            <h1 className="text-2xl font-bold text-slate-900">{data.name}</h1>
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

      <dl className="grid grid-cols-1 gap-5 sm:grid-cols-2 bg-white rounded-xl border border-slate-200 shadow-sm p-6">
        <div className="space-y-1">
          <dt className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Floor Number</dt>
          <dd className="text-sm font-medium">{data.floor_number}</dd>
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
          <dt className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Parent Building</dt>
          <dd className="text-sm text-slate-800 mt-1">
            {parentBuilding ? (
              <Link to={`/dashboard/buildings/${parentBuilding.id}`} className="text-orange hover:underline">
                {parentBuilding.name}
              </Link>
            ) : (
              <span className="text-slate-400">—</span>
            )}
          </dd>
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
          <dd className="text-sm text-slate-600">{new Date(data.created_at).toLocaleString()}</dd>
        </div>
        <div className="space-y-1">
          <dt className="text-xs font-semibold text-slate-500 uppercase tracking-wider">Updated At</dt>
          <dd className="text-sm text-slate-600">{new Date(data.updated_at).toLocaleString()}</dd>
        </div>
      </dl>

      <div className="bg-white rounded-xl border border-slate-200 shadow-sm p-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-semibold flex items-center gap-2">
            <DoorOpen size={18} className="text-orange" />
            Rooms on this Floor
            <span className="text-sm font-normal text-slate-400">({rooms.length})</span>
          </h2>
          <div className="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => navigate(`/dashboard/rooms?floor_id=${data.id}`)}
            >
              View All Rooms
            </Button>
            <Button
              size="sm"
              onClick={() => navigate(`/dashboard/rooms/new?floor_id=${data.id}${parentProperty ? `&property_id=${parentProperty.id}` : ''}`)}
            >
              <Plus size={14} className="mr-1" />
              New Room
            </Button>
          </div>
        </div>
        {rooms.length === 0 ? (
          <p className="text-sm text-slate-400">No rooms created yet on this floor.</p>
        ) : (
          <div className="space-y-2">
            {rooms.map((r) => (
              <div
                key={r.id}
                onClick={() => navigate(`/dashboard/rooms/${r.id}`)}
                className="flex items-center justify-between p-3 rounded-lg border border-slate-200 hover:border-slate-300 hover:bg-slate-50 cursor-pointer transition-all"
              >
                <div className="flex items-center gap-2">
                  <DoorOpen size={16} className="text-slate-400" />
                  <div>
                    <p className="font-semibold text-sm text-slate-900">{r.name}</p>
                    <p className="text-xs text-slate-400">Cap: {r.capacity} · Rp {r.price?.toLocaleString()}</p>
                  </div>
                </div>
                <span className={`px-2 py-0.5 rounded-full text-[10px] font-semibold ${
                  r.is_available
                    ? "bg-green-50 text-green-700 border border-green-200"
                    : "bg-red-50 text-red-700 border border-red-200"
                }`}>
                  {r.is_available ? "Available" : "Occupied"}
                </span>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

import { useState } from "react";
import { useNavigate, useSearchParams, Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { DoorOpen, Plus } from "lucide-react";
import { RoomTable } from "../components/RoomTable";
import { useRooms, useDeleteRoom } from "../hooks";
import { useFloors } from "../../floor/hooks";
import { useBuildings } from "../../building/hooks";
import { usePropertys } from "../../property/hooks";
import type { Room } from "../types";

export function RoomListPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const propertyIdFromUrl = searchParams.get("property_id") ?? "";
  const buildingIdFromUrl = searchParams.get("building_id") ?? "";
  const floorIdFromUrl = searchParams.get("floor_id") ?? "";

  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [searchInput, setSearchInput] = useState("");
  const [propertyId, setPropertyId] = useState(propertyIdFromUrl);
  const [buildingId, setBuildingId] = useState(buildingIdFromUrl);
  const [floorId, setFloorId] = useState(floorIdFromUrl);

  const { data, isLoading } = useRooms({
    page,
    per_page: 20,
    search: search || undefined,
    floor_id: floorId || undefined,
    property_id: propertyId || undefined,
    building_id: buildingId || undefined,
  });

  const { data: propertiesData } = usePropertys({ per_page: 100 });
  const properties = Array.isArray(propertiesData?.data)
    ? propertiesData.data
    : Array.isArray(propertiesData)
    ? propertiesData
    : [];

  const { data: buildingsData } = useBuildings({
    per_page: 100,
    property_id: propertyId || undefined,
  });
  const allBuildings = Array.isArray(buildingsData?.data)
    ? buildingsData.data
    : Array.isArray(buildingsData)
    ? buildingsData
    : [];

  const { data: floorsData } = useFloors({
    per_page: 100,
    building_id: buildingId || undefined,
  });
  const allFloors = Array.isArray(floorsData?.data)
    ? floorsData.data
    : Array.isArray(floorsData)
    ? floorsData
    : [];

  const deleteMutation = useDeleteRoom();

  const handleRowClick = (row: Room) => {
    navigate(`/dashboard/rooms/${row.id}`);
  };

  const handleEdit = (row: Room) => {
    navigate(`/dashboard/rooms/${row.id}/edit`);
  };

  const handleDelete = (row: Room) => {
    if (!confirm(`Delete room "${row.name}"?`)) return;
    deleteMutation.mutate(row.id);
  };

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setSearch(searchInput);
    setPage(1);
  };

  const handlePropertyChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    setPropertyId(e.target.value);
    setBuildingId("");
    setFloorId("");
    setPage(1);
  };

  const handleBuildingChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    setBuildingId(e.target.value);
    setFloorId("");
    setPage(1);
  };

  const handleFloorChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    setFloorId(e.target.value);
    setPage(1);
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <div className="flex items-center gap-2 text-sm text-slate-500 mb-1">
            <DoorOpen size={14} />
            <Link to="/dashboard/properties" className="hover:text-slate-900">Properties</Link>
            <span>/</span>
            <span className="text-slate-900 font-medium">Rooms & Units</span>
          </div>
          <h1 className="text-2xl font-bold text-slate-900">Rooms & Units</h1>
        </div>
        <Button onClick={() => navigate("/dashboard/rooms/new")}>
          <Plus className="w-4 h-4 mr-1" />
          New Room
        </Button>
      </div>

      <form onSubmit={handleSearch} className="flex flex-wrap gap-2">
        <Input
          placeholder="Search rooms..."
          value={searchInput}
          onChange={(e) => setSearchInput(e.target.value)}
          className="max-w-sm"
        />
        <Button type="submit" variant="outline">Search</Button>
        <select
          value={propertyId}
          onChange={handlePropertyChange}
          className="sm:max-w-xs rounded-lg border border-slate-300 bg-white px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-orange/30 text-slate-800"
        >
          <option value="">All Properties</option>
          {properties.map((p) => (
            <option key={p.id} value={p.id}>{p.name}</option>
          ))}
        </select>
        <select
          value={buildingId}
          onChange={handleBuildingChange}
          disabled={!propertyId}
          className="sm:max-w-xs rounded-lg border border-slate-300 bg-white px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-orange/30 text-slate-800 disabled:opacity-50"
        >
          <option value="">All Buildings</option>
          {allBuildings.map((b) => (
            <option key={b.id} value={b.id}>{b.name}</option>
          ))}
        </select>
        <select
          value={floorId}
          onChange={handleFloorChange}
          disabled={!buildingId}
          className="sm:max-w-xs rounded-lg border border-slate-300 bg-white px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-orange/30 text-slate-800 disabled:opacity-50"
        >
          <option value="">All Floors</option>
          {allFloors.map((f) => (
            <option key={f.id} value={f.id}>{f.name}</option>
          ))}
        </select>
      </form>

      {isLoading ? (
        <div className="text-center py-12 text-slate-400">Loading rooms...</div>
      ) : (
        <>
          <RoomTable
            data={data?.data ?? []}
            floors={allFloors}
            buildings={allBuildings}
            properties={properties}
            onRowClick={handleRowClick}
            onEdit={handleEdit}
            onDelete={handleDelete}
          />

          {data && data.total_pages > 1 && (
            <div className="flex items-center justify-between">
              <p className="text-sm text-slate-500">
                Page {data.page} of {data.total_pages} ({data.total} total)
              </p>
              <div className="flex gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  disabled={page <= 1}
                  onClick={() => setPage(page - 1)}
                >
                  Previous
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  disabled={page >= data.total_pages}
                  onClick={() => setPage(page + 1)}
                >
                  Next
                </Button>
              </div>
            </div>
          )}
        </>
      )}
    </div>
  );
}

import { useState } from "react";
import { useNavigate, useSearchParams, Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Layers, Plus } from "lucide-react";
import { FloorTable } from "../components/FloorTable";
import { useFloors, useDeleteFloor } from "../hooks";
import { useBuildings } from "../../building/hooks";
import type { Floor } from "../types";

export function FloorListPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const buildingIdFromUrl = searchParams.get("building_id") ?? "";

  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [searchInput, setSearchInput] = useState("");
  const [buildingId, setBuildingId] = useState(buildingIdFromUrl);

  const { data, isLoading } = useFloors({
    page,
    per_page: 20,
    search: search || undefined,
    building_id: buildingId || undefined,
  });
  const { data: buildingsData } = useBuildings({ per_page: 100 });
  const buildings = Array.isArray(buildingsData?.data)
    ? buildingsData.data
    : Array.isArray(buildingsData)
    ? buildingsData
    : [];

  const deleteMutation = useDeleteFloor();

  const handleRowClick = (row: Floor) => {
    navigate(`/dashboard/floors/${row.id}`);
  };

  const handleEdit = (row: Floor) => {
    navigate(`/dashboard/floors/${row.id}/edit`);
  };

  const handleDelete = (row: Floor) => {
    if (!confirm(`Delete floor "${row.name}"?`)) return;
    deleteMutation.mutate(row.id);
  };

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setSearch(searchInput);
    setPage(1);
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <div className="flex items-center gap-2 text-sm text-slate-500 mb-1">
            <Layers size={14} />
            <Link to="/dashboard/buildings" className="hover:text-slate-900">Buildings</Link>
            <span>/</span>
            <span className="text-slate-900 font-medium">Floors</span>
          </div>
          <h1 className="text-2xl font-bold text-slate-900">Floors</h1>
        </div>
        <Button onClick={() => navigate("/dashboard/floors/new")}>
          <Plus className="w-4 h-4 mr-1" />
          New Floor
        </Button>
      </div>

      <form onSubmit={handleSearch} className="flex gap-2">
        <Input
          placeholder="Search floors..."
          value={searchInput}
          onChange={(e) => setSearchInput(e.target.value)}
          className="max-w-sm"
        />
        <Button type="submit" variant="outline">Search</Button>
        <select
          value={buildingId}
          onChange={(e) => { setBuildingId(e.target.value); setPage(1); }}
          className="sm:max-w-xs rounded-lg border border-slate-300 bg-white px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-orange/30 text-slate-800"
        >
          <option value="">All Buildings</option>
          {buildings.map((b) => (
            <option key={b.id} value={b.id}>{b.name}</option>
          ))}
        </select>
      </form>

      {isLoading ? (
        <div className="text-center py-12 text-slate-400">Loading floors...</div>
      ) : (
        <>
          <FloorTable
            data={data?.data ?? []}
            buildings={buildings}
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

import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Layers, Plus, Search, ChevronRight } from "lucide-react";
import { FloorTable } from "../components/FloorTable";
import { useFloors, useDeleteFloor } from "../hooks";
import { useBuildings } from "../../building/hooks";
import type { Floor } from "../types";

export function FloorListPage() {
  const navigate = useNavigate();
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [searchInput, setSearchInput] = useState("");
  const [buildingId, setBuildingId] = useState("");

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
      <div className="space-y-2">
        <div className="flex items-center gap-2 text-sm text-black/40">
          <Link to="/dashboard" className="hover:text-black/60">Dashboard</Link>
          <ChevronRight size={14} />
          <span className="text-black/60">Floors</span>
        </div>
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-xl bg-orange/10 flex items-center justify-center">
              <Layers size={20} className="text-orange" />
            </div>
            <div>
              <h1 className="text-2xl font-bold">Floors</h1>
              <p className="text-sm text-black/40">Manage floors across your buildings</p>
            </div>
          </div>
          <Button onClick={() => navigate("/dashboard/floors/new")}>
            <Plus size={16} className="mr-1" />
            New Floor
          </Button>
        </div>
      </div>

      <div className="flex flex-col sm:flex-row gap-3">
        <form onSubmit={handleSearch} className="flex gap-2 flex-1">
          <div className="relative flex-1">
            <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-black/30" />
            <Input
              placeholder="Search floors by name…"
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
              className="pl-9"
            />
          </div>
          <Button type="submit" variant="outline">Search</Button>
        </form>
        <select
          value={buildingId}
          onChange={(e) => { setBuildingId(e.target.value); setPage(1); }}
          className="sm:max-w-xs rounded-lg border border-black/10 bg-white px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-orange/30"
        >
          <option value="">All Buildings</option>
          {buildings.map((b) => (
            <option key={b.id} value={b.id}>{b.name}</option>
          ))}
        </select>
      </div>

      {isLoading ? (
        <div className="text-center py-12 text-black/30">Loading floors…</div>
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
              <p className="text-sm text-black/40">
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

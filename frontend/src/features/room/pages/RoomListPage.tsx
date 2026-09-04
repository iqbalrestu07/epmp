import { useState } from "react";
import { useNavigate, useSearchParams, Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { DoorOpen, Plus, Search, ChevronRight } from "lucide-react";
import { RoomTable } from "../components/RoomTable";
import { useRooms, useDeleteRoom } from "../hooks";
import type { Room } from "../types";

export function RoomListPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const floorIdFromUrl = searchParams.get("floor_id") ?? "";

  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [searchInput, setSearchInput] = useState("");
  const [floorId, setFloorId] = useState(floorIdFromUrl);

  const { data, isLoading } = useRooms({
    page,
    per_page: 20,
    search: search || undefined,
    floor_id: floorId || undefined,
  });

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

  return (
    <div className="space-y-6">
      <div className="space-y-2">
        <div className="flex items-center gap-2 text-sm text-black/40">
          <Link to="/dashboard" className="hover:text-black/60">Dashboard</Link>
          <ChevronRight size={14} />
          <span className="text-black/60">Rooms</span>
        </div>
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-xl bg-orange/10 flex items-center justify-center">
              <DoorOpen size={20} className="text-orange" />
            </div>
            <div>
              <h1 className="text-2xl font-bold">Rooms & Units</h1>
              <p className="text-sm text-black/40">Manage rooms and units across your properties</p>
            </div>
          </div>
          <Button onClick={() => navigate("/dashboard/rooms/new")}>
            <Plus size={16} className="mr-1" />
            New Room
          </Button>
        </div>
      </div>

      <div className="flex flex-col sm:flex-row gap-3">
        <form onSubmit={handleSearch} className="flex gap-2 flex-1">
          <div className="relative flex-1">
            <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-black/30" />
            <Input
              placeholder="Search rooms by name…"
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
              className="pl-9"
            />
          </div>
          <Button type="submit" variant="outline">Search</Button>
        </form>
        <Input
          placeholder="Filter by floor ID…"
          value={floorId}
          onChange={(e) => { setFloorId(e.target.value); setPage(1); }}
          className="sm:max-w-xs"
        />
      </div>

      {isLoading ? (
        <div className="text-center py-12 text-black/30">Loading rooms…</div>
      ) : (
        <>
          <RoomTable
            data={data?.data ?? []}
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

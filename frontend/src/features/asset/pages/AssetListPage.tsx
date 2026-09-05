import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { AssetTable } from "../components/AssetTable";
import { useAssets, useDeleteAsset } from "../hooks";
import type { Asset } from "../types";
import { ChevronRight } from "lucide-react";

export function AssetListPage() {
  const navigate = useNavigate();
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [searchInput, setSearchInput] = useState("");

  const { data, isLoading } = useAssets({
    page,
    per_page: 20,
    search: search || undefined,
  });

  const deleteMutation = useDeleteAsset();

  const handleRowClick = (row: Asset) => {
    navigate(`/dashboard/assets/${row.id}`);
  };
  
  const handleEditClick = (row: Asset, e: React.MouseEvent) => {
    e.stopPropagation();
    navigate(`/dashboard/assets/${row.id}/edit`);
  };

  const handleDeleteClick = (row: Asset, e: React.MouseEvent) => {
    e.stopPropagation();
    if (confirm("Are you sure you want to delete this asset?")) {
      deleteMutation.mutate(row.id);
    }
  };

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setSearch(searchInput);
    setPage(1);
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-2 text-sm text-slate-500">
        <Link to="/dashboard" className="hover:text-slate-600">Dashboard</Link>
        <ChevronRight size={14} />
        <span className="text-slate-600">Assets</span>
      </div>

      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <h1 className="text-2xl font-bold">Assets</h1>
        <Button onClick={() => navigate("/dashboard/assets/new")} className="bg-orange hover:bg-orange/90 text-white">
          + Add Asset
        </Button>
      </div>

      <div className="bg-white rounded-xl shadow-sm border border-slate-200 overflow-hidden">
        <div className="p-4 border-b border-slate-200">
          <form onSubmit={handleSearch} className="flex gap-2 max-w-md">
            <Input
              placeholder="Search assets..."
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
              className="flex-1 rounded-lg border-slate-300 focus:ring-orange/30"
            />
            <Button type="submit" variant="outline" className="rounded-lg">Search</Button>
          </form>
        </div>

        {isLoading ? (
          <div className="text-center py-12 text-slate-500 text-sm">Loading assets...</div>
        ) : (
          <AssetTable 
            data={data?.data ?? []} 
            onRowClick={handleRowClick} 
            onEditClick={handleEditClick}
            onDeleteClick={handleDeleteClick}
          />
        )}
      </div>

      {data && data.total_pages > 1 && (
        <div className="flex items-center justify-between pt-4">
          <p className="text-sm text-slate-500">
            Showing page {data.page} of {data.total_pages} ({data.total} total)
          </p>
          <div className="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              disabled={page <= 1}
              onClick={() => setPage(page - 1)}
              className="rounded-lg"
            >
              Previous
            </Button>
            <Button
              variant="outline"
              size="sm"
              disabled={page >= data.total_pages}
              onClick={() => setPage(page + 1)}
              className="rounded-lg"
            >
              Next
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}

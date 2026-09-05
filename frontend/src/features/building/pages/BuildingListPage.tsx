import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { BuildingTable } from "../components/BuildingTable";
import { useBuildings } from "../hooks";
import { usePropertys } from "../../property/hooks";
import type { Building } from "../types";
import type { Property } from "../../property/types";
import { Plus, Building2 } from "lucide-react";

export function BuildingListPage() {
  const navigate = useNavigate();
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [searchInput, setSearchInput] = useState("");
  const [propertyId, setPropertyId] = useState("");

  const { data, isLoading } = useBuildings({
    page,
    per_page: 20,
    search: search || undefined,
    property_id: propertyId || undefined,
  });
  const { data: propertiesData } = usePropertys({ per_page: 100 });
  const properties = propertiesData?.data ?? [];

  const handleRowClick = (row: Building) => {
    navigate(`/dashboard/buildings/${row.id}`);
  };
  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setSearch(searchInput);
    setPage(1);
  };

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <div className="flex items-center gap-2 text-sm text-gray-500 mb-1">
            <Building2 size={14} />
            <button
              onClick={() => navigate("/dashboard/properties")}
              className="hover:text-black"
            >
              Properties
            </button>
            <span>/</span>
            <span className="text-black font-medium">Buildings</span>
          </div>
          <h1 className="text-2xl font-bold">Buildings</h1>
        </div>
        <Button onClick={() => navigate("/dashboard/buildings/new")}>
          <Plus className="w-4 h-4 mr-1" />
          New Building
        </Button>
      </div>
      <form onSubmit={handleSearch} className="flex gap-2">
        <Input
          placeholder="Search..."
          value={searchInput}
          onChange={(e) => setSearchInput(e.target.value)}
          className="max-w-sm"
        />
        <Button type="submit" variant="outline">Search</Button>
        <select
          value={propertyId}
          onChange={(e) => { setPropertyId(e.target.value); setPage(1); }}
          className="sm:max-w-xs rounded-md border border-gray-300 px-3 py-2 text-sm"
        >
          <option value="">All Properties</option>
          {properties.map((p: Property) => (
            <option key={p.id} value={p.id}>{p.name}</option>
          ))}
        </select>
      </form>

      {isLoading ? (
        <div className="text-center py-8">Loading...</div>
      ) : (
        <>
          <BuildingTable data={data?.data ?? []} properties={properties} onRowClick={handleRowClick} />

          {data && data.total_pages > 1 && (
            <div className="flex items-center justify-between">
              <p className="text-sm text-gray-500">
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

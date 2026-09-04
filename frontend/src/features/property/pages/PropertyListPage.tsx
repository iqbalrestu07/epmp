import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { PropertyTable } from "../components/PropertyTable";
import { usePropertys } from "../hooks";
import type { Property } from "../types";
import { Eye, Plus, Globe2 } from "lucide-react";

export function PropertyListPage() {
  const navigate = useNavigate();
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [searchInput, setSearchInput] = useState("");

  const { data, isLoading } = usePropertys({
    page,
    per_page: 20,
    search: search || undefined,
  });

  const handleRowClick = (row: Property) => {
    navigate(`/dashboard/properties/${row.id}`);
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
            <Globe2 size={14} />
            <button
              onClick={() => navigate("/dashboard/organizations")}
              className="hover:text-black"
            >
              Organizations
            </button>
            <span>/</span>
            <span className="text-black font-medium">Properties</span>
          </div>
          <div className="flex items-center gap-3">
            <h1 className="text-2xl font-bold">Properties</h1>
            <Button
              variant="outline"
              size="sm"
              onClick={() => navigate("/dashboard/properties/interactive")}
            >
              <Eye className="w-4 h-4 mr-1" />
              Interactive 3D
            </Button>
          </div>
        </div>
        <Button onClick={() => navigate("/dashboard/properties/new")}>
          <Plus className="w-4 h-4 mr-1" />
          New Property
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
      </form>

      {isLoading ? (
        <div className="text-center py-8">Loading...</div>
      ) : (
        <>
          <PropertyTable data={data?.data ?? []} onRowClick={handleRowClick} />

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

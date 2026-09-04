import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { OrganizationTable } from "../components/OrganizationTable";
import { useOrganizations } from "../hooks";
import type { Organization } from "../types";
import { Plus } from "lucide-react";

export function OrganizationListPage() {
  const navigate = useNavigate();
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [searchInput, setSearchInput] = useState("");

  const { data, isLoading } = useOrganizations({
    page,
    per_page: 20,
    search: search || undefined,
  });

  const handleRowClick = (row: Organization) => {
    navigate(`/dashboard/organizations/${row.id}`);
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
          <h1 className="text-2xl font-bold">Organizations</h1>
          <p className="text-sm text-gray-500 mt-1">Create and manage your organizations first before adding properties.</p>
        </div>
        <Button onClick={() => navigate("/dashboard/organizations/new")}>
          <Plus className="w-4 h-4 mr-1" />
          New Organization
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
          <OrganizationTable data={data?.data ?? []} onRowClick={handleRowClick} />

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

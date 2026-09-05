import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { OrganizationTable } from "../components/OrganizationTable";
import { useOrg } from "../context/OrgContext";
import type { Organization } from "../types";
import { Plus, ChevronRight } from "lucide-react";

export function OrganizationListPage() {
  const navigate = useNavigate();
  const { orgs, isLoading, refreshOrgs } = useOrg();
  const [searchInput, setSearchInput] = useState("");

  const handleRowClick = (row: Organization) => {
    navigate(`/dashboard/organizations/${row.id}`);
  };

  const filtered = searchInput
    ? orgs.filter(o =>
        o.name.toLowerCase().includes(searchInput.toLowerCase()) ||
        o.domain?.toLowerCase().includes(searchInput.toLowerCase())
      )
    : orgs;

  return (
    <div className="space-y-4">
      <div className="flex items-center gap-2 text-sm text-slate-500">
        <Link to="/dashboard" className="hover:text-slate-600">Dashboard</Link>
        <ChevronRight size={14} />
        <span className="text-slate-600">Organizations</span>
      </div>
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold">Organizations</h1>
          <p className="text-sm text-slate-500 mt-1">Your organizations. Switch between them using the header dropdown.</p>
        </div>
        <Button onClick={() => navigate("/dashboard/organizations/new")}>
          <Plus className="w-4 h-4 mr-1" />
          New Organization
        </Button>
      </div>
      <form onSubmit={(e) => { e.preventDefault(); }} className="flex gap-2">
        <Input
          placeholder="Search..."
          value={searchInput}
          onChange={(e) => setSearchInput(e.target.value)}
          className="max-w-sm"
        />
        <Button type="button" variant="outline" onClick={() => refreshOrgs()}>Refresh</Button>
      </form>

      {isLoading ? (
        <div className="text-center py-12 text-slate-400">Loading...</div>
      ) : filtered.length === 0 ? (
        <div className="text-center py-12 bg-white rounded-xl border">
          <p className="text-slate-500 mb-4">No organizations found.</p>
          <Button onClick={() => navigate("/dashboard/organizations/new")}>
            <Plus className="w-4 h-4 mr-1" />
            Create your first organization
          </Button>
        </div>
      ) : (
        <OrganizationTable data={filtered} onRowClick={handleRowClick} />
      )}
    </div>
  );
}

import { useNavigate, Link } from "react-router-dom";
import { useCreateOrganization } from "../hooks";
import { OrganizationForm } from "../components/OrganizationForm";
import type { CreateOrganizationFormData } from "../schema";
import { useOrg } from "../context/OrgContext";
import { ChevronRight, Sparkles, Building2 } from "lucide-react";

export function OrganizationCreatePage() {
  const navigate = useNavigate();
  const createMutation = useCreateOrganization();
  const { orgs, refreshOrgs } = useOrg();
  const isFirstOrg = orgs.length === 0;

  const handleSubmit = (data: CreateOrganizationFormData) => {
    createMutation.mutate(data, {
      onSuccess: async () => {
        await refreshOrgs();
        navigate("/dashboard");
      },
    });
  };

  return (
    <div className="space-y-6 max-w-3xl mx-auto py-2">
      {/* Breadcrumbs */}
      <div className="flex items-center gap-2 text-sm text-slate-400">
        <Link to="/dashboard" className="hover:text-slate-600 transition-colors">Dashboard</Link>
        <ChevronRight size={14} />
        <Link to="/dashboard/organizations" className="hover:text-slate-600 transition-colors">Organizations</Link>
        <ChevronRight size={14} />
        <span className="text-slate-700 font-medium">New Organization</span>
      </div>

      {/* Header Banner */}
      <div className="bg-white p-6 sm:p-8 rounded-2xl border border-slate-200 shadow-sm relative overflow-hidden">
        <div className="absolute -right-8 -bottom-8 opacity-5 pointer-events-none">
          <Building2 size={200} />
        </div>
        <div className="flex items-center gap-2 text-xs font-bold uppercase tracking-wider text-orange mb-2">
          <Sparkles size={14} />
          {isFirstOrg ? "Welcome to EPMP Platform" : "Workspace Management"}
        </div>
        <h1 className="text-2xl sm:text-3xl font-bold text-slate-900 tracking-tight">
          {isFirstOrg ? "Create your primary workspace" : "Add a new Organization"}
        </h1>
        <p className="text-sm text-slate-500 mt-1 max-w-xl">
          {isFirstOrg
            ? "An organization connects your team, properties, rental units, financial contracts, and 3D digital twins under one roof."
            : "Set up an isolated workspace for a separate property group, subsidiary, or client portfolio."}
        </p>
      </div>

      {/* Form Card */}
      <div className="bg-white p-6 sm:p-8 rounded-2xl border border-slate-200 shadow-sm">
        <OrganizationForm
          onSubmit={handleSubmit}
          isSubmitting={createMutation.isPending}
        />
      </div>
    </div>
  );
}

import { useParams, useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { useOrganization, useDeleteOrganization } from "../hooks";
import { usePropertys } from "../../property/hooks";
import { useOrg } from "../context/OrgContext";
import { Building2, Eye, Plus, Check, ArrowRight, MapPin } from "lucide-react";

export function OrganizationDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { data, isLoading } = useOrganization(id);
  const deleteMutation = useDeleteOrganization();
  const { currentOrg, switchOrg } = useOrg();

  // Fetch properties belonging to this specific organization
  const { data: propertyData, isLoading: isLoadingProperties } = usePropertys({
    organization_id: id,
    per_page: 100,
  });
  const properties = propertyData?.data ?? [];

  const isCurrentActiveOrg = currentOrg?.id === id;

  const handleDelete = () => {
    if (!confirm("Are you sure you want to delete this organization?")) return;
    deleteMutation.mutate(id!, {
      onSuccess: () => navigate("/dashboard/organizations"),
    });
  };

  if (isLoading) return <div className="p-8 text-center text-slate-500">Loading organization details...</div>;
  if (!data) return <div className="p-8 text-center text-slate-500">Organization not found</div>;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white p-6 rounded-2xl border shadow-sm">
        <div>
          <div className="flex items-center gap-3">
            <h1 className="text-2xl font-bold text-slate-900">{data.name}</h1>
            <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold ${
              data.is_active
                ? "bg-green-50 text-green-700 border border-green-200"
                : "bg-red-50 text-red-700 border border-red-200"
            }`}>
              {data.is_active ? "Active" : "Inactive"}
            </span>
            {isCurrentActiveOrg && (
              <span className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-semibold bg-orange/10 text-orange border border-orange/20">
                <Check size={12} /> Active Context
              </span>
            )}
          </div>
          <p className="text-sm text-slate-500 mt-1">Domain: <span className="font-mono text-slate-700">{data.domain || "—"}</span></p>
        </div>

        <div className="flex items-center gap-2 flex-wrap">
          {!isCurrentActiveOrg && (
            <Button
              variant="outline"
              size="sm"
              onClick={() => {
                switchOrg(data.id);
              }}
              className="text-orange hover:bg-orange/5 border-orange/30"
            >
              Set as Active
            </Button>
          )}
          <Button variant="outline" size="sm" onClick={() => navigate("/dashboard/organizations")}>
            Back
          </Button>
          <Button variant="outline" size="sm" onClick={() => navigate(`/dashboard/organizations/${id}/edit`)}>
            Edit
          </Button>
          <Button variant="destructive" size="sm" onClick={handleDelete}>
            Delete
          </Button>
        </div>
      </div>

      {/* Properties in this Organization Section */}
      <div className="bg-white p-6 rounded-2xl border shadow-sm space-y-4">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b">
          <div>
            <div className="flex items-center gap-2">
              <Building2 className="text-orange" size={20} />
              <h2 className="text-lg font-bold text-slate-900">Properties in this Organization</h2>
              <span className="px-2 py-0.5 bg-slate-100 text-slate-700 rounded-full text-xs font-semibold">
                {properties.length}
              </span>
            </div>
            <p className="text-xs text-slate-500 mt-0.5">
              List of all properties and estates managed under {data.name}.
            </p>
          </div>

          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => navigate("/dashboard/properties/interactive")}
              className="flex items-center gap-1.5"
            >
              <Eye size={14} className="text-orange" />
              Interactive 3D
            </Button>
            <Button
              size="sm"
              onClick={() => navigate("/dashboard/properties/new")}
              className="flex items-center gap-1.5"
            >
              <Plus size={14} />
              New Property
            </Button>
          </div>
        </div>

        {/* Property List / Grid */}
        {isLoadingProperties ? (
          <div className="py-8 text-center text-slate-400 text-sm">Loading properties...</div>
        ) : properties.length === 0 ? (
          <div className="py-12 text-center rounded-xl border border-dashed border-slate-200">
            <Building2 className="mx-auto text-gray-300 mb-2" size={36} />
            <p className="text-slate-700 font-medium">No properties found in this organization</p>
            <p className="text-xs text-slate-400 mt-1 mb-4">Create your first property to start adding buildings, floors, and rooms.</p>
            <Button
              size="sm"
              onClick={() => navigate("/dashboard/properties/new")}
            >
              <Plus size={14} className="mr-1" />
              Create First Property
            </Button>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {properties.map((property) => (
              <div
                key={property.id}
                className="group p-4 rounded-xl border border-slate-200 hover:border-orange/50 hover:shadow-md transition-all bg-white flex flex-col justify-between"
              >
                <div>
                  <div className="flex items-start justify-between gap-2">
                    <h3 className="font-semibold text-slate-900 group-hover:text-orange transition-colors">
                      {property.name}
                    </h3>
                    <span className="text-[11px] font-medium px-2 py-0.5 rounded-full bg-slate-100 text-slate-700 capitalize shrink-0">
                      {property.property_type.replace("_", " ")}
                    </span>
                  </div>
                  {property.address && (
                    <p className="text-xs text-slate-500 mt-1 flex items-center gap-1">
                      <MapPin size={12} className="text-slate-400 shrink-0" />
                      <span className="truncate">{property.address}</span>
                    </p>
                  )}
                  {property.description && (
                    <p className="text-xs text-slate-400 mt-2 line-clamp-2">
                      {property.description}
                    </p>
                  )}
                </div>

                <div className="mt-4 pt-3 border-t border-gray-100 flex items-center justify-between text-xs">
                  <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-medium ${
                    property.is_active ? "bg-green-50 text-green-700" : "bg-red-50 text-red-600"
                  }`}>
                    {property.is_active ? "Active" : "Inactive"}
                  </span>
                  <div className="flex items-center gap-2">
                    <button
                      onClick={() => navigate(`/dashboard/properties/${property.id}`)}
                      className="font-medium text-slate-600 hover:text-slate-900 flex items-center gap-0.5"
                    >
                      Detail <ArrowRight size={12} />
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

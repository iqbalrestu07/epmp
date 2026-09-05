import { useState } from "react";
import { useNavigate, Link, useSearchParams } from "react-router-dom";
import { useCreateBuilding } from "../hooks";
import { BuildingForm } from "../components/BuildingForm";
import { InteractiveBuildingCreator } from "../components/InteractiveBuildingCreator";
import type { CreateBuildingFormData } from "../schema";
import { ChevronRight, Building2, Box, FileText } from "lucide-react";

export function BuildingCreatePage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const createMutation = useCreateBuilding();
  const [viewMode, setViewMode] = useState<"interactive" | "basic">("interactive");

  const defaultValues: Partial<CreateBuildingFormData> = {};
  const prefillPropertyId = searchParams.get("property_id");
  if (prefillPropertyId) {
    defaultValues.property_id = prefillPropertyId;
  }

  const handleSubmit = (data: CreateBuildingFormData) => {
    createMutation.mutate(data, {
      onSuccess: () => navigate("/dashboard/buildings"),
    });
  };

  return (
    <div className="space-y-6">
      {/* Breadcrumb Navigation */}
      <div className="flex items-center gap-2 text-sm text-slate-500">
        <Link to="/dashboard" className="hover:text-slate-800 transition-colors">Dashboard</Link>
        <ChevronRight size={14} />
        <Link to="/dashboard/buildings" className="hover:text-slate-800 transition-colors">Buildings</Link>
        <ChevronRight size={14} />
        <span className="text-slate-800 font-medium">New Building</span>
      </div>

      {/* Page Title & View Switcher */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-900 flex items-center gap-2">
            <Building2 className="text-orange" size={26} />
            Create Building
          </h1>
          <p className="text-sm text-slate-500 mt-0.5">
            {viewMode === "interactive"
              ? "Position the building in a real 3D environment using the digital twin."
              : "Standard form mode for quick data entry."}
          </p>
        </div>

        {/* View Toggle Tabs */}
        <div className="flex bg-slate-200/80 p-1 rounded-xl shadow-inner border border-slate-300/60 self-start">
          <button
            type="button"
            onClick={() => setViewMode("interactive")}
            className={`flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-semibold transition-all ${
              viewMode === "interactive"
                ? "bg-white text-slate-900 shadow-sm"
                : "text-slate-600 hover:text-slate-900"
            }`}
          >
            <Box size={15} className={viewMode === "interactive" ? "text-orange" : ""} />
            Interactive 3D View
          </button>
          <button
            type="button"
            onClick={() => setViewMode("basic")}
            className={`flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-semibold transition-all ${
              viewMode === "basic"
                ? "bg-white text-slate-900 shadow-sm"
                : "text-slate-600 hover:text-slate-900"
            }`}
          >
            <FileText size={15} className={viewMode === "basic" ? "text-orange" : ""} />
            Basic View (Form)
          </button>
        </div>
      </div>

      {/* Content Canvas or Form */}
      <div className="flex-1">
        {viewMode === "basic" ? (
          <div className="max-w-2xl bg-white p-8 rounded-2xl border border-slate-200 shadow-sm">
            <h2 className="text-lg font-bold text-slate-800 mb-2">Building Information</h2>
            <p className="text-sm text-slate-500 mb-6">Enter the property and floor capacity for this building structure.</p>
            <BuildingForm
              onSubmit={handleSubmit}
              defaultValues={defaultValues}
              isSubmitting={createMutation.isPending}
            />
          </div>
        ) : (
          <InteractiveBuildingCreator
            onSubmit={handleSubmit}
            isSubmitting={createMutation.isPending}
          />
        )}
      </div>
    </div>
  );
}

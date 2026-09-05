import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { ChevronRight } from "lucide-react";
import { useCreateAssetAssignment } from "../hooks";
import { AssetAssignmentForm } from "../components/AssetAssignmentForm";
import { InteractiveAssetAssignment } from "../components/InteractiveAssetAssignment";
import { Button } from "@/components/ui/button";

export function AssetAssignmentCreatePage() {
  const navigate = useNavigate();
  const createMutation = useCreateAssetAssignment();
  const [viewMode, setViewMode] = useState<"basic" | "interactive">("interactive");

  const handleSubmit = (data: any) => {
    // The data might come from Interactive as just { asset_id, room_id } 
    // We append the current date for assigned_date if needed by schema.
    const submissionData = {
      ...data,
      assigned_date: data.assigned_date || new Date().toISOString(),
    };

    createMutation.mutate(submissionData, {
      onSuccess: () => navigate("/dashboard/asset-assignments"),
    });
  };

  return (
    <div className="space-y-6 h-full flex flex-col">
      <div className="flex items-center gap-2 text-sm text-black/40">
        <Link to="/dashboard" className="hover:text-black/60">Dashboard</Link>
        <ChevronRight size={14} />
        <Link to="/dashboard/asset-assignments" className="hover:text-black/60">Asset Assignments</Link>
        <ChevronRight size={14} />
        <span className="text-black/60">New Assignment</span>
      </div>

      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">New Assignment</h1>
        <div className="flex gap-2 bg-slate-100 p-1 rounded-md">
          <Button
            variant={viewMode === "basic" ? "default" : "ghost"}
            size="sm"
            onClick={() => setViewMode("basic")}
            className="text-sm"
          >
            Basic View
          </Button>
          <Button
            variant={viewMode === "interactive" ? "default" : "ghost"}
            size="sm"
            onClick={() => setViewMode("interactive")}
            className="text-sm"
          >
            Interactive Planner
          </Button>
        </div>
      </div>
      
      <div className="flex-1 min-h-[600px]">
        {viewMode === "basic" ? (
          <div className="max-w-2xl bg-white p-6 rounded-xl border border-black/5 shadow-sm">
            <AssetAssignmentForm
              onSubmit={handleSubmit as any}
              isSubmitting={createMutation.isPending}
            />
          </div>
        ) : (
          <InteractiveAssetAssignment
            onSubmit={handleSubmit}
            isSubmitting={createMutation.isPending}
          />
        )}
      </div>
    </div>
  );
}

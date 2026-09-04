import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useCreateBuilding } from "../hooks";
import { BuildingForm } from "../components/BuildingForm";
import { InteractiveBuildingCreator } from "../components/InteractiveBuildingCreator";
import type { CreateBuildingFormData } from "../schema";
import { Button } from "@/components/ui/button";

export function BuildingCreatePage() {
  const navigate = useNavigate();
  const createMutation = useCreateBuilding();
  const [viewMode, setViewMode] = useState<"basic" | "interactive">("interactive");

  const handleSubmit = (data: CreateBuildingFormData) => {
    createMutation.mutate(data, {
      onSuccess: () => navigate("/dashboard/buildings"),
    });
  };

  return (
    <div className="space-y-4 h-full flex flex-col">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">New Building</h1>
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
            Interactive View
          </Button>
        </div>
      </div>
      
      <div className="flex-1 min-h-[600px]">
        {viewMode === "basic" ? (
          <div className="bg-white p-6 rounded-xl border border-slate-200">
            <BuildingForm
              onSubmit={handleSubmit}
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

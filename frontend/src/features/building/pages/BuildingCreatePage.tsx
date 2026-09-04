import { useNavigate } from "react-router-dom";
import { useCreateBuilding } from "../hooks";
import { BuildingForm } from "../components/BuildingForm";
import type { CreateBuildingFormData } from "../schema";

export function BuildingCreatePage() {
  const navigate = useNavigate();
  const createMutation = useCreateBuilding();

  const handleSubmit = (data: CreateBuildingFormData) => {
    createMutation.mutate(data, {
      onSuccess: () => navigate("/dashboard/buildings"),
    });
  };

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold">New Building</h1>
      <BuildingForm
        onSubmit={handleSubmit}
        isSubmitting={createMutation.isPending}
      />
    </div>
  );
}

import { useNavigate, Link, useSearchParams } from "react-router-dom";
import { ChevronRight } from "lucide-react";
import { useCreateFloor } from "../hooks";
import { FloorForm } from "../components/FloorForm";
import type { CreateFloorFormData } from "../schema";

export function FloorCreatePage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const createMutation = useCreateFloor();

  const defaultValues: Partial<CreateFloorFormData> = {};
  const prefillBuildingId = searchParams.get("building_id");
  if (prefillBuildingId) {
    defaultValues.building_id = prefillBuildingId;
  }

  const handleSubmit = (data: CreateFloorFormData) => {
    createMutation.mutate(data, {
      onSuccess: () => navigate("/dashboard/floors"),
    });
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-2 text-sm text-slate-500">
        <Link to="/dashboard" className="hover:text-slate-600">Dashboard</Link>
        <ChevronRight size={14} />
        <Link to="/dashboard/floors" className="hover:text-slate-600">Floors</Link>
        <ChevronRight size={14} />
        <span className="text-slate-600">New</span>
      </div>

      <h1 className="text-2xl font-bold">New Floor</h1>

      <div className="max-w-2xl bg-white rounded-xl border border-slate-200 p-6">
        <FloorForm
          onSubmit={handleSubmit}
          defaultValues={defaultValues}
          isSubmitting={createMutation.isPending}
        />
      </div>
    </div>
  );
}

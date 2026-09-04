import { useNavigate, Link } from "react-router-dom";
import { ChevronRight } from "lucide-react";
import { useCreateFloor } from "../hooks";
import { FloorForm } from "../components/FloorForm";
import type { CreateFloorFormData } from "../schema";

export function FloorCreatePage() {
  const navigate = useNavigate();
  const createMutation = useCreateFloor();

  const handleSubmit = (data: CreateFloorFormData) => {
    createMutation.mutate(data, {
      onSuccess: () => navigate("/dashboard/floors"),
    });
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-2 text-sm text-black/40">
        <Link to="/dashboard" className="hover:text-black/60">Dashboard</Link>
        <ChevronRight size={14} />
        <Link to="/dashboard/floors" className="hover:text-black/60">Floors</Link>
        <ChevronRight size={14} />
        <span className="text-black/60">New</span>
      </div>

      <h1 className="text-2xl font-bold">New Floor</h1>

      <div className="max-w-2xl bg-white rounded-xl border border-black/5 p-6">
        <FloorForm
          onSubmit={handleSubmit}
          isSubmitting={createMutation.isPending}
        />
      </div>
    </div>
  );
}

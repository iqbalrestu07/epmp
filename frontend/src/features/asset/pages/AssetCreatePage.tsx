import { useNavigate, Link } from "react-router-dom";
import { ChevronRight } from "lucide-react";
import { useCreateAsset } from "../hooks";
import { AssetForm } from "../components/AssetForm";
import type { CreateAssetFormData } from "../schema";

export function AssetCreatePage() {
  const navigate = useNavigate();
  const createMutation = useCreateAsset();

  const handleSubmit = (data: CreateAssetFormData) => {
    createMutation.mutate(data, {
      onSuccess: () => navigate("/dashboard/assets"),
    });
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-2 text-sm text-slate-500">
        <Link to="/dashboard" className="hover:text-slate-600">Dashboard</Link>
        <ChevronRight size={14} />
        <Link to="/dashboard/assets" className="hover:text-slate-600">Assets</Link>
        <ChevronRight size={14} />
        <span className="text-slate-600">New</span>
      </div>

      <h1 className="text-2xl font-bold">New Asset</h1>
      
      <div className="max-w-2xl bg-white rounded-xl border border-slate-200 p-6">
        <AssetForm
          onSubmit={handleSubmit}
          isSubmitting={createMutation.isPending}
        />
      </div>
    </div>
  );
}

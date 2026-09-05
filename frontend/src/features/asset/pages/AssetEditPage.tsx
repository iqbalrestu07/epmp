import { useParams, useNavigate, Link } from "react-router-dom";
import { ChevronRight } from "lucide-react";
import { useAsset, useUpdateAsset } from "../hooks";
import { AssetForm } from "../components/AssetForm";
import type { UpdateAssetFormData } from "../schema";

export function AssetEditPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { data: asset, isLoading } = useAsset(id!);
  const updateMutation = useUpdateAsset();

  const handleSubmit = (data: UpdateAssetFormData) => {
    updateMutation.mutate({ id: id!, data }, {
      onSuccess: () => navigate("/dashboard/assets"),
    });
  };

  if (isLoading) return <div className="p-8 text-center text-slate-500">Loading asset...</div>;
  if (!asset) return <div className="p-8 text-center text-red-500">Asset not found.</div>;

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-2 text-sm text-slate-500">
        <Link to="/dashboard" className="hover:text-slate-600">Dashboard</Link>
        <ChevronRight size={14} />
        <Link to="/dashboard/assets" className="hover:text-slate-600">Assets</Link>
        <ChevronRight size={14} />
        <span className="text-slate-600">Edit Asset</span>
      </div>

      <h1 className="text-2xl font-bold">Edit Asset: {asset.name}</h1>
      
      <div className="max-w-2xl bg-white rounded-xl border border-slate-200 p-6">
        <AssetForm
          onSubmit={handleSubmit as any}
          defaultValues={asset as any}
          isSubmitting={updateMutation.isPending}
        />
      </div>
    </div>
  );
}

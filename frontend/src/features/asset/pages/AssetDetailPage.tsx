import { useParams, useNavigate, Link } from "react-router-dom";
import { ChevronRight, Edit2, Tag } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useAsset } from "../hooks";

export function AssetDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { data: asset, isLoading } = useAsset(id!);

  if (isLoading) return <div className="p-8 text-center text-slate-500">Loading asset...</div>;
  if (!asset) return <div className="p-8 text-center text-red-500">Asset not found.</div>;

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-2 text-sm text-slate-500">
        <Link to="/dashboard" className="hover:text-slate-600">Dashboard</Link>
        <ChevronRight size={14} />
        <Link to="/dashboard/assets" className="hover:text-slate-600">Assets</Link>
        <ChevronRight size={14} />
        <span className="text-slate-600">{asset.name}</span>
      </div>

      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold flex items-center gap-2">
          <Tag className="text-orange" /> {asset.name}
        </h1>
        <Button onClick={() => navigate(`/dashboard/assets/${id}/edit`)} variant="outline">
          <Edit2 size={16} className="mr-2" />
          Edit Asset
        </Button>
      </div>

      <div className="bg-white rounded-xl shadow-sm border border-slate-200 p-6 space-y-6">
        <div className="grid grid-cols-2 md:grid-cols-3 gap-6">
          <div>
            <h3 className="text-sm font-medium text-slate-500 mb-1">Status</h3>
            <p className="font-medium">{asset.status}</p>
          </div>
          <div>
            <h3 className="text-sm font-medium text-slate-500 mb-1">Category</h3>
            <p className="font-medium">{asset.category}</p>
          </div>
          <div>
            <h3 className="text-sm font-medium text-slate-500 mb-1">Purchase Price</h3>
            <p className="font-medium">${asset.purchase_price.toLocaleString()}</p>
          </div>
        </div>
      </div>
    </div>
  );
}

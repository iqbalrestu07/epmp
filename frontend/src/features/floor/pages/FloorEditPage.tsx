import { useNavigate, useParams, Link } from "react-router-dom";
import { ChevronRight } from "lucide-react";
import { useFloor, useUpdateFloor } from "../hooks";
import { FloorForm } from "../components/FloorForm";
import type { CreateFloorFormData } from "../schema";

export function FloorEditPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { data, isLoading } = useFloor(id);
  const updateMutation = useUpdateFloor();

  const handleSubmit = (formData: CreateFloorFormData) => {
    updateMutation.mutate(
      { id: id!, data: formData },
      { onSuccess: () => navigate(`/dashboard/floors/${id}`) }
    );
  };

  if (isLoading) return <div className="text-center py-12 text-slate-400">Loading…</div>;
  if (!data) return <div className="text-center py-12 text-slate-400">Not found</div>;

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-2 text-sm text-slate-500">
        <Link to="/dashboard" className="hover:text-slate-600">Dashboard</Link>
        <ChevronRight size={14} />
        <Link to="/dashboard/floors" className="hover:text-slate-600">Floors</Link>
        <ChevronRight size={14} />
        <Link to={`/dashboard/floors/${id}`} className="hover:text-slate-600 truncate max-w-32">{data.name}</Link>
        <ChevronRight size={14} />
        <span className="text-slate-600">Edit</span>
      </div>

      <h1 className="text-2xl font-bold">Edit Floor</h1>

      <div className="max-w-2xl bg-white rounded-xl border border-slate-200 p-6">
        <FloorForm
          onSubmit={handleSubmit}
          defaultValues={data}
          isSubmitting={updateMutation.isPending}
        />
      </div>
    </div>
  );
}

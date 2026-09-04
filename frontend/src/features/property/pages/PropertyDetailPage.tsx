import { useParams, useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { useProperty, useDeleteProperty } from "../hooks";

export function PropertyDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { data, isLoading } = useProperty(id);
  const deleteMutation = useDeleteProperty();

  const handleDelete = () => {
    if (!confirm("Are you sure you want to delete this property?")) return;
    deleteMutation.mutate(id!, {
      onSuccess: () => navigate("/dashboard/properties"),
    });
  };

  if (isLoading) return <div>Loading...</div>;
  if (!data) return <div>Not found</div>;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">{data.name}</h1>
        <div className="flex gap-2">
          <Button variant="outline" onClick={() => navigate("/dashboard/properties")}>
            Back
          </Button>
          <Button variant="outline" onClick={() => navigate(`/dashboard/properties/${id}/edit`)}>
            Edit
          </Button>
          <Button variant="destructive" onClick={handleDelete}>
            Delete
          </Button>
        </div>
      </div>

      <dl className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <div className="space-y-1">
          <dt className="text-sm font-medium text-gray-500">Type</dt>
          <dd className="text-sm capitalize">{data.property_type.replace(/_/g, " ")}</dd>
        </div>
        <div className="space-y-1">
          <dt className="text-sm font-medium text-gray-500">Status</dt>
          <dd className="text-sm">
            <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${
              data.is_active
                ? "bg-green-100 text-green-700"
                : "bg-red-100 text-red-700"
            }`}>
              {data.is_active ? "Active" : "Inactive"}
            </span>
          </dd>
        </div>
        <div className="space-y-1">
          <dt className="text-sm font-medium text-gray-500">Address</dt>
          <dd className="text-sm">{data.address || "—"}</dd>
        </div>
        <div className="space-y-1">
          <dt className="text-sm font-medium text-gray-500">Description</dt>
          <dd className="text-sm">{data.description || "—"}</dd>
        </div>
      </dl>
    </div>
  );
}

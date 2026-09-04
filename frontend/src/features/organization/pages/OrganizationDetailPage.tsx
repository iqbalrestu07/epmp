import { useParams, useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { useOrganization, useDeleteOrganization } from "../hooks";

export function OrganizationDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { data, isLoading } = useOrganization(id);
  const deleteMutation = useDeleteOrganization();

  const handleDelete = () => {
    if (!confirm("Are you sure you want to delete this organization?")) return;
    deleteMutation.mutate(id!, {
      onSuccess: () => navigate("/dashboard/organizations"),
    });
  };

  if (isLoading) return <div>Loading...</div>;
  if (!data) return <div>Not found</div>;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">{data.name}</h1>
        <div className="flex gap-2">
          <Button variant="outline" onClick={() => navigate("/dashboard/organizations")}>
            Back
          </Button>
          <Button variant="outline" onClick={() => navigate(`/dashboard/organizations/${id}/edit`)}>
            Edit
          </Button>
          <Button variant="destructive" onClick={handleDelete}>
            Delete
          </Button>
        </div>
      </div>

      <dl className="grid grid-cols-1 gap-4 sm:grid-cols-2">
        <div className="space-y-1">
          <dt className="text-sm font-medium text-gray-500">Domain</dt>
          <dd className="text-sm">{data.domain || "—"}</dd>
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
      </dl>

      <div className="border-t pt-4">
        <h2 className="text-lg font-semibold mb-2">Properties in this Organization</h2>
        <p className="text-sm text-gray-400">Properties linked to this organization will appear here.</p>
        <Button
          variant="outline"
          className="mt-3"
          onClick={() => navigate("/dashboard/properties")}
        >
          View Properties
        </Button>
      </div>
    </div>
  );
}

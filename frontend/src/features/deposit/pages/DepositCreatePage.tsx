import { useNavigate, useSearchParams } from "react-router-dom";
import { useCreateDeposit } from "../hooks";
import { DepositForm } from "../components/DepositForm";
import type { CreateDepositFormData } from "../schema";

export function DepositCreatePage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const createMutation = useCreateDeposit();

  const prefilledValues: Partial<CreateDepositFormData> = {
    contract_id: searchParams.get("contract_id") || undefined,
    tenant_id: searchParams.get("tenant_id") || undefined,
    amount: searchParams.get("amount") ? Number(searchParams.get("amount")) : undefined,
  };

  const handleSubmit = (data: CreateDepositFormData) => {
    createMutation.mutate(data, {
      onSuccess: () => navigate("/dashboard/deposits"),
    });
  };

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold">Catat Uang Jaminan (Deposit)</h1>
      <DepositForm
        onSubmit={handleSubmit}
        defaultValues={prefilledValues}
        isSubmitting={createMutation.isPending}
      />
    </div>
  );
}

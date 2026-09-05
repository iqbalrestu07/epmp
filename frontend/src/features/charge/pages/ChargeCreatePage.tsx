import { useNavigate, useSearchParams } from "react-router-dom";
import { useCreateCharge } from "../hooks";
import { ChargeForm } from "../components/ChargeForm";
import type { CreateChargeFormData } from "../schema";

export function ChargeCreatePage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const createMutation = useCreateCharge();

  const prefilledValues: Partial<CreateChargeFormData> = {
    contract_id: searchParams.get("contract_id") || undefined,
    invoice_id: searchParams.get("invoice_id") || undefined,
    amount: searchParams.get("amount") ? Number(searchParams.get("amount")) : undefined,
  };

  const handleSubmit = (data: CreateChargeFormData) => {
    createMutation.mutate(data, {
      onSuccess: () => navigate("/dashboard/charges"),
    });
  };

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold">Tambah Biaya Tambahan (Charge)</h1>
      <ChargeForm
        onSubmit={handleSubmit}
        defaultValues={prefilledValues}
        isSubmitting={createMutation.isPending}
      />
    </div>
  );
}

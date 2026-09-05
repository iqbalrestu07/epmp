import { useNavigate, useSearchParams } from "react-router-dom";
import { useCreatePayment } from "../hooks";
import { PaymentForm } from "../components/PaymentForm";
import type { CreatePaymentFormData } from "../schema";

export function PaymentCreatePage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const createMutation = useCreatePayment();

  const prefilledValues: Partial<CreatePaymentFormData> = {
    invoice_id: searchParams.get("invoice_id") || undefined,
    tenant_id: searchParams.get("tenant_id") || undefined,
    amount: searchParams.get("amount") ? Number(searchParams.get("amount")) : undefined,
  };

  const handleSubmit = (data: CreatePaymentFormData) => {
    createMutation.mutate(data, {
      onSuccess: () => navigate("/dashboard/payments"),
    });
  };

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold">Catat Pembayaran (Payment)</h1>
      <PaymentForm
        onSubmit={handleSubmit}
        defaultValues={prefilledValues}
        isSubmitting={createMutation.isPending}
      />
    </div>
  );
}

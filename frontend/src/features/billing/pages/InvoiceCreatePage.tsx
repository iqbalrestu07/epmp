import { useNavigate, useSearchParams } from "react-router-dom";
import { useCreateInvoice } from "../hooks";
import { InvoiceForm } from "../components/InvoiceForm";
import type { CreateInvoiceFormData } from "../schema";

export function InvoiceCreatePage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const createMutation = useCreateInvoice();

  const prefilledValues: Partial<CreateInvoiceFormData> = {
    contract_id: searchParams.get("contract_id") || undefined,
    tenant_id: searchParams.get("tenant_id") || undefined,
    amount: searchParams.get("amount") ? Number(searchParams.get("amount")) : undefined,
    currency: searchParams.get("currency") || "IDR",
  };

  const handleSubmit = (data: CreateInvoiceFormData) => {
    createMutation.mutate(data, {
      onSuccess: () => navigate("/dashboard/invoices"),
    });
  };

  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold">Buat Tagihan Baru (Invoice)</h1>
      <InvoiceForm
        onSubmit={handleSubmit}
        defaultValues={prefilledValues}
        isSubmitting={createMutation.isPending}
      />
    </div>
  );
}

import { useMemo } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useInvoices } from "@/features/billing/hooks";
import { useTenants } from "@/features/tenant/hooks";
import { SearchableSelect } from "@/components/ui/SearchableSelect";
import { formatCurrency } from "@/utils/currency";
import {
  createPaymentSchema,
  type CreatePaymentFormData,
} from "../schema";

interface PaymentFormProps {
  onSubmit: (data: CreatePaymentFormData) => void;
  defaultValues?: Partial<CreatePaymentFormData>;
  isSubmitting?: boolean;
}

export function PaymentForm({
  onSubmit,
  defaultValues,
  isSubmitting,
}: PaymentFormProps) {
  const {
    register,
    handleSubmit,
    setValue,
    watch,
    formState: { errors },
  } = useForm<CreatePaymentFormData>({
    resolver: zodResolver(createPaymentSchema),
    defaultValues: {
      payment_method: "Transfer",
      status: "Success",
      payment_date: new Date().toISOString().split("T")[0],
      ...defaultValues,
    },
  });

  const selectedInvoiceId = watch("invoice_id");
  const selectedTenantId = watch("tenant_id");

  const { data: invoicesData } = useInvoices({ per_page: 100 });
  const { data: tenantsData } = useTenants({ per_page: 100 });
  const invoices = useMemo(() => {
    return Array.isArray(invoicesData?.data) ? invoicesData.data : Array.isArray(invoicesData) ? invoicesData : [];
  }, [invoicesData]);
  const tenants = useMemo(() => {
    return Array.isArray(tenantsData?.data) ? tenantsData.data : Array.isArray(tenantsData) ? tenantsData : [];
  }, [tenantsData]);

  const tenantsMap = useMemo(() => {
    return new Map(tenants.map((t: any) => [t.id, t.full_name]));
  }, [tenants]);

  // Options for Tenant
  const tenantOptions = useMemo(() => {
    return tenants.map((t: any) => ({
      value: t.id,
      label: t.full_name,
      subLabel: `${t.phone || ''} ${t.email ? `• ${t.email}` : ''}`.trim(),
    }));
  }, [tenants]);

  // Filter invoices if tenant selected
  const filteredInvoices = useMemo(() => {
    if (!selectedTenantId) return invoices;
    return invoices.filter((inv: any) => inv.tenant_id === selectedTenantId);
  }, [invoices, selectedTenantId]);

  // Options for Invoice SearchableSelect
  const invoiceOptions = useMemo(() => {
    return filteredInvoices.map((inv: any) => {
      const tenantName = tenantsMap.get(inv.tenant_id) || 'Penyewa';
      const formattedAmount = formatCurrency(inv.amount || 0, (inv as any).currency || 'IDR');
      return {
        value: inv.id,
        label: `${tenantName} - ${formattedAmount}`,
        subLabel: `Invoice #${inv.id.slice(0, 8)} (${inv.status || 'Unpaid'})`,
        badge: inv.status,
      };
    });
  }, [filteredInvoices, tenantsMap]);

  // Handle invoice change -> auto fill tenant and amount
  const handleInvoiceChange = (invId: string) => {
    setValue("invoice_id", invId, { shouldValidate: true });
    const inv = invoices.find((item: any) => item.id === invId);
    if (inv) {
      if (inv.tenant_id) {
        setValue("tenant_id", inv.tenant_id, { shouldValidate: true });
      }
      if (inv.amount) {
        setValue("amount", Number(inv.amount), { shouldValidate: true });
      }
    }
  };

  const handleTenantChange = (tId: string) => {
    setValue("tenant_id", tId, { shouldValidate: true });
    const tenantInvoices = invoices.filter((item: any) => item.tenant_id === tId && item.status !== 'Paid');
    if (tenantInvoices.length > 0 && !selectedInvoiceId) {
      handleInvoiceChange(tenantInvoices[0].id);
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
      <div className="space-y-2">
        <Label htmlFor="tenant_id">Penyewa (Tenant) *</Label>
        <SearchableSelect
          id="tenant_id"
          options={tenantOptions}
          value={selectedTenantId}
          onChange={handleTenantChange}
          placeholder="Cari & pilih nama penyewa..."
          searchPlaceholder="Ketik nama atau no HP penyewa..."
          error={!!errors.tenant_id}
        />
        {errors.tenant_id && (
          <p className="text-xs text-red-600 mt-1">{errors.tenant_id.message}</p>
        )}
      </div>

      <div className="space-y-2">
        <div className="flex items-center justify-between">
          <Label htmlFor="invoice_id">Tagihan yang Dibayar (Invoice) *</Label>
          {selectedTenantId && (
            <span className="text-xs text-slate-500 font-normal">
              Menampilkan {filteredInvoices.length} tagihan untuk penyewa ini
            </span>
          )}
        </div>
        <SearchableSelect
          id="invoice_id"
          options={invoiceOptions}
          value={selectedInvoiceId}
          onChange={handleInvoiceChange}
          placeholder="Cari tagihan invoice yang akan dilunasi..."
          searchPlaceholder="Cari invoice berdasarkan nama, nominal, status..."
          error={!!errors.invoice_id}
        />
        {errors.invoice_id && (
          <p className="text-xs text-red-600 mt-1">{errors.invoice_id.message}</p>
        )}
      </div>
      <div className="space-y-2">
        <Label htmlFor="amount">Amount</Label>
        <Input
          id="amount"
          type="number"
          step="any"
          min="0"
          placeholder="0.00"
          {...register("amount", { valueAsNumber: true })}
        />
        {errors.amount && (
          <p className="text-xs text-red-600 mt-1">{errors.amount.message}</p>
        )}
      </div>
      <div className="space-y-2">
        <Label htmlFor="payment_date">Payment Date</Label>
        <Input
          id="payment_date"
          type="date"
          {...register("payment_date")}
        />
        {errors.payment_date && (
          <p className="text-xs text-red-600 mt-1">{errors.payment_date.message}</p>
        )}
      </div>
      <div className="space-y-2">
        <Label htmlFor="payment_method">Payment Method</Label>
        <select
          id="payment_method"
          {...register("payment_method")}
          className="w-full rounded-md border border-slate-300 px-3 py-2 text-sm bg-white text-slate-900 focus:outline-none focus:ring-2 focus:ring-orange"
        >
          <option value="Transfer">Transfer</option>
          <option value="CreditCard">CreditCard</option>
          <option value="Cash">Cash</option>
          <option value="EWallet">EWallet</option>
        </select>
        {errors.payment_method && (
          <p className="text-xs text-red-600 mt-1">{errors.payment_method.message}</p>
        )}
      </div>
      <div className="space-y-2">
        <Label htmlFor="status">Status</Label>
        <select
          id="status"
          {...register("status")}
          className="w-full rounded-md border border-slate-300 px-3 py-2 text-sm bg-white text-slate-900 focus:outline-none focus:ring-2 focus:ring-orange"
        >
          <option value="Pending">Pending</option>
          <option value="Success">Success</option>
          <option value="Failed">Failed</option>
          <option value="Refunded">Refunded</option>
        </select>
        {errors.status && (
          <p className="text-xs text-red-600 mt-1">{errors.status.message}</p>
        )}
      </div>
      <div className="space-y-2">
        <Label htmlFor="reference_number">Reference Number</Label>
        <Input
          id="reference_number"
          placeholder="e.g. TRX-12345"
          {...register("reference_number")}
        />
        {errors.reference_number && (
          <p className="text-xs text-red-600 mt-1">{errors.reference_number.message}</p>
        )}
      </div>

      <Button type="submit" disabled={isSubmitting} className="bg-orange hover:bg-orange/90 text-white w-full sm:w-auto">
        {isSubmitting ? "Saving..." : "Save"}
      </Button>
    </form>
  );
}

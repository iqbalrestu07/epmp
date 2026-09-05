import { useMemo } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  createInvoiceSchema,
  type CreateInvoiceFormData,
} from "../schema";

import { useContracts } from "../../contract/hooks";
import { useTenants } from "../../tenant/hooks";
import { useRooms } from "../../room/hooks";
import { CurrencySelect } from "@/components/ui/CurrencySelect";
import { SearchableSelect } from "@/components/ui/SearchableSelect";
import { formatCurrency } from "@/utils/currency";

interface InvoiceFormProps {
  onSubmit: (data: CreateInvoiceFormData) => void;
  defaultValues?: Partial<CreateInvoiceFormData>;
  isSubmitting?: boolean;
}

export function InvoiceForm({
  onSubmit,
  defaultValues,
  isSubmitting,
}: InvoiceFormProps) {
  const {
    register,
    handleSubmit,
    setValue,
    watch,
    formState: { errors },
  } = useForm<CreateInvoiceFormData>({
    resolver: zodResolver(createInvoiceSchema),
    defaultValues: {
      status: "Unpaid",
      amount: 0,
      payment_method: "Bank Transfer",
      currency: "IDR",
      ...defaultValues,
    },
  });

  const selectedTenantId = watch("tenant_id");
  const selectedContractId = watch("contract_id");

  const { data: contractsData } = useContracts({ per_page: 100 });
  const contracts = useMemo(() => {
    return Array.isArray(contractsData?.data) ? contractsData.data : Array.isArray(contractsData) ? contractsData : [];
  }, [contractsData]);

  const { data: tenantsData } = useTenants({ per_page: 100 });
  const tenants = useMemo(() => {
    return Array.isArray(tenantsData?.data) ? tenantsData.data : Array.isArray(tenantsData) ? tenantsData : [];
  }, [tenantsData]);

  const { data: roomsData } = useRooms({ per_page: 100 });
  const roomsMap = useMemo(() => {
    const list = Array.isArray(roomsData?.data) ? roomsData.data : Array.isArray(roomsData) ? roomsData : [];
    return new Map(list.map((r: any) => [r.id, r.name]));
  }, [roomsData]);

  // Options for Tenant SearchableSelect
  const tenantOptions = useMemo(() => {
    return tenants.map((t: any) => ({
      value: t.id,
      label: t.full_name,
      subLabel: `${t.phone || ''} ${t.email ? `• ${t.email}` : ''}`.trim(),
    }));
  }, [tenants]);

  // Filtered contracts based on selected tenant
  const availableContracts = useMemo(() => {
    if (!selectedTenantId) return contracts;
    return contracts.filter((c: any) => c.tenant_id === selectedTenantId);
  }, [contracts, selectedTenantId]);

  const contractOptions = useMemo(() => {
    return availableContracts.map((c: any) => {
      const roomName = roomsMap.get(c.room_id) || 'Unit Kamar';
      const priceText = formatCurrency(c.price_per_month || 0, c.currency || 'IDR');
      return {
        value: c.id,
        label: `${roomName} - ${priceText}/bln`,
        subLabel: `Kontrak #${c.id.slice(0, 8)} (${c.status})`,
        badge: c.status,
      };
    });
  }, [availableContracts, roomsMap]);

  // When tenant changes, auto-select contract if available
  const handleTenantChange = (tenantId: string) => {
    setValue("tenant_id", tenantId, { shouldValidate: true });
    const tenantContracts = contracts.filter((c: any) => c.tenant_id === tenantId);
    if (tenantContracts.length > 0 && !selectedContractId) {
      const c = tenantContracts[0];
      setValue("contract_id", c.id, { shouldValidate: true });
      const rent = c.monthly_rent || c.price_per_month;
      if (rent) {
        setValue("amount", Number(rent));
      }
      if (c.currency) {
        setValue("currency", c.currency);
      }
    }
  };

  // When contract selected, auto-fill amount & currency if empty
  const handleContractChange = (contractId: string) => {
    setValue("contract_id", contractId, { shouldValidate: true });
    const c = contracts.find((item: any) => item.id === contractId);
    if (c) {
      if (c.tenant_id) {
        setValue("tenant_id", c.tenant_id, { shouldValidate: true });
      }
      const rent = c.monthly_rent || c.price_per_month;
      if (rent) {
        setValue("amount", Number(rent));
      }
      if (c.currency) {
        setValue("currency", c.currency);
      }
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
      {/* 1. Pilih Penyewa (Tenant) Terlebih Dahulu dengan Live Search */}
      <div className="space-y-2">
        <Label htmlFor="tenant_id">Penyewa (Tenant) *</Label>
        <SearchableSelect
          id="tenant_id"
          options={tenantOptions}
          value={selectedTenantId}
          onChange={handleTenantChange}
          placeholder="Ketik untuk mencari nama atau telepon penyewa..."
          searchPlaceholder="Cari nama, no HP, email penyewa..."
          error={!!errors.tenant_id}
        />
        {errors.tenant_id && (
          <p className="text-xs text-red-600 mt-1">{errors.tenant_id.message}</p>
        )}
      </div>

      {/* 2. Pilih Kontrak yang Otomatis Terfilter */}
      <div className="space-y-2">
        <div className="flex items-center justify-between">
          <Label htmlFor="contract_id">Kontrak Sewa *</Label>
          {selectedTenantId && (
            <span className="text-xs text-slate-500 font-normal">
              Menampilkan {availableContracts.length} kontrak untuk penyewa ini
            </span>
          )}
        </div>
        <SearchableSelect
          id="contract_id"
          options={contractOptions}
          value={selectedContractId}
          onChange={handleContractChange}
          placeholder={
            selectedTenantId
              ? "Pilih kontrak sewa penyewa ini..."
              : "Pilih penyewa terlebih dahulu atau cari kontrak..."
          }
          searchPlaceholder="Cari berdasarkan kamar, nomor kontrak..."
          error={!!errors.contract_id}
        />
        {errors.contract_id && (
          <p className="text-xs text-red-600 mt-1">{errors.contract_id.message}</p>
        )}
      </div>

      <div className="space-y-2">
        <Label htmlFor="amount">Jumlah Tagihan & Mata Uang *</Label>
        <div className="flex gap-2">
          <CurrencySelect
            {...register("currency")}
            className="w-48 shrink-0"
          />
          <Input
            id="amount"
            type="number"
            step="any"
            placeholder="0"
            className="flex-1"
            {...register("amount", { valueAsNumber: true })}
          />
        </div>
        {errors.amount && (
          <p className="text-xs text-red-600 mt-1">{errors.amount.message}</p>
        )}
      </div>

      <div className="space-y-2">
        <Label htmlFor="status">Status</Label>
        <select
          id="status"
          {...register("status")}
          className="w-full rounded-lg border border-slate-300 bg-white px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-orange/30 text-slate-800"
        >
          <option value="Unpaid">Unpaid</option>
          <option value="Paid">Paid</option>
          <option value="Overdue">Overdue</option>
          <option value="Cancelled">Cancelled</option>
        </select>
        {errors.status && (
          <p className="text-xs text-red-600 mt-1">{errors.status.message}</p>
        )}
      </div>

      <div className="grid grid-cols-2 gap-4">
        <div className="space-y-2">
          <Label htmlFor="due_date">Due Date</Label>
          <Input
            id="due_date"
            type="date"
            {...register("due_date")}
          />
          {errors.due_date && (
            <p className="text-xs text-red-600 mt-1">{errors.due_date.message}</p>
          )}
        </div>

        <div className="space-y-2">
          <Label htmlFor="paid_date">Paid Date (Optional)</Label>
          <Input
            id="paid_date"
            type="date"
            {...register("paid_date")}
          />
          {errors.paid_date && (
            <p className="text-xs text-red-600 mt-1">{errors.paid_date.message}</p>
          )}
        </div>
      </div>

      <div className="space-y-2">
        <Label htmlFor="payment_method">Payment Method</Label>
        <select
          id="payment_method"
          {...register("payment_method")}
          className="w-full rounded-lg border border-slate-300 bg-white px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-orange/30 text-slate-800"
        >
          <option value="Bank Transfer">Bank Transfer</option>
          <option value="Cash">Cash</option>
          <option value="Credit Card">Credit Card</option>
          <option value="E-Wallet">E-Wallet (QRIS / GoPay / OVO)</option>
        </select>
        {errors.payment_method && (
          <p className="text-xs text-red-600 mt-1">{errors.payment_method.message}</p>
        )}
      </div>

      <div className="space-y-2">
        <Label htmlFor="notes">Notes (Optional)</Label>
        <Input
          id="notes"
          placeholder="e.g. Rent for September 2026"
          {...register("notes")}
        />
        {errors.notes && (
          <p className="text-xs text-red-600 mt-1">{errors.notes.message}</p>
        )}
      </div>

      <Button type="submit" disabled={isSubmitting} className="bg-orange hover:bg-orange/90 text-white w-full sm:w-auto">
        {isSubmitting ? "Saving..." : "Save"}
      </Button>
    </form>
  );
}

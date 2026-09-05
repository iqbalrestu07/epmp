import { useMemo } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useContracts } from "@/features/contract/hooks";
import { useTenants } from "@/features/tenant/hooks";
import { useRooms } from "@/features/room/hooks";
import { SearchableSelect } from "@/components/ui/SearchableSelect";
import { formatCurrency } from "@/utils/currency";
import {
  createDepositSchema,
  type CreateDepositFormData,
} from "../schema";

interface DepositFormProps {
  onSubmit: (data: CreateDepositFormData) => void;
  defaultValues?: Partial<CreateDepositFormData>;
  isSubmitting?: boolean;
}

export function DepositForm({
  onSubmit,
  defaultValues,
  isSubmitting,
}: DepositFormProps) {
  const {
    register,
    handleSubmit,
    setValue,
    watch,
    formState: { errors },
  } = useForm<CreateDepositFormData>({
    resolver: zodResolver(createDepositSchema),
    defaultValues: {
      status: "Collected",
      collection_date: new Date().toISOString().split("T")[0],
      ...defaultValues,
    },
  });

  const selectedTenantId = watch("tenant_id");
  const selectedContractId = watch("contract_id");

  const { data: contractsData } = useContracts({ per_page: 100 });
  const { data: tenantsData } = useTenants({ per_page: 100 });
  const { data: roomsData } = useRooms({ per_page: 100 });

  const contracts = useMemo(() => {
    return Array.isArray(contractsData?.data) ? contractsData.data : Array.isArray(contractsData) ? contractsData : [];
  }, [contractsData]);
  const tenants = useMemo(() => {
    return Array.isArray(tenantsData?.data) ? tenantsData.data : Array.isArray(tenantsData) ? tenantsData : [];
  }, [tenantsData]);
  const roomsMap = useMemo(() => {
    const list = Array.isArray(roomsData?.data) ? roomsData.data : Array.isArray(roomsData) ? roomsData : [];
    return new Map(list.map((r: any) => [r.id, r.name]));
  }, [roomsData]);

  // Options for Tenant
  const tenantOptions = useMemo(() => {
    return tenants.map((t: any) => ({
      value: t.id,
      label: t.full_name,
      subLabel: `${t.phone || ''} ${t.email ? `• ${t.email}` : ''}`.trim(),
    }));
  }, [tenants]);

  // Filter contracts based on tenant
  const filteredContracts = useMemo(() => {
    if (!selectedTenantId) return contracts;
    return contracts.filter((c: any) => c.tenant_id === selectedTenantId);
  }, [contracts, selectedTenantId]);

  const contractOptions = useMemo(() => {
    return filteredContracts.map((c: any) => {
      const roomName = roomsMap.get(c.room_id) || 'Unit Kamar';
      return {
        value: c.id,
        label: `${roomName} - Kontrak #${c.id.slice(0, 8)}`,
        subLabel: `Deposit: ${formatCurrency(c.deposit_amount || c.price_per_month || 0, c.currency || 'IDR')}`,
        badge: c.status,
      };
    });
  }, [filteredContracts, roomsMap]);

  const handleTenantChange = (tenantId: string) => {
    setValue("tenant_id", tenantId, { shouldValidate: true });
    const tenantContracts = contracts.filter((c: any) => c.tenant_id === tenantId);
    if (tenantContracts.length > 0 && !selectedContractId) {
      const c = tenantContracts[0];
      setValue("contract_id", c.id, { shouldValidate: true });
      const dep = c.deposit_amount || c.monthly_rent || c.price_per_month;
      if (dep) setValue("amount", Number(dep));
    }
  };

  const handleContractChange = (contractId: string) => {
    setValue("contract_id", contractId, { shouldValidate: true });
    const c = contracts.find((item: any) => item.id === contractId);
    if (c) {
      if (c.tenant_id) {
        setValue("tenant_id", c.tenant_id, { shouldValidate: true });
      }
      const dep = c.deposit_amount || c.monthly_rent || c.price_per_month;
      if (dep) setValue("amount", Number(dep));
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
          <Label htmlFor="contract_id">Kontrak Sewa *</Label>
          {selectedTenantId && (
            <span className="text-xs text-slate-500 font-normal">
              Menampilkan {filteredContracts.length} kontrak untuk penyewa ini
            </span>
          )}
        </div>
        <SearchableSelect
          id="contract_id"
          options={contractOptions}
          value={selectedContractId}
          onChange={handleContractChange}
          placeholder="Pilih kontrak sewa..."
          searchPlaceholder="Cari kamar atau nomor kontrak..."
          error={!!errors.contract_id}
        />
        {errors.contract_id && (
          <p className="text-xs text-red-600 mt-1">{errors.contract_id.message}</p>
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
        <Label htmlFor="status">Status</Label>
        <select
          id="status"
          {...register("status")}
          className="w-full rounded-md border border-slate-300 px-3 py-2 text-sm bg-white text-slate-900 focus:outline-none focus:ring-2 focus:ring-orange"
        >
          <option value="Collected">Collected</option>
          <option value="Returned">Returned</option>
          <option value="Forfeited">Forfeited</option>
        </select>
        {errors.status && (
          <p className="text-xs text-red-600 mt-1">{errors.status.message}</p>
        )}
      </div>
      <div className="space-y-2">
        <Label htmlFor="collection_date">Collection Date</Label>
        <Input
          id="collection_date"
          type="date"
          {...register("collection_date")}
        />
        {errors.collection_date && (
          <p className="text-xs text-red-600 mt-1">{errors.collection_date.message}</p>
        )}
      </div>
      <div className="space-y-2">
        <Label htmlFor="refund_date">Refund Date (Optional)</Label>
        <Input
          id="refund_date"
          type="date"
          {...register("refund_date")}
        />
        {errors.refund_date && (
          <p className="text-xs text-red-600 mt-1">{errors.refund_date.message}</p>
        )}
      </div>
      <div className="space-y-2">
        <Label htmlFor="notes">Notes</Label>
        <Input
          id="notes"
          placeholder="Optional notes"
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

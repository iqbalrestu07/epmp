import { useParams, useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Receipt, CreditCard, PlusCircle, ArrowLeft, Edit, Trash2 } from "lucide-react";
import { useContract, useDeleteContract } from "../hooks";
import { useTenant } from "@/features/tenant/hooks";
import { useRoom } from "@/features/room/hooks";
import { useProperty } from "@/features/property/hooks";
import { formatDate } from "@/utils/date";
import { formatCurrency } from "@/utils/currency";
import { AlertBadge } from "@/components/ui/AlertBadge";

export function ContractDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { data, isLoading } = useContract(id);
  const deleteMutation = useDeleteContract();

  const { data: tenantData } = useTenant(data?.tenant_id);
  const { data: roomData } = useRoom(data?.room_id);
  const { data: propertyData } = useProperty(data?.property_id);

  const handleDelete = () => {
    if (!confirm("Apakah Anda yakin ingin menghapus kontrak sewa ini?")) return;
    deleteMutation.mutate(id!, {
      onSuccess: () => navigate("/dashboard/contracts"),
    });
  };

  if (isLoading) return <div className="text-center py-12 text-slate-400">Memuat detail kontrak...</div>;
  if (!data) return <div className="text-center py-12 text-slate-400">Kontrak tidak ditemukan</div>;

  const tenantName = (tenantData as any)?.data?.full_name || (tenantData as any)?.full_name || "Memuat penyewa...";
  const tenantPhone = (tenantData as any)?.data?.phone || (tenantData as any)?.phone || "";
  const roomName = (roomData as any)?.data?.name || (roomData as any)?.name || "Kamar";
  const propName = (propertyData as any)?.data?.name || (propertyData as any)?.name || "Properti";

  return (
    <div className="space-y-6">
      {/* Header Bar */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-3">
            <h1 className="text-2xl font-bold text-slate-900">
              Kontrak #{id?.slice(0, 8)}
            </h1>
            <AlertBadge variant="active" label={data.status || "Aktif"} />
          </div>
          <p className="text-sm text-slate-500 mt-1">
            {propName} • {roomName} • {tenantName}
          </p>
        </div>

        <div className="flex items-center gap-2">
          <Button variant="outline" onClick={() => navigate("/dashboard/contracts")}>
            <ArrowLeft className="w-4 h-4 mr-1.5" /> Kembali
          </Button>
          <Button variant="outline" onClick={() => navigate(`/dashboard/contracts/${id}/edit`)}>
            <Edit className="w-4 h-4 mr-1.5" /> Edit
          </Button>
          <Button variant="destructive" onClick={handleDelete}>
            <Trash2 className="w-4 h-4 mr-1.5" /> Hapus
          </Button>
        </div>
      </div>

      {/* Quick Financial & Operations Actions Banner */}
      <div className="bg-gradient-to-r from-orange/10 via-amber-50 to-white rounded-2xl border border-orange/20 p-5 shadow-sm">
        <h3 className="text-xs font-bold text-orange uppercase tracking-wider mb-3">
          Aksi Cepat Keuangan & Operasional
        </h3>
        <div className="flex flex-wrap gap-3">
          <Button
            onClick={() =>
              navigate(
                `/invoices/new?contract_id=${id}&tenant_id=${data.tenant_id}&amount=${data.monthly_rent || 0}&currency=${(data as any).currency || 'IDR'}`
              )
            }
            className="bg-orange hover:bg-orange/90 text-slate-900 font-semibold shadow-sm"
          >
            <Receipt className="w-4 h-4 mr-1.5" /> Buat Tagihan (Invoice)
          </Button>

          <Button
            variant="outline"
            onClick={() =>
              navigate(
                `/deposits/new?contract_id=${id}&tenant_id=${data.tenant_id}&amount=${data.deposit_amount || data.monthly_rent || 0}`
              )
            }
            className="bg-white hover:bg-slate-50 text-slate-800 border-slate-300 font-medium"
          >
            <CreditCard className="w-4 h-4 mr-1.5" /> Catat Uang Jaminan (Deposit)
          </Button>

          <Button
            variant="outline"
            onClick={() => navigate(`/charges/new?contract_id=${id}`)}
            className="bg-white hover:bg-slate-50 text-slate-800 border-slate-300 font-medium"
          >
            <PlusCircle className="w-4 h-4 mr-1.5" /> Tambah Biaya (Charge)
          </Button>
        </div>
      </div>

      {/* Contract Specification Card */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-5 bg-white rounded-2xl border border-slate-200 shadow-sm p-6">
        <div className="space-y-4">
          <div>
            <dt className="text-xs font-semibold text-slate-400 uppercase tracking-wider">Penyewa (Tenant)</dt>
            <dd className="text-base font-semibold text-slate-900 mt-1">
              {tenantName} {tenantPhone ? `(${tenantPhone})` : ''}
            </dd>
          </div>

          <div>
            <dt className="text-xs font-semibold text-slate-400 uppercase tracking-wider">Unit Kamar & Properti</dt>
            <dd className="text-sm font-medium text-slate-800 mt-1">
              {roomName} — {propName}
            </dd>
          </div>

          <div>
            <dt className="text-xs font-semibold text-slate-400 uppercase tracking-wider">Periode Kontrak Sewa</dt>
            <dd className="text-sm text-slate-800 mt-1">
              {formatDate(data.start_date)} s/d {formatDate(data.end_date)}
            </dd>
          </div>
        </div>

        <div className="space-y-4 md:border-l md:border-slate-100 md:pl-6">
          <div>
            <dt className="text-xs font-semibold text-slate-400 uppercase tracking-wider">Harga Sewa Bulanan</dt>
            <dd className="text-lg font-bold text-slate-900 mt-1">
              {formatCurrency(data.monthly_rent, (data as any).currency || "IDR")} / bulan
            </dd>
          </div>

          <div>
            <dt className="text-xs font-semibold text-slate-400 uppercase tracking-wider">Uang Jaminan (Deposit)</dt>
            <dd className="text-base font-semibold text-slate-800 mt-1">
              {formatCurrency(data.deposit_amount, (data as any).currency || "IDR")}
            </dd>
          </div>

          <div>
            <dt className="text-xs font-semibold text-slate-400 uppercase tracking-wider">Syarat & Ketentuan (Terms)</dt>
            <dd className="text-xs text-slate-600 mt-1 whitespace-pre-line bg-slate-50 p-2.5 rounded-lg border border-slate-100">
              {data.terms || "Standard tenancy terms apply."}
            </dd>
          </div>
        </div>
      </div>
    </div>
  );
}

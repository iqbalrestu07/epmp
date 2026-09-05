import { useParams, useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { CreditCard, AlertTriangle, ArrowLeft, Edit, Trash2, Percent } from "lucide-react";
import { useInvoice, useDeleteInvoice } from "../hooks";
import { useTenant } from "@/features/tenant/hooks";
import { formatDate } from "@/utils/date";
import { formatCurrency } from "@/utils/currency";
import { AlertBadge } from "@/components/ui/AlertBadge";

import { usePayments } from "@/features/payment/hooks";

export function InvoiceDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { data, isLoading } = useInvoice(id);
  const deleteMutation = useDeleteInvoice();
  const { data: paymentsData } = usePayments({ per_page: 200 });

  const { data: tenantData } = useTenant(data?.tenant_id);

  const handleDelete = () => {
    if (!confirm("Apakah Anda yakin ingin menghapus tagihan invoice ini?")) return;
    deleteMutation.mutate(id!, {
      onSuccess: () => navigate("/dashboard/invoices"),
    });
  };

  if (isLoading) return <div className="text-center py-12 text-slate-400">Memuat tagihan...</div>;
  if (!data) return <div className="text-center py-12 text-slate-400">Tagihan tidak ditemukan</div>;

  const payments = Array.isArray(paymentsData?.data) ? paymentsData.data : Array.isArray(paymentsData) ? paymentsData : [];
  const hasSuccessfulPayment = payments.some(
    (p: any) => (p.status === "Success" || p.status === "Completed") && p.invoice_id === id
  );
  const isPaid = data.status === "Paid" || hasSuccessfulPayment;

  const tenantName = (tenantData as any)?.data?.full_name || (tenantData as any)?.full_name || "Penyewa";
  const tenantPhone = (tenantData as any)?.data?.phone || (tenantData as any)?.phone || "";
  const contractNum = data.contract_id ? `#${data.contract_id.slice(0, 8)}` : "—";

  return (
    <div className="space-y-6">
      {/* Header Bar */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-3">
            <h1 className="text-2xl font-bold text-slate-900">
              Tagihan #{id?.slice(0, 8)}
            </h1>
            {isPaid ? (
              <AlertBadge variant="paid" label="Lunas" />
            ) : (
              <AlertBadge variant="pending" label={data.status || "Belum Lunas"} />
            )}
          </div>
          <p className="text-sm text-slate-500 mt-1">
            Penyewa: {tenantName} • Kontrak {contractNum}
          </p>
        </div>

        <div className="flex items-center gap-2">
          <Button variant="outline" onClick={() => navigate("/dashboard/invoices")}>
            <ArrowLeft className="w-4 h-4 mr-1.5" /> Kembali
          </Button>
          <Button variant="outline" onClick={() => navigate(`/dashboard/invoices/${id}/edit`)}>
            <Edit className="w-4 h-4 mr-1.5" /> Edit
          </Button>
          <Button variant="destructive" onClick={handleDelete}>
            <Trash2 className="w-4 h-4 mr-1.5" /> Hapus
          </Button>
        </div>
      </div>

      {/* Quick Action Banner */}
      {!isPaid && (
        <div className="bg-gradient-to-r from-emerald-50 via-teal-50 to-white rounded-2xl border border-emerald-200 p-5 shadow-sm">
          <h3 className="text-xs font-bold text-emerald-800 uppercase tracking-wider mb-3">
            Tindakan Pembayaran & Penyesuaian
          </h3>
          <div className="flex flex-wrap gap-3">
            <Button
              onClick={() =>
                navigate(
                  `/payments/new?invoice_id=${id}&tenant_id=${data.tenant_id}&amount=${data.amount}`
                )
              }
              className="bg-emerald-600 hover:bg-emerald-700 text-white font-semibold shadow-sm"
            >
              <CreditCard className="w-4 h-4 mr-1.5" /> Catat Pembayaran (Pay Now)
            </Button>

            <Button
              variant="outline"
              onClick={() => navigate(`/penalties/new?invoice_id=${id}`)}
              className="bg-white hover:bg-slate-50 text-slate-800 border-slate-300 font-medium"
            >
              <AlertTriangle className="w-4 h-4 mr-1.5 text-amber-500" /> Denda Keterlambatan
            </Button>

            <Button
              variant="outline"
              onClick={() => navigate(`/adjustments/new?invoice_id=${id}`)}
              className="bg-white hover:bg-slate-50 text-slate-800 border-slate-300 font-medium"
            >
              <Percent className="w-4 h-4 mr-1.5 text-blue-500" /> Beri Diskon / Keringanan
            </Button>
          </div>
        </div>
      )}

      {/* Invoice Details Card */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-5 bg-white rounded-2xl border border-slate-200 shadow-sm p-6">
        <div className="space-y-4">
          <div>
            <dt className="text-xs font-semibold text-slate-400 uppercase tracking-wider">Penyewa (Tenant)</dt>
            <dd className="text-base font-semibold text-slate-900 mt-1">
              {tenantName} {tenantPhone ? `(${tenantPhone})` : ''}
            </dd>
          </div>

          <div>
            <dt className="text-xs font-semibold text-slate-400 uppercase tracking-wider">Kontrak Sewa</dt>
            <dd className="text-sm font-medium text-slate-800 mt-1">
              Kontrak {contractNum}
            </dd>
          </div>

          <div>
            <dt className="text-xs font-semibold text-slate-400 uppercase tracking-wider">Metode Pembayaran</dt>
            <dd className="text-sm text-slate-800 mt-1">
              {data.payment_method || "Bank Transfer"}
            </dd>
          </div>
        </div>

        <div className="space-y-4 md:border-l md:border-slate-100 md:pl-6">
          <div>
            <dt className="text-xs font-semibold text-slate-400 uppercase tracking-wider">Total Tagihan</dt>
            <dd className="text-2xl font-bold text-slate-900 mt-1">
              {formatCurrency(data.amount, (data as any).currency || "IDR")}
            </dd>
          </div>

          <div>
            <dt className="text-xs font-semibold text-slate-400 uppercase tracking-wider">Tanggal Jatuh Tempo</dt>
            <dd className="text-sm font-medium text-slate-800 mt-1">
              {formatDate(data.due_date)}
            </dd>
          </div>

          <div>
            <dt className="text-xs font-semibold text-slate-400 uppercase tracking-wider">Tanggal Pelunasan</dt>
            <dd className="text-sm text-slate-800 mt-1">
              {data.paid_date ? formatDate(data.paid_date) : "Belum Lunas"}
            </dd>
          </div>

          <div>
            <dt className="text-xs font-semibold text-slate-400 uppercase tracking-wider">Catatan</dt>
            <dd className="text-xs text-slate-600 mt-1">
              {data.notes || "—"}
            </dd>
          </div>
        </div>
      </div>
    </div>
  );
}

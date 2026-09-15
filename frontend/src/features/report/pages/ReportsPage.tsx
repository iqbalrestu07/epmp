import { Fragment, useState } from "react";
import { Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useOccupancyReport, useRevenueReport, useArAgingReport } from "../hooks";
import { usePropertys } from "@/features/property/hooks";
import { formatCurrency } from "@/utils/currency";
import { BarChart3, TrendingUp, AlertCircle } from "lucide-react";

type Tab = "occupancy" | "revenue" | "aging";

const BUCKET_LABELS: Record<string, string> = {
  Current: "Belum Jatuh Tempo",
  "1-30": "1–30 hari",
  "31-60": "31–60 hari",
  "61-90": "61–90 hari",
  "90+": "90+ hari",
};

function StatusBadge({ label, value, total }: { label: string; value: number; total: number }) {
  const pct = total > 0 ? Math.round((value / total) * 100) : 0;
  return (
    <div className="rounded-lg border bg-white p-4 shadow-sm">
      <p className="text-xs text-slate-500 uppercase tracking-wide">{label}</p>
      <p className="mt-1 text-2xl font-bold text-slate-900">{value}</p>
      <p className="text-xs text-slate-400">{pct}% dari total</p>
    </div>
  );
}

function OccupancyReport() {
  const [propertyId, setPropertyId] = useState("");
  const { data: properties } = usePropertys({ page: 1, per_page: 100 });
  const { data, isLoading } = useOccupancyReport({
    property_id: propertyId || undefined,
  });

  return (
    <div className="space-y-4">
      <div className="flex gap-2">
        <select
          className="h-10 rounded-md border border-input bg-background px-3 text-sm"
          value={propertyId}
          onChange={(e) => setPropertyId(e.target.value)}
        >
          <option value="">Semua Properti</option>
          {(properties?.data ?? []).map((p: any) => (
            <option key={p.id} value={p.id}>{p.name}</option>
          ))}
        </select>
      </div>

      {isLoading ? (
        <div className="text-center py-12 text-slate-400">Memuat laporan okupansi...</div>
      ) : !data ? null : (
        <>
          <div className="grid grid-cols-2 md:grid-cols-5 gap-3">
            <StatusBadge label="Total Kamar" value={data.totals.total_rooms} total={data.totals.total_rooms} />
            <StatusBadge label="Terisi" value={data.totals.occupied} total={data.totals.total_rooms} />
            <StatusBadge label="Direservasi" value={data.totals.reserved} total={data.totals.total_rooms} />
            <StatusBadge label="Maintenance" value={data.totals.maintenance} total={data.totals.total_rooms} />
            <StatusBadge label="Tersedia" value={data.totals.available} total={data.totals.total_rooms} />
          </div>
          <p className="text-sm text-slate-600">
            Tingkat okupansi keseluruhan: <span className="font-semibold">{data.totals.occupancy_rate}%</span>
          </p>

          <div className="rounded-lg border bg-white shadow-sm overflow-hidden">
            <table className="w-full text-sm">
              <thead className="bg-slate-50 text-slate-600">
                <tr>
                  <th className="text-left px-4 py-3 font-medium">Properti / Gedung</th>
                  <th className="text-right px-4 py-3 font-medium">Total</th>
                  <th className="text-right px-4 py-3 font-medium">Terisi</th>
                  <th className="text-right px-4 py-3 font-medium">Direservasi</th>
                  <th className="text-right px-4 py-3 font-medium">Maintenance</th>
                  <th className="text-right px-4 py-3 font-medium">Tersedia</th>
                  <th className="text-right px-4 py-3 font-medium">Okupansi</th>
                </tr>
              </thead>
              <tbody>
                {data.properties.length === 0 && (
                  <tr><td colSpan={7} className="px-4 py-8 text-center text-slate-400">Tidak ada data</td></tr>
                )}
                {data.properties.map((p) => (
                  <Fragment key={p.property_id}>
                    <tr className="border-t hover:bg-slate-50">
                      <td className="px-4 py-3 font-medium text-slate-900">{p.property_name}</td>
                      <td className="px-4 py-3 text-right">{p.total_rooms}</td>
                      <td className="px-4 py-3 text-right">{p.occupied}</td>
                      <td className="px-4 py-3 text-right">{p.reserved}</td>
                      <td className="px-4 py-3 text-right">{p.maintenance}</td>
                      <td className="px-4 py-3 text-right">{p.available}</td>
                      <td className="px-4 py-3 text-right font-medium">{p.occupancy_rate}%</td>
                    </tr>
                    {p.buildings.map((b) => (
                      <tr key={b.building_id} className="border-t bg-slate-50/50 text-slate-600">
                        <td className="px-4 py-2 pl-8">↳ {b.building_name}</td>
                        <td className="px-4 py-2 text-right">{b.total_rooms}</td>
                        <td className="px-4 py-2 text-right">{b.occupied}</td>
                        <td className="px-4 py-2 text-right">{b.reserved}</td>
                        <td className="px-4 py-2 text-right">{b.maintenance}</td>
                        <td className="px-4 py-2 text-right">{b.available}</td>
                        <td className="px-4 py-2 text-right">{b.occupancy_rate}%</td>
                      </tr>
                    ))}
                  </Fragment>
                ))}
              </tbody>
            </table>
          </div>
        </>
      )}
    </div>
  );
}

function RevenueReport() {
  const today = new Date();
  const sixMonthsAgo = new Date(today.getFullYear(), today.getMonth() - 5, 1);
  const [from, setFrom] = useState(sixMonthsAgo.toISOString().slice(0, 10));
  const [to, setTo] = useState(
    new Date(today.getFullYear(), today.getMonth() + 1, 1).toISOString().slice(0, 10)
  );

  const { data, isLoading } = useRevenueReport({ from, to });

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-end gap-3">
        <div>
          <label className="block text-xs text-slate-500 mb-1">Dari</label>
          <Input type="date" value={from} onChange={(e) => setFrom(e.target.value)} />
        </div>
        <div>
          <label className="block text-xs text-slate-500 mb-1">Sampai</label>
          <Input type="date" value={to} onChange={(e) => setTo(e.target.value)} />
        </div>
      </div>

      {isLoading ? (
        <div className="text-center py-12 text-slate-400">Memuat laporan pendapatan...</div>
      ) : !data ? null : (
        <div className="rounded-lg border bg-white shadow-sm overflow-hidden">
          <table className="w-full text-sm">
            <thead className="bg-slate-50 text-slate-600">
              <tr>
                <th className="text-left px-4 py-3 font-medium">Periode</th>
                <th className="text-left px-4 py-3 font-medium">Mata Uang</th>
                <th className="text-right px-4 py-3 font-medium"># Invoice</th>
                <th className="text-right px-4 py-3 font-medium">Ditagihkan</th>
                <th className="text-right px-4 py-3 font-medium"># Lunas</th>
                <th className="text-right px-4 py-3 font-medium"># Pembayaran</th>
                <th className="text-right px-4 py-3 font-medium">Diterima</th>
              </tr>
            </thead>
            <tbody>
              {data.rows.length === 0 && (
                <tr><td colSpan={7} className="px-4 py-8 text-center text-slate-400">Tidak ada data pada rentang ini</td></tr>
              )}
              {data.rows.map((r) => (
                <tr key={`${r.period}-${r.currency}`} className="border-t hover:bg-slate-50">
                  <td className="px-4 py-3 font-medium text-slate-900">{r.period}</td>
                  <td className="px-4 py-3">{r.currency}</td>
                  <td className="px-4 py-3 text-right">{r.invoice_count}</td>
                  <td className="px-4 py-3 text-right">{formatCurrency(r.invoiced_amount, r.currency)}</td>
                  <td className="px-4 py-3 text-right">{r.paid_invoice_count}</td>
                  <td className="px-4 py-3 text-right">{r.payment_count}</td>
                  <td className="px-4 py-3 text-right font-medium text-emerald-700">{formatCurrency(r.received_amount, r.currency)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}

function ArAgingReport() {
  const { data, isLoading } = useArAgingReport();

  return (
    <div className="space-y-4">
      {isLoading ? (
        <div className="text-center py-12 text-slate-400">Memuat laporan aging piutang...</div>
      ) : !data ? null : (
        <>
          <div className="grid grid-cols-2 md:grid-cols-5 gap-3">
            {["Current", "1-30", "31-60", "61-90", "90+"].map((bucket) => {
              const rows = data.buckets.filter((b) => b.bucket === bucket);
              const count = rows.reduce((s, b) => s + b.invoice_count, 0);
              return (
                <div key={bucket} className="rounded-lg border bg-white p-4 shadow-sm">
                  <p className="text-xs text-slate-500 uppercase tracking-wide">{BUCKET_LABELS[bucket] ?? bucket}</p>
                  <p className="mt-1 text-xl font-bold text-slate-900">{count} invoice</p>
                  {rows.map((b) => (
                    <p key={b.currency} className="text-xs text-slate-500">
                      {formatCurrency(b.outstanding, b.currency)}
                    </p>
                  ))}
                </div>
              );
            })}
          </div>

          <div className="rounded-lg border bg-white shadow-sm overflow-hidden">
            <table className="w-full text-sm">
              <thead className="bg-slate-50 text-slate-600">
                <tr>
                  <th className="text-left px-4 py-3 font-medium">Penyewa</th>
                  <th className="text-left px-4 py-3 font-medium">Status</th>
                  <th className="text-left px-4 py-3 font-medium">Jatuh Tempo</th>
                  <th className="text-right px-4 py-3 font-medium">Hari Telat</th>
                  <th className="text-left px-4 py-3 font-medium">Bucket</th>
                  <th className="text-right px-4 py-3 font-medium">Tagihan</th>
                  <th className="text-right px-4 py-3 font-medium">Terbayar</th>
                  <th className="text-right px-4 py-3 font-medium">Sisa</th>
                </tr>
              </thead>
              <tbody>
                {data.invoices.length === 0 && (
                  <tr><td colSpan={8} className="px-4 py-8 text-center text-slate-400">Tidak ada piutang outstanding</td></tr>
                )}
                {data.invoices.map((inv) => (
                  <tr key={inv.invoice_id} className="border-t hover:bg-slate-50">
                    <td className="px-4 py-3 font-medium text-slate-900">{inv.tenant_name}</td>
                    <td className="px-4 py-3">{inv.status}</td>
                    <td className="px-4 py-3">{new Date(inv.due_date).toLocaleDateString("id-ID")}</td>
                    <td className="px-4 py-3 text-right">{inv.days_overdue}</td>
                    <td className="px-4 py-3">{BUCKET_LABELS[inv.bucket] ?? inv.bucket}</td>
                    <td className="px-4 py-3 text-right">{formatCurrency(inv.amount, inv.currency)}</td>
                    <td className="px-4 py-3 text-right">{formatCurrency(inv.paid_amount, inv.currency)}</td>
                    <td className="px-4 py-3 text-right font-medium text-red-600">{formatCurrency(inv.outstanding, inv.currency)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </>
      )}
    </div>
  );
}

export function ReportsPage() {
  const [tab, setTab] = useState<Tab>("occupancy");

  const tabs: { key: Tab; label: string; icon: React.ReactNode }[] = [
    { key: "occupancy", label: "Okupansi", icon: <BarChart3 className="w-4 h-4" /> },
    { key: "revenue", label: "Pendapatan", icon: <TrendingUp className="w-4 h-4" /> },
    { key: "aging", label: "Aging Piutang", icon: <AlertCircle className="w-4 h-4" /> },
  ];

  return (
    <div className="space-y-6">
      <div>
        <div className="flex items-center gap-2 text-sm text-slate-500 mb-1">
          <Link to="/dashboard" className="hover:text-slate-900">Dashboard</Link>
          <span>/</span>
          <span className="text-slate-900 font-medium">Laporan</span>
        </div>
        <h1 className="text-2xl font-bold text-slate-900">Laporan</h1>
      </div>

      <div className="flex gap-2 border-b">
        {tabs.map((t) => (
          <Button
            key={t.key}
            variant={tab === t.key ? "default" : "ghost"}
            size="sm"
            onClick={() => setTab(t.key)}
            className="rounded-b-none"
          >
            <span className="mr-1">{t.icon}</span>
            {t.label}
          </Button>
        ))}
      </div>

      {tab === "occupancy" && <OccupancyReport />}
      {tab === "revenue" && <RevenueReport />}
      {tab === "aging" && <ArAgingReport />}
    </div>
  );
}

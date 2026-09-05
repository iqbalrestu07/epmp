import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Building2, Users, DoorOpen, Banknote, Receipt,
  AlertCircle, Clock, AlertTriangle, XCircle, ArrowRight,
  MessageSquare, Send, CheckCircle2, Box, RefreshCw
} from 'lucide-react';
import { motion } from 'framer-motion';
import { api } from '../services/api';
import { AlertBadge } from '../components/ui/AlertBadge';
import { TenantRoomRoster } from '@/features/occupancy/components/TenantRoomRoster';

interface DashboardMetrics {
  total_properties: number;
  total_rooms: number;
  occupied_rooms: number;
  vacant_rooms: number;
  occupancy_rate: number;
  active_tenants: number;
  monthly_revenue: number;
  unpaid_invoices_amount: number;
  open_work_orders: number;
}

interface OverdueInvoiceAlert {
  id: string;
  tenant_id: string;
  tenant_name: string;
  tenant_email: string;
  tenant_phone: string;
  amount: number;
  due_date: string;
  days_overdue: number;
  property_name: string;
  room_name: string;
}

interface DueSoonInvoiceAlert {
  id: string;
  tenant_id: string;
  tenant_name: string;
  tenant_email: string;
  tenant_phone: string;
  amount: number;
  due_date: string;
  days_until_due: number;
  property_name: string;
  room_name: string;
}

interface ExpiringContractAlert {
  id: string;
  tenant_id: string;
  tenant_name: string;
  tenant_email: string;
  tenant_phone: string;
  property_name: string;
  room_name: string;
  start_date: string;
  end_date: string;
  days_until_expiration: number;
  monthly_rent: number;
}

interface ExpiredContractAlert {
  id: string;
  tenant_id: string;
  tenant_name: string;
  tenant_email: string;
  tenant_phone: string;
  property_name: string;
  room_name: string;
  start_date: string;
  end_date: string;
  days_expired: number;
  monthly_rent: number;
}

interface DashboardAlerts {
  overdue_invoices: OverdueInvoiceAlert[];
  due_soon_invoices: DueSoonInvoiceAlert[];
  expiring_contracts: ExpiringContractAlert[];
  expired_contracts: ExpiredContractAlert[];
}

interface PropertyOccupancySummary {
  property_id: string;
  property_name: string;
  total_rooms: number;
  occupied: number;
  rate: number;
}

interface RecentActivity {
  id: string;
  type: string;
  title: string;
  description: string;
  timestamp: string;
}

interface DashboardSummaryResponse {
  metrics: DashboardMetrics;
  alerts: DashboardAlerts;
  property_occupancy: PropertyOccupancySummary[];
  recent_activities: RecentActivity[];
}

export default function DashboardPage() {
  const navigate = useNavigate();
  const [data, setData] = useState<DashboardSummaryResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<'all' | 'overdue' | 'due_soon' | 'expiring' | 'expired'>('all');

  const fetchSummary = async () => {
    setLoading(true);
    try {
      const res = await api.get<DashboardSummaryResponse>('/dashboard/summary');
      setData(res);
    } catch (err) {
      console.error('Failed to load dashboard summary:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchSummary();
  }, []);

  const overdue = data?.alerts?.overdue_invoices || [];
  const dueSoon = data?.alerts?.due_soon_invoices || [];
  const expiring = data?.alerts?.expiring_contracts || [];
  const expired = data?.alerts?.expired_contracts || [];
  const totalAlertsCount = overdue.length + dueSoon.length + expiring.length + expired.length;

  const formatRupiah = (val: number) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      maximumFractionDigits: 0,
    }).format(val || 0);
  };

  const getCleanPhone = (phone: string) => {
    if (!phone) return '';
    let clean = phone.replace(/[^0-9]/g, '');
    if (clean.startsWith('0')) {
      clean = '62' + clean.slice(1);
    }
    return clean;
  };

  const handleSendReminder = (phone: string, name: string, type: 'overdue' | 'due_soon' | 'expiring', amountOrRent: number, deadlineStr: string) => {
    const cleanPhone = getCleanPhone(phone);
    let message = '';
    if (type === 'overdue') {
      message = `Halo ${name}, mengingatkan bahwa tagihan sewa Anda sebesar ${formatRupiah(amountOrRent)} telah melewati jatuh tempo (${new Date(deadlineStr).toLocaleDateString('id-ID')}). Mohon segera menyelesaikan pembayaran. Terima kasih.`;
    } else if (type === 'due_soon') {
      message = `Halo ${name}, tagihan sewa Anda sebesar ${formatRupiah(amountOrRent)} akan jatuh tempo pada ${new Date(deadlineStr).toLocaleDateString('id-ID')}. Mohon persiapkan pembayaran tepat waktu. Terima kasih.`;
    } else {
      message = `Halo ${name}, masa sewa unit Anda akan berakhir pada ${new Date(deadlineStr).toLocaleDateString('id-ID')}. Apabila Anda ingin memperpanjang masa sewa, silakan hubungi pengelola. Terima kasih.`;
    }

    if (cleanPhone) {
      window.open(`https://wa.me/${cleanPhone}?text=${encodeURIComponent(message)}`, '_blank');
    } else {
      alert(`Pesan Peringatan:\n\n${message}`);
    }
  };

  return (
    <div className="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500 pb-12">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-3xl font-display font-bold text-slate-900 tracking-tight">Dashboard Ringkasan</h1>
          <p className="text-slate-500 text-sm mt-1">Pusat komando operasional, status tenant, dan kontrol finansial real-time.</p>
        </div>
        <div className="flex items-center gap-3">
          <button
            onClick={fetchSummary}
            disabled={loading}
            className="flex items-center gap-2 px-4 py-2 bg-white border border-slate-200 text-slate-700 hover:bg-slate-50 rounded-xl text-sm font-medium shadow-2xs transition-all"
          >
            <RefreshCw size={15} className={loading ? 'animate-spin' : ''} />
            Muat Ulang
          </button>
          <button
            onClick={() => navigate('/dashboard/properties/interactive')}
            className="flex items-center gap-2 px-4 py-2 bg-orange hover:bg-orange/90 text-white rounded-xl text-sm font-medium shadow-xs transition-all"
          >
            <Box size={16} />
            3D Spatial View
          </button>
        </div>
      </div>

      {/* KPI Metrics Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
        {/* Okupansi */}
        <motion.div
          whileHover={{ y: -3 }}
          className="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-2xs flex flex-col justify-between"
        >
          <div className="flex justify-between items-start mb-3">
            <span className="text-xs font-semibold uppercase tracking-wider text-slate-500">Tingkat Okupansi</span>
            <div className="w-10 h-10 rounded-xl bg-purple-50 text-purple-600 flex items-center justify-center">
              <DoorOpen size={20} />
            </div>
          </div>
          <div>
            <div className="flex items-baseline gap-2">
              <h3 className="text-3xl font-bold text-slate-900">{data?.metrics?.occupancy_rate ?? 0}%</h3>
              <span className="text-xs text-slate-500 font-medium">terisi</span>
            </div>
            <div className="w-full bg-slate-100 rounded-full h-2 mt-3 overflow-hidden">
              <div
                className="bg-purple-600 h-2 rounded-full transition-all duration-800"
                style={{ width: `${Math.min(100, data?.metrics?.occupancy_rate ?? 0)}%` }}
              />
            </div>
            <p className="text-xs text-slate-500 mt-2 flex justify-between">
              <span>{data?.metrics?.occupied_rooms ?? 0} Unit Terisi</span>
              <span>{data?.metrics?.vacant_rooms ?? 0} Kosong</span>
            </p>
          </div>
        </motion.div>

        {/* Penyewa Aktif */}
        <motion.div
          whileHover={{ y: -3 }}
          className="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-2xs flex flex-col justify-between"
        >
          <div className="flex justify-between items-start mb-3">
            <span className="text-xs font-semibold uppercase tracking-wider text-slate-500">Penyewa Aktif</span>
            <div className="w-10 h-10 rounded-xl bg-emerald-50 text-emerald-600 flex items-center justify-center">
              <Users size={20} />
            </div>
          </div>
          <div>
            <h3 className="text-3xl font-bold text-slate-900">{data?.metrics?.active_tenants ?? 0}</h3>
            <p className="text-xs text-slate-500 mt-1">Tenant dengan kontrak sewa aktif</p>
            <div className="mt-4 pt-3 border-t border-slate-100 flex items-center justify-between text-xs text-slate-600 font-medium">
              <span>Properti Dikelola</span>
              <span className="font-semibold text-slate-900">{data?.metrics?.total_properties ?? 0} Properti</span>
            </div>
          </div>
        </motion.div>

        {/* Pendapatan Terkumpul Bulan Ini */}
        <motion.div
          whileHover={{ y: -3 }}
          className="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-2xs flex flex-col justify-between"
        >
          <div className="flex justify-between items-start mb-3">
            <span className="text-xs font-semibold uppercase tracking-wider text-slate-500">Penerimaan Bulan Ini</span>
            <div className="w-10 h-10 rounded-xl bg-emerald-50 text-emerald-600 flex items-center justify-center">
              <Banknote size={20} />
            </div>
          </div>
          <div>
            <h3 className="text-2xl font-bold text-slate-900 truncate">
              {formatRupiah(data?.metrics?.monthly_revenue ?? 0)}
            </h3>
            <p className="text-xs text-emerald-700 font-medium mt-1">Penerimaan sewa terverifikasi</p>
            <div className="mt-4 pt-3 border-t border-slate-100 flex items-center justify-between text-xs text-slate-600 font-medium">
              <span>Tiket Perbaikan</span>
              <span className="font-semibold text-slate-900">{data?.metrics?.open_work_orders ?? 0} Aktif</span>
            </div>
          </div>
        </motion.div>

        {/* Tagihan Belum Dibayar / Piutang */}
        <motion.div
          whileHover={{ y: -3 }}
          className="bg-white p-5 rounded-2xl border border-slate-200/80 shadow-2xs flex flex-col justify-between"
        >
          <div className="flex justify-between items-start mb-3">
            <span className="text-xs font-semibold uppercase tracking-wider text-slate-500">Total Piutang Sewa</span>
            <div className="w-10 h-10 rounded-xl bg-rose-50 text-rose-600 flex items-center justify-center">
              <Receipt size={20} />
            </div>
          </div>
          <div>
            <h3 className="text-2xl font-bold text-rose-600 truncate">
              {formatRupiah(data?.metrics?.unpaid_invoices_amount ?? 0)}
            </h3>
            <p className="text-xs text-slate-500 mt-1">Akumulasi tagihan belum lunas</p>
            <div className="mt-4 pt-3 border-t border-slate-100 flex items-center justify-between text-xs text-slate-600 font-medium">
              <span>Overdue Tagihan</span>
              <span className="font-bold text-red-600">{overdue.length} Invoice</span>
            </div>
          </div>
        </motion.div>
      </div>

      {/* ─── PUSAT PERHATIAN KHUSUS TENANT (ALERT CENTER) ─── */}
      <div className="bg-white rounded-2xl border border-slate-200/80 shadow-2xs overflow-hidden">
        <div className="p-6 border-b border-slate-100">
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
            <div>
              <div className="flex items-center gap-2.5">
                <h2 className="text-xl font-bold text-slate-900">Pusat Peringatan & Perhatian Khusus</h2>
                {totalAlertsCount > 0 ? (
                  <span className="px-2.5 py-0.5 rounded-full text-xs font-bold bg-red-100 text-red-700 animate-pulse">
                    {totalAlertsCount} Tindakan Diperlukan
                  </span>
                ) : (
                  <span className="px-2.5 py-0.5 rounded-full text-xs font-semibold bg-emerald-100 text-emerald-700">
                    Semua Terkendali
                  </span>
                )}
              </div>
              <p className="text-xs text-slate-500 mt-1">
                Monitoring otomatis penyewa jatuh tempo, mendekati batas pembayaran, dan status masa kontrak.
              </p>
            </div>
          </div>

          {/* Filter Tabs */}
          <div className="flex flex-wrap gap-2 mt-5">
            <button
              onClick={() => setActiveTab('all')}
              className={`px-3.5 py-1.5 rounded-xl text-xs font-semibold transition-all flex items-center gap-1.5 ${
                activeTab === 'all'
                  ? 'bg-slate-900 text-white shadow-xs'
                  : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
              }`}
            >
              Semua Perhatian
              <span className={`px-1.5 py-0.2 rounded-full text-[10px] ${activeTab === 'all' ? 'bg-white/20 text-white' : 'bg-slate-200 text-slate-700'}`}>
                {totalAlertsCount}
              </span>
            </button>

            <button
              onClick={() => setActiveTab('overdue')}
              className={`px-3.5 py-1.5 rounded-xl text-xs font-semibold transition-all flex items-center gap-1.5 ${
                activeTab === 'overdue'
                  ? 'bg-red-600 text-white shadow-xs'
                  : 'bg-red-50 text-red-700 hover:bg-red-100 border border-red-200/50'
              }`}
            >
              <AlertCircle size={14} />
              Jatuh Tempo (Overdue)
              <span className={`px-1.5 py-0.2 rounded-full text-[10px] ${activeTab === 'overdue' ? 'bg-white/20 text-white' : 'bg-red-200 text-red-800'}`}>
                {overdue.length}
              </span>
            </button>

            <button
              onClick={() => setActiveTab('due_soon')}
              className={`px-3.5 py-1.5 rounded-xl text-xs font-semibold transition-all flex items-center gap-1.5 ${
                activeTab === 'due_soon'
                  ? 'bg-amber-600 text-white shadow-xs'
                  : 'bg-amber-50 text-amber-800 hover:bg-amber-100 border border-amber-200/50'
              }`}
            >
              <Clock size={14} />
              Segera Jatuh Tempo
              <span className={`px-1.5 py-0.2 rounded-full text-[10px] ${activeTab === 'due_soon' ? 'bg-white/20 text-white' : 'bg-amber-200 text-amber-800'}`}>
                {dueSoon.length}
              </span>
            </button>

            <button
              onClick={() => setActiveTab('expiring')}
              className={`px-3.5 py-1.5 rounded-xl text-xs font-semibold transition-all flex items-center gap-1.5 ${
                activeTab === 'expiring'
                  ? 'bg-orange text-white shadow-xs'
                  : 'bg-orange/10 text-orange hover:bg-orange/20 border border-orange/20'
              }`}
            >
              <AlertTriangle size={14} />
              Akan Habis Kontrak
              <span className={`px-1.5 py-0.2 rounded-full text-[10px] ${activeTab === 'expiring' ? 'bg-white/20 text-white' : 'bg-orange/20 text-orange'}`}>
                {expiring.length}
              </span>
            </button>

            <button
              onClick={() => setActiveTab('expired')}
              className={`px-3.5 py-1.5 rounded-xl text-xs font-semibold transition-all flex items-center gap-1.5 ${
                activeTab === 'expired'
                  ? 'bg-rose-800 text-white shadow-xs'
                  : 'bg-rose-50 text-rose-800 hover:bg-rose-100 border border-rose-200/50'
              }`}
            >
              <XCircle size={14} />
              Habis Kontrak
              <span className={`px-1.5 py-0.2 rounded-full text-[10px] ${activeTab === 'expired' ? 'bg-white/20 text-white' : 'bg-rose-200 text-rose-900'}`}>
                {expired.length}
              </span>
            </button>
          </div>
        </div>

        {/* Tab Content Cards */}
        <div className="p-6">
          {totalAlertsCount === 0 ? (
            <div className="text-center py-12">
              <CheckCircle2 size={42} className="text-emerald-500 mx-auto mb-3" />
              <h4 className="text-base font-semibold text-slate-800">Tidak Ada Peringatan Kritis</h4>
              <p className="text-xs text-slate-500 mt-1">Seluruh tagihan sewa dan masa kontrak tenant dalam status aman.</p>
            </div>
          ) : (
            <div className="space-y-3">
              {/* Overdue Items */}
              {(activeTab === 'all' || activeTab === 'overdue') &&
                overdue.map((item) => (
                  <div
                    key={`overdue-${item.id}`}
                    className="flex flex-col md:flex-row md:items-center justify-between p-4 rounded-xl border border-red-200 bg-red-50/40 hover:bg-red-50/70 transition-all gap-4"
                  >
                    <div className="flex items-start gap-3.5">
                      <div className="w-10 h-10 rounded-xl bg-red-100 text-red-700 flex items-center justify-center shrink-0 mt-0.5">
                        <AlertCircle size={20} />
                      </div>
                      <div>
                        <div className="flex items-center gap-2 flex-wrap">
                          <h4 className="font-bold text-slate-900 text-sm">{item.tenant_name}</h4>
                          <AlertBadge type="overdue" sublabel={`${item.days_overdue} hari lewat`} />
                          <span className="text-xs text-slate-500">
                            • {item.property_name} ({item.room_name})
                          </span>
                        </div>
                        <p className="text-xs text-slate-600 mt-1">
                          Nominal Tunggakan: <strong className="text-red-700">{formatRupiah(item.amount)}</strong> • Jatuh Tempo: {new Date(item.due_date).toLocaleDateString('id-ID')}
                        </p>
                      </div>
                    </div>

                    <div className="flex items-center gap-2 shrink-0 self-end md:self-center">
                      <button
                        onClick={() => handleSendReminder(item.tenant_phone, item.tenant_name, 'overdue', item.amount, item.due_date)}
                        className="px-3 py-1.5 rounded-lg bg-red-600 hover:bg-red-700 text-white text-xs font-semibold flex items-center gap-1.5 shadow-2xs transition-all"
                      >
                        <MessageSquare size={13} />
                        Kirim Peringatan WA
                      </button>
                      <button
                        onClick={() => navigate('/dashboard/invoices')}
                        className="px-3 py-1.5 rounded-lg bg-white border border-slate-300 text-slate-700 hover:bg-slate-50 text-xs font-semibold transition-all"
                      >
                        Lihat Tagihan
                      </button>
                    </div>
                  </div>
                ))}

              {/* Due Soon Items */}
              {(activeTab === 'all' || activeTab === 'due_soon') &&
                dueSoon.map((item) => (
                  <div
                    key={`due-${item.id}`}
                    className="flex flex-col md:flex-row md:items-center justify-between p-4 rounded-xl border border-amber-200 bg-amber-50/40 hover:bg-amber-50/70 transition-all gap-4"
                  >
                    <div className="flex items-start gap-3.5">
                      <div className="w-10 h-10 rounded-xl bg-amber-100 text-amber-700 flex items-center justify-center shrink-0 mt-0.5">
                        <Clock size={20} />
                      </div>
                      <div>
                        <div className="flex items-center gap-2 flex-wrap">
                          <h4 className="font-bold text-slate-900 text-sm">{item.tenant_name}</h4>
                          <AlertBadge type="due_soon" sublabel={`Sisa ${item.days_until_due} hari`} />
                          <span className="text-xs text-slate-500">
                            • {item.property_name} ({item.room_name})
                          </span>
                        </div>
                        <p className="text-xs text-slate-600 mt-1">
                          Nominal: <strong>{formatRupiah(item.amount)}</strong> • Jatuh Tempo: {new Date(item.due_date).toLocaleDateString('id-ID')}
                        </p>
                      </div>
                    </div>

                    <div className="flex items-center gap-2 shrink-0 self-end md:self-center">
                      <button
                        onClick={() => handleSendReminder(item.tenant_phone, item.tenant_name, 'due_soon', item.amount, item.due_date)}
                        className="px-3 py-1.5 rounded-lg bg-amber-600 hover:bg-amber-700 text-white text-xs font-semibold flex items-center gap-1.5 shadow-2xs transition-all"
                      >
                        <Send size={13} />
                        Kirim Reminder
                      </button>
                      <button
                        onClick={() => navigate('/dashboard/invoices')}
                        className="px-3 py-1.5 rounded-lg bg-white border border-slate-300 text-slate-700 hover:bg-slate-50 text-xs font-semibold transition-all"
                      >
                        Detail Invoice
                      </button>
                    </div>
                  </div>
                ))}

              {/* Expiring Contract Items */}
              {(activeTab === 'all' || activeTab === 'expiring') &&
                expiring.map((item) => (
                  <div
                    key={`expiring-${item.id}`}
                    className="flex flex-col md:flex-row md:items-center justify-between p-4 rounded-xl border border-orange/30 bg-orange/5 hover:bg-orange/10 transition-all gap-4"
                  >
                    <div className="flex items-start gap-3.5">
                      <div className="w-10 h-10 rounded-xl bg-orange/20 text-orange flex items-center justify-center shrink-0 mt-0.5">
                        <AlertTriangle size={20} />
                      </div>
                      <div>
                        <div className="flex items-center gap-2 flex-wrap">
                          <h4 className="font-bold text-slate-900 text-sm">{item.tenant_name}</h4>
                          <AlertBadge type="expiring_soon" sublabel={`Sisa ${item.days_until_expiration} hari`} />
                          <span className="text-xs text-slate-500">
                            • {item.property_name} ({item.room_name})
                          </span>
                        </div>
                        <p className="text-xs text-slate-600 mt-1">
                          Sewa: {formatRupiah(item.monthly_rent)}/bln • Masa Sewa Habis: {new Date(item.end_date).toLocaleDateString('id-ID')}
                        </p>
                      </div>
                    </div>

                    <div className="flex items-center gap-2 shrink-0 self-end md:self-center">
                      <button
                        onClick={() => handleSendReminder(item.tenant_phone, item.tenant_name, 'expiring', item.monthly_rent, item.end_date)}
                        className="px-3 py-1.5 rounded-lg bg-orange hover:bg-orange/90 text-white text-xs font-semibold flex items-center gap-1.5 shadow-2xs transition-all"
                      >
                        <MessageSquare size={13} />
                        Tawarkan Perpanjang
                      </button>
                      <button
                        onClick={() => navigate('/dashboard/contracts')}
                        className="px-3 py-1.5 rounded-lg bg-white border border-slate-300 text-slate-700 hover:bg-slate-50 text-xs font-semibold transition-all"
                      >
                        Lihat Kontrak
                      </button>
                    </div>
                  </div>
                ))}

              {/* Expired Contract Items */}
              {(activeTab === 'all' || activeTab === 'expired') &&
                expired.map((item) => (
                  <div
                    key={`expired-${item.id}`}
                    className="flex flex-col md:flex-row md:items-center justify-between p-4 rounded-xl border border-rose-300 bg-rose-50/50 hover:bg-rose-50/80 transition-all gap-4"
                  >
                    <div className="flex items-start gap-3.5">
                      <div className="w-10 h-10 rounded-xl bg-rose-200 text-rose-800 flex items-center justify-center shrink-0 mt-0.5">
                        <XCircle size={20} />
                      </div>
                      <div>
                        <div className="flex items-center gap-2 flex-wrap">
                          <h4 className="font-bold text-slate-900 text-sm">{item.tenant_name}</h4>
                          <AlertBadge type="expired" sublabel={`Lewat ${item.days_expired} hari`} />
                          <span className="text-xs text-slate-500">
                            • {item.property_name} ({item.room_name})
                          </span>
                        </div>
                        <p className="text-xs text-slate-600 mt-1">
                          Kontrak Berakhir: {new Date(item.end_date).toLocaleDateString('id-ID')} • Perlu diproses check-out atau perpanjangan.
                        </p>
                      </div>
                    </div>

                    <div className="flex items-center gap-2 shrink-0 self-end md:self-center">
                      <button
                        onClick={() => navigate('/dashboard/occupancies')}
                        className="px-3 py-1.5 rounded-lg bg-rose-800 hover:bg-rose-900 text-white text-xs font-semibold flex items-center gap-1.5 shadow-2xs transition-all"
                      >
                        Proses Check-out
                      </button>
                      <button
                        onClick={() => navigate('/dashboard/contracts/new')}
                        className="px-3 py-1.5 rounded-lg bg-white border border-slate-300 text-slate-700 hover:bg-slate-50 text-xs font-semibold transition-all"
                      >
                        Buat Kontrak Baru
                      </button>
                    </div>
                  </div>
                ))}
            </div>
          )}
        </div>
      </div>

      {/* ─── DIREKTORI HUNIAN & ROSTER KAMAR TERPADU ─── */}
      <div className="space-y-3">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
          <div>
            <h3 className="text-lg font-bold text-slate-900 flex items-center gap-2">
              <DoorOpen className="w-5 h-5 text-orange" />
              Direktori Hunian & Status Kamar (Live Roster)
            </h3>
            <p className="text-xs text-slate-500">
              Pantau status sewa setiap kamar (Kosong, Dipesan, Terisi), relasi penyewa, masa kontrak, dan status tagihan dalam satu pandangan.
            </p>
          </div>
        </div>
        <TenantRoomRoster />
      </div>

      {/* Grid 2 Kolom: Okupansi Properti & Aktivitas Terbaru */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Okupansi Per Properti */}
        <div className="lg:col-span-2 bg-white p-6 rounded-2xl border border-slate-200/80 shadow-2xs">
          <div className="flex items-center justify-between mb-5">
            <div>
              <h3 className="text-base font-bold text-slate-900">Distribusi Okupansi Properti</h3>
              <p className="text-xs text-slate-500 mt-0.5">Performa okupansi unit per masing-masing properti cabang.</p>
            </div>
            <button
              onClick={() => navigate('/dashboard/properties')}
              className="text-xs font-semibold text-orange hover:text-orange/80 flex items-center gap-1"
            >
              Lihat Properti <ArrowRight size={13} />
            </button>
          </div>

          <div className="space-y-4">
            {data?.property_occupancy && data.property_occupancy.length > 0 ? (
              data.property_occupancy.map((prop) => (
                <div key={prop.property_id} className="p-3.5 rounded-xl border border-slate-100 hover:bg-slate-50/60 transition-all">
                  <div className="flex items-center justify-between text-sm mb-1.5">
                    <div className="flex items-center gap-2">
                      <Building2 size={16} className="text-slate-400" />
                      <span className="font-semibold text-slate-900">{prop.property_name}</span>
                    </div>
                    <span className="text-xs font-bold text-slate-700">{prop.rate}% Terisi</span>
                  </div>
                  <div className="w-full bg-slate-100 rounded-full h-2 overflow-hidden">
                    <div
                      className={`h-2 rounded-full transition-all duration-800 ${
                        prop.rate >= 80 ? 'bg-emerald-500' : prop.rate >= 50 ? 'bg-amber-500' : 'bg-orange'
                      }`}
                      style={{ width: `${Math.min(100, prop.rate)}%` }}
                    />
                  </div>
                  <div className="flex justify-between text-[11px] text-slate-400 mt-1.5 font-medium">
                    <span>{prop.occupied} Unit Terpakai</span>
                    <span>Total {prop.total_rooms} Unit Kamar</span>
                  </div>
                </div>
              ))
            ) : (
              <p className="text-xs text-slate-400 text-center py-6">Belum ada data properti yang tercatat.</p>
            )}
          </div>
        </div>

        {/* Aktivitas Terkini Feed */}
        <div className="bg-white p-6 rounded-2xl border border-slate-200/80 shadow-2xs">
          <h3 className="text-base font-bold text-slate-900 mb-1">Aktivitas Operasional</h3>
          <p className="text-xs text-slate-500 mb-4">Feed transaksi dan check-in terakhir.</p>

          <div className="space-y-3">
            {data?.recent_activities && data.recent_activities.length > 0 ? (
              data.recent_activities.map((act) => (
                <div key={act.id} className="flex items-start gap-3 p-2.5 rounded-xl hover:bg-slate-50 transition-colors">
                  <div
                    className={`w-7 h-7 rounded-lg flex items-center justify-center shrink-0 text-xs font-bold ${
                      act.type === 'payment'
                        ? 'bg-emerald-100 text-emerald-700'
                        : act.type === 'occupancy'
                        ? 'bg-purple-100 text-purple-700'
                        : 'bg-amber-100 text-amber-700'
                    }`}
                  >
                    {act.type === 'payment' ? 'Rp' : act.type === 'occupancy' ? 'IN' : 'WO'}
                  </div>
                  <div className="min-w-0 flex-1">
                    <p className="text-xs font-semibold text-slate-800 truncate">{act.title}</p>
                    <p className="text-[11px] text-slate-500 truncate">{act.description}</p>
                    <span className="text-[10px] text-slate-400">
                      {new Date(act.timestamp).toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })} • {new Date(act.timestamp).toLocaleDateString('id-ID')}
                    </span>
                  </div>
                </div>
              ))
            ) : (
              <p className="text-xs text-slate-400 text-center py-6">Belum ada aktivitas tercatat.</p>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

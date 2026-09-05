import { useState, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  DoorOpen,
  Search,
  MessageSquare,
  FileText,
  Building2,
  Phone,
  Eye,
  Receipt,
  CreditCard
} from 'lucide-react';
import { useRooms } from '@/features/room/hooks';
import { useTenants } from '@/features/tenant/hooks';
import { useContracts } from '@/features/contract/hooks';
import { useOccupancys } from '@/features/occupancy/hooks';
import { useInvoices } from '@/features/billing/hooks';
import { usePropertys } from '@/features/property/hooks';
import { formatCurrency } from '@/utils/currency';
import { AlertBadge } from '@/components/ui/AlertBadge';
import { Input } from '@/components/ui/input';

export function TenantRoomRoster() {
  const navigate = useNavigate();
  const [filterStatus, setFilterStatus] = useState<string>('all');
  const [searchQuery, setSearchQuery] = useState('');

  const { data: roomsData, isLoading: loadingRooms } = useRooms({ per_page: 100 });
  const { data: tenantsData } = useTenants({ per_page: 100 });
  const { data: contractsData } = useContracts({ per_page: 100 });
  const { data: occupanciesData } = useOccupancys({ per_page: 100 });
  const { data: invoicesData } = useInvoices({ per_page: 100 });
  const { data: propertiesData } = usePropertys({ per_page: 100 });

  const rooms = useMemo(() => Array.isArray(roomsData?.data) ? roomsData.data : [], [roomsData]);
  const tenants = useMemo(() => Array.isArray(tenantsData?.data) ? tenantsData.data : [], [tenantsData]);
  const contracts = useMemo(() => Array.isArray(contractsData?.data) ? contractsData.data : [], [contractsData]);
  const occupancies = useMemo(() => Array.isArray(occupanciesData?.data) ? occupanciesData.data : [], [occupanciesData]);
  const invoices = useMemo(() => Array.isArray(invoicesData?.data) ? invoicesData.data : [], [invoicesData]);
  const properties = useMemo(() => Array.isArray(propertiesData?.data) ? propertiesData.data : [], [propertiesData]);

  const propertiesMap = useMemo(() => new Map(properties.map(p => [p.id, p.name])), [properties]);
  const tenantsMap = useMemo(() => new Map(tenants.map(t => [t.id, t])), [tenants]);

  // Combine unified roster rows per room
  const rosterItems = useMemo(() => {
    return rooms.map((room) => {
      const propertyName = propertiesMap.get(room.property_id) || 'Properti';

      // Find active contract for this room
      const activeContract = contracts.find(
        (c) => c.room_id === room.id && (c.status === 'Active' || c.status === 'Pending')
      );

      // Find active checked-in occupancy for this room
      const activeOccupancy = occupancies.find(
        (o) => o.room_id === room.id && o.status === 'CheckedIn'
      );

      // Find tenant
      const tenantId = activeOccupancy?.tenant_id || activeContract?.tenant_id;
      const tenant = tenantId ? tenantsMap.get(tenantId) : null;

      // Determine precise operational status
      let roomStatus: 'Occupied' | 'Reserved' | 'Available' | 'Maintenance' = 'Available';
      if (room.status === 'Maintenance') {
        roomStatus = 'Maintenance';
      } else if (activeOccupancy) {
        roomStatus = 'Occupied';
      } else if (activeContract) {
        roomStatus = 'Reserved';
      } else if (room.status === 'Occupied') {
        roomStatus = 'Occupied';
      } else if (room.status === 'Reserved') {
        roomStatus = 'Reserved';
      } else if (!room.is_available) {
        roomStatus = 'Occupied';
      }

      // Find latest invoice for this tenant/contract
      const roomInvoices = activeContract
        ? invoices.filter((i) => i.contract_id === activeContract.id)
        : [];
      const latestInvoice = roomInvoices.length > 0
        ? roomInvoices.sort((a, b) => new Date(b.created_at || '').getTime() - new Date(a.created_at || '').getTime())[0]
        : null;

      // Calculate remaining days for contract
      let daysRemaining: number | null = null;
      let contractEndStr = '';
      if (activeContract?.end_date) {
        const end = new Date(activeContract.end_date);
        const now = new Date();
        daysRemaining = Math.ceil((end.getTime() - now.getTime()) / (1000 * 60 * 60 * 24));
        contractEndStr = end.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' });
      }

      return {
        room,
        propertyName,
        roomStatus,
        tenant,
        activeContract,
        activeOccupancy,
        latestInvoice,
        daysRemaining,
        contractEndStr,
      };
    });
  }, [rooms, propertiesMap, contracts, occupancies, tenantsMap, invoices]);

  // Filter items based on tab and search query
  const filteredItems = useMemo(() => {
    return rosterItems.filter((item) => {
      if (filterStatus !== 'all' && item.roomStatus !== filterStatus) {
        return false;
      }
      if (searchQuery.trim()) {
        const q = searchQuery.toLowerCase();
        const matchRoom = item.room.name.toLowerCase().includes(q);
        const matchProp = item.propertyName.toLowerCase().includes(q);
        const matchTenant = item.tenant?.full_name?.toLowerCase().includes(q) || false;
        const matchPhone = item.tenant?.phone?.includes(q) || false;
        return matchRoom || matchProp || matchTenant || matchPhone;
      }
      return true;
    });
  }, [rosterItems, filterStatus, searchQuery]);

  const stats = useMemo(() => {
    const total = rosterItems.length;
    const occupied = rosterItems.filter(i => i.roomStatus === 'Occupied').length;
    const reserved = rosterItems.filter(i => i.roomStatus === 'Reserved').length;
    const available = rosterItems.filter(i => i.roomStatus === 'Available').length;
    const maintenance = rosterItems.filter(i => i.roomStatus === 'Maintenance').length;
    return { total, occupied, reserved, available, maintenance };
  }, [rosterItems]);

  const openWhatsApp = (phone?: string, tenantName?: string, roomName?: string) => {
    if (!phone) return;
    const cleanPhone = phone.replace(/[^0-9]/g, '');
    const intlPhone = cleanPhone.startsWith('0') ? '62' + cleanPhone.slice(1) : cleanPhone;
    const text = encodeURIComponent(`Halo ${tenantName || 'Penyewa'}, ini dari pengelola mengenai kamar/unit ${roomName || ''}. Ada yang dapat kami bantu?`);
    window.open(`https://wa.me/${intlPhone}?text=${text}`, '_blank');
  };

  return (
    <div className="space-y-4">
      {/* Overview Metric Bar */}
      <div className="grid grid-cols-2 sm:grid-cols-5 gap-3">
        <div
          onClick={() => setFilterStatus('all')}
          className={`cursor-pointer p-4 rounded-2xl border transition-all ${
            filterStatus === 'all'
              ? 'bg-slate-900 text-white border-slate-900 shadow-md'
              : 'bg-white text-slate-700 border-slate-200 hover:border-slate-300'
          }`}
        >
          <p className="text-xs font-medium opacity-75">Total Kamar</p>
          <p className="text-2xl font-bold mt-0.5">{stats.total}</p>
        </div>

        <div
          onClick={() => setFilterStatus('Occupied')}
          className={`cursor-pointer p-4 rounded-2xl border transition-all ${
            filterStatus === 'Occupied'
              ? 'bg-red-600 text-white border-red-600 shadow-md'
              : 'bg-red-50/50 text-red-900 border-red-200 hover:border-red-300'
          }`}
        >
          <p className="text-xs font-semibold flex items-center gap-1.5 text-red-600">
            <span className="w-2 h-2 rounded-full bg-red-500 animate-pulse"></span>
            Terisi (Occupied)
          </p>
          <p className="text-2xl font-bold mt-0.5">{stats.occupied}</p>
        </div>

        <div
          onClick={() => setFilterStatus('Reserved')}
          className={`cursor-pointer p-4 rounded-2xl border transition-all ${
            filterStatus === 'Reserved'
              ? 'bg-amber-600 text-white border-amber-600 shadow-md'
              : 'bg-amber-50/50 text-amber-900 border-amber-200 hover:border-amber-300'
          }`}
        >
          <p className="text-xs font-semibold flex items-center gap-1.5 text-amber-600">
            <span className="w-2 h-2 rounded-full bg-amber-500"></span>
            Dipesan (Belum Masuk)
          </p>
          <p className="text-2xl font-bold mt-0.5">{stats.reserved}</p>
        </div>

        <div
          onClick={() => setFilterStatus('Available')}
          className={`cursor-pointer p-4 rounded-2xl border transition-all ${
            filterStatus === 'Available'
              ? 'bg-emerald-600 text-white border-emerald-600 shadow-md'
              : 'bg-emerald-50/50 text-emerald-900 border-emerald-200 hover:border-emerald-300'
          }`}
        >
          <p className="text-xs font-semibold flex items-center gap-1.5 text-emerald-600">
            <span className="w-2 h-2 rounded-full bg-emerald-500"></span>
            Tersedia (Kosong)
          </p>
          <p className="text-2xl font-bold mt-0.5">{stats.available}</p>
        </div>

        <div
          onClick={() => setFilterStatus('Maintenance')}
          className={`cursor-pointer p-4 rounded-2xl border transition-all ${
            filterStatus === 'Maintenance'
              ? 'bg-slate-700 text-white border-slate-700 shadow-md'
              : 'bg-slate-100/70 text-slate-700 border-slate-300 hover:border-slate-400'
          }`}
        >
          <p className="text-xs font-medium opacity-75">Dalam Perbaikan</p>
          <p className="text-2xl font-bold mt-0.5">{stats.maintenance}</p>
        </div>
      </div>

      {/* Filter and Search Bar */}
      <div className="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3 bg-white p-3.5 rounded-2xl border border-slate-200 shadow-sm">
        <div className="relative flex-1 max-w-md">
          <Search className="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
          <Input
            placeholder="Cari nomor kamar, nama tenant, atau telepon..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="pl-9 bg-slate-50 border-slate-200 focus:bg-white"
          />
        </div>

        <div className="flex items-center gap-2 overflow-x-auto text-xs">
          <span className="text-slate-400 font-medium whitespace-nowrap">Filter Status:</span>
          {[
            { id: 'all', label: 'Semua' },
            { id: 'Occupied', label: '🔴 Terisi' },
            { id: 'Reserved', label: '🟡 Dipesan' },
            { id: 'Available', label: '🟢 Tersedia' },
            { id: 'Maintenance', label: '⚪ Perbaikan' },
          ].map((tab) => (
            <button
              key={tab.id}
              onClick={() => setFilterStatus(tab.id)}
              className={`px-3 py-1.5 rounded-xl font-semibold transition-colors whitespace-nowrap ${
                filterStatus === tab.id
                  ? 'bg-slate-900 text-white shadow-sm'
                  : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
              }`}
            >
              {tab.label}
            </button>
          ))}
        </div>
      </div>

      {/* Roster Table */}
      <div className="rounded-2xl border border-slate-200 bg-white overflow-hidden shadow-sm">
        {loadingRooms ? (
          <div className="p-12 text-center text-slate-400">
            Memuat direktori kamar dan hunian...
          </div>
        ) : filteredItems.length === 0 ? (
          <div className="p-12 text-center text-slate-500">
            Tidak ada kamar yang sesuai dengan filter pencarian.
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm text-left">
              <thead className="border-b border-slate-100 bg-slate-50/75 text-slate-600 font-semibold text-xs uppercase tracking-wider">
                <tr>
                  <th className="px-5 py-3.5">Kamar & Properti</th>
                  <th className="px-5 py-3.5">Status Unit</th>
                  <th className="px-5 py-3.5">Penyewa (Tenant)</th>
                  <th className="px-5 py-3.5">Periode Sewa & Sisa Kontrak</th>
                  <th className="px-5 py-3.5">Status Tagihan Terakhir</th>
                  <th className="px-5 py-3.5 text-right">Aksi</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100">
                {filteredItems.map((item) => (
                  <tr
                    key={item.room.id}
                    className="hover:bg-slate-50/60 transition-colors cursor-pointer"
                    onClick={() => navigate(`/dashboard/rooms/${item.room.id}`)}
                  >
                    {/* Kamar & Properti */}
                    <td className="px-5 py-4">
                      <div className="flex items-center gap-3">
                        <div className="w-9 h-9 rounded-xl bg-orange/10 flex items-center justify-center text-orange font-bold shrink-0">
                          <DoorOpen className="w-4 h-4" />
                        </div>
                        <div>
                          <p className="font-bold text-slate-900">{item.room.name}</p>
                          <p className="text-xs text-slate-400 flex items-center gap-1">
                            <Building2 className="w-3 h-3" />
                            {item.propertyName} · {formatCurrency(item.room.price, item.room.currency || 'IDR')}
                          </p>
                        </div>
                      </div>
                    </td>

                    {/* Status Unit */}
                    <td className="px-5 py-4">
                      {item.roomStatus === 'Occupied' && (
                        <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold bg-red-50 text-red-700 border border-red-200">
                          <span className="w-1.5 h-1.5 rounded-full bg-red-500 animate-pulse"></span>
                          Terisi (Occupied)
                        </span>
                      )}
                      {item.roomStatus === 'Reserved' && (
                        <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold bg-amber-50 text-amber-700 border border-amber-200">
                          <span className="w-1.5 h-1.5 rounded-full bg-amber-500"></span>
                          Dipesan (Belum Masuk)
                        </span>
                      )}
                      {item.roomStatus === 'Available' && (
                        <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold bg-emerald-50 text-emerald-700 border border-emerald-200">
                          <span className="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
                          Tersedia (Kosong)
                        </span>
                      )}
                      {item.roomStatus === 'Maintenance' && (
                        <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold bg-slate-100 text-slate-700 border border-slate-200">
                          <span className="w-1.5 h-1.5 rounded-full bg-slate-500"></span>
                          Perbaikan
                        </span>
                      )}
                    </td>

                    {/* Tenant Info */}
                    <td className="px-5 py-4">
                      {item.tenant ? (
                        <div>
                          <p className="font-semibold text-slate-900">{item.tenant.full_name}</p>
                          <div className="flex items-center gap-2 text-xs text-slate-500 mt-0.5">
                            <span className="flex items-center gap-1">
                              <Phone className="w-3 h-3 text-slate-400" />
                              {item.tenant.phone || '—'}
                            </span>
                          </div>
                        </div>
                      ) : (
                        <span className="text-xs text-slate-400 italic">Belum ada penyewa</span>
                      )}
                    </td>

                    {/* Masa Sewa */}
                    <td className="px-5 py-4">
                      {item.activeContract ? (
                        <div>
                          <p className="text-xs text-slate-700 font-medium">
                            s/d {item.contractEndStr}
                          </p>
                          {item.daysRemaining !== null && (
                            <p className="text-xs mt-0.5">
                              {item.daysRemaining < 0 ? (
                                <span className="text-red-600 font-semibold">Habis Kontrak</span>
                              ) : item.daysRemaining <= 30 ? (
                                <span className="text-amber-600 font-semibold">Sisa {item.daysRemaining} hari</span>
                              ) : (
                                <span className="text-slate-500">Sisa {Math.round(item.daysRemaining / 30)} bulan</span>
                              )}
                            </p>
                          )}
                        </div>
                      ) : (
                        <span className="text-xs text-slate-400">—</span>
                      )}
                    </td>

                    {/* Status Tagihan */}
                    <td className="px-5 py-4">
                      {item.latestInvoice ? (
                        <div>
                          {item.latestInvoice.status === 'Paid' ? (
                            <AlertBadge variant="paid" label="Lunas" size="sm" />
                          ) : (
                            <AlertBadge
                              variant="pending"
                              label={`Tagihan ${formatCurrency(item.latestInvoice.amount, (item.latestInvoice as any).currency || 'IDR')}`}
                              size="sm"
                            />
                          )}
                        </div>
                      ) : item.tenant ? (
                        <span className="text-xs text-slate-400">Tidak ada tagihan aktif</span>
                      ) : (
                        <span className="text-xs text-slate-400">—</span>
                      )}
                    </td>

                    {/* Aksi */}
                    <td className="px-5 py-4 text-right" onClick={(e) => e.stopPropagation()}>
                      <div className="flex items-center justify-end gap-1.5">
                        {/* Aksi Cepat Tagihan / Pembayaran */}
                        {item.latestInvoice && item.latestInvoice.status !== 'Paid' ? (
                          <button
                            onClick={() =>
                              navigate(
                                `/payments/new?invoice_id=${item.latestInvoice!.id}&tenant_id=${item.tenant?.id}&amount=${item.latestInvoice!.amount}`
                              )
                            }
                            title="Bayar Tagihan Sekarang"
                            className="p-2 rounded-xl text-emerald-700 bg-emerald-50 hover:bg-emerald-100 border border-emerald-200 transition-colors"
                          >
                            <CreditCard className="w-3.5 h-3.5" />
                          </button>
                        ) : item.activeContract ? (
                          <button
                            onClick={() =>
                              navigate(
                                `/invoices/new?contract_id=${item.activeContract!.id}&tenant_id=${item.tenant?.id}&amount=${(item.activeContract as any)?.monthly_rent || 0}`
                              )
                            }
                            title="Buat Tagihan Baru"
                            className="p-2 rounded-xl text-orange bg-orange/5 hover:bg-orange/10 border border-orange/20 transition-colors"
                          >
                            <Receipt className="w-3.5 h-3.5" />
                          </button>
                        ) : null}

                        {item.tenant?.phone && (
                          <button
                            onClick={() => openWhatsApp(item.tenant!.phone, item.tenant!.full_name, item.room.name)}
                            title="Chat via WhatsApp"
                            className="p-2 rounded-xl text-emerald-600 hover:bg-emerald-50 border border-emerald-200 transition-colors"
                          >
                            <MessageSquare className="w-3.5 h-3.5" />
                          </button>
                        )}
                        {item.activeContract && (
                          <button
                            onClick={() => navigate(`/dashboard/contracts/${item.activeContract!.id}`)}
                            title="Buka Surat Kontrak"
                            className="p-2 rounded-xl text-slate-600 hover:bg-slate-100 border border-slate-200 transition-colors"
                          >
                            <FileText className="w-3.5 h-3.5" />
                          </button>
                        )}
                        <button
                          onClick={() => navigate(`/dashboard/rooms/${item.room.id}`)}
                          title="Detail Kamar"
                          className="p-2 rounded-xl text-slate-600 hover:bg-slate-100 border border-slate-200 transition-colors"
                        >
                          <Eye className="w-3.5 h-3.5" />
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
}

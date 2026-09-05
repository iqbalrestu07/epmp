import { useMemo } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { useRoom, useDeleteRoom } from "../hooks";
import { useFloors } from "../../floor/hooks";
import { useBuildings } from "../../building/hooks";
import { usePropertys } from "../../property/hooks";
import { useBeds, useDeleteBed } from "../../bed/hooks";
import { useOccupancys } from "../../occupancy/hooks";
import { useContracts } from "../../contract/hooks";
import { useTenants } from "../../tenant/hooks";
import { formatCurrency } from "@/utils/currency";
import { formatDate } from "@/utils/date";
import {
  DoorOpen,
  ChevronRight,
  Pencil,
  Trash2,
  ArrowLeft,
  Users,
  DollarSign,
  Layers,
  Building2,
  Bed as BedIcon,
  Plus,
  MessageSquare,
  Phone,
  Mail,
  Receipt,
  CreditCard,
  FileText,
  UserCheck,
  Clock,
  Sparkles,
  ExternalLink,
} from "lucide-react";

export function RoomDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { data, isLoading } = useRoom(id);
  const deleteMutation = useDeleteRoom();
  const deleteBedMutation = useDeleteBed();

  const { data: floorsData } = useFloors({ per_page: 100 });
  const { data: buildingsData } = useBuildings({ per_page: 100 });
  const { data: propertiesData } = usePropertys({ per_page: 100 });
  const { data: bedsData } = useBeds({ per_page: 100 });
  const { data: occupanciesData } = useOccupancys({ per_page: 100 });
  const { data: contractsData } = useContracts({ per_page: 100 });
  const { data: tenantsData } = useTenants({ per_page: 100 });

  const floors = useMemo(() => {
    return Array.isArray(floorsData?.data) ? floorsData.data : Array.isArray(floorsData) ? floorsData : [];
  }, [floorsData]);

  const buildings = useMemo(() => {
    return Array.isArray(buildingsData?.data) ? buildingsData.data : Array.isArray(buildingsData) ? buildingsData : [];
  }, [buildingsData]);

  const properties = useMemo(() => {
    return Array.isArray(propertiesData?.data) ? propertiesData.data : Array.isArray(propertiesData) ? propertiesData : [];
  }, [propertiesData]);

  const allBeds = useMemo(() => {
    return Array.isArray(bedsData?.data) ? bedsData.data : Array.isArray(bedsData) ? bedsData : [];
  }, [bedsData]);

  const occupancies = useMemo(() => {
    return Array.isArray(occupanciesData?.data) ? occupanciesData.data : Array.isArray(occupanciesData) ? occupanciesData : [];
  }, [occupanciesData]);

  const contracts = useMemo(() => {
    return Array.isArray(contractsData?.data) ? contractsData.data : Array.isArray(contractsData) ? contractsData : [];
  }, [contractsData]);

  const tenants = useMemo(() => {
    return Array.isArray(tenantsData?.data) ? tenantsData.data : Array.isArray(tenantsData) ? tenantsData : [];
  }, [tenantsData]);

  const parentFloor = floors.find((f: any) => f.id === data?.floor_id);
  const parentBuilding = parentFloor
    ? buildings.find((b: any) => b.id === parentFloor.building_id)
    : undefined;
  const parentProperty = properties.find((p: any) => p.id === data?.property_id);

  // Filter beds for this room
  const roomBeds = useMemo(() => {
    if (!id) return [];
    return allBeds.filter((b: any) => b.room_id === id);
  }, [allBeds, id]);

  // Find active occupancy and contract for this room
  const activeOccupancy = useMemo(() => {
    if (!id) return null;
    return (
      occupancies.find(
        (o: any) => o.room_id === id && (o.status === "CheckedIn" || o.status === "Active")
      ) || null
    );
  }, [occupancies, id]);

  const activeContract = useMemo(() => {
    if (!id) return null;
    return (
      contracts.find(
        (c: any) => c.room_id === id && (c.status === "Active" || c.status === "Signed")
      ) || null
    );
  }, [contracts, id]);

  const activeTenantId = activeOccupancy?.tenant_id || activeContract?.tenant_id;
  const activeTenant = useMemo(() => {
    if (!activeTenantId) return null;
    return tenants.find((t: any) => t.id === activeTenantId) || null;
  }, [tenants, activeTenantId]);

  // Operational room status
  const isOccupied = !!(activeOccupancy || activeContract || (data && !data.is_available) || data?.status === "Occupied");
  const isReserved = !isOccupied && (data?.status === "Reserved" || (activeContract && activeContract.status === "Signed" && !activeOccupancy));

  // Remaining days for contract
  const contractDaysRemaining = useMemo(() => {
    if (!activeContract?.end_date) return null;
    const end = new Date(activeContract.end_date);
    const now = new Date();
    return Math.ceil((end.getTime() - now.getTime()) / (1000 * 60 * 60 * 24));
  }, [activeContract]);

  const handleDelete = () => {
    if (!confirm(`Hapus kamar "${data?.name}"?`)) return;
    deleteMutation.mutate(id!, {
      onSuccess: () => navigate("/dashboard/rooms"),
    });
  };

  const handleDeleteBed = (bedId: string, bedName: string) => {
    if (!confirm(`Hapus tempat tidur "${bedName}" dari kamar ini?`)) return;
    deleteBedMutation.mutate(bedId);
  };

  const openWhatsApp = (phone?: string, tenantName?: string) => {
    if (!phone) return;
    let clean = phone.replace(/[^0-9]/g, "");
    if (clean.startsWith("0")) {
      clean = "62" + clean.slice(1);
    }
    const msg = encodeURIComponent(
      `Halo Kak ${tenantName || ""}, kami dari pengelola properti terkait unit kamar ${data?.name || ""}...`
    );
    window.open(`https://wa.me/${clean}?text=${msg}`, "_blank");
  };

  if (isLoading) return <div className="text-center py-12 text-slate-400">Memuat data kamar…</div>;
  if (!data) return <div className="text-center py-12 text-slate-400">Kamar tidak ditemukan</div>;

  return (
    <div className="space-y-6 pb-12">
      {/* Breadcrumbs */}
      <div className="flex items-center gap-2 text-sm text-slate-500 flex-wrap">
        <Link to="/dashboard/properties" className="hover:text-slate-900 transition-colors">Properti</Link>
        <ChevronRight size={14} />
        {parentProperty ? (
          <>
            <Link to={`/dashboard/properties/${parentProperty.id}`} className="hover:text-slate-900 transition-colors">
              {parentProperty.name}
            </Link>
            <ChevronRight size={14} />
          </>
        ) : null}
        {parentBuilding ? (
          <>
            <Link to={`/dashboard/buildings/${parentBuilding.id}`} className="hover:text-slate-900 transition-colors">
              {parentBuilding.name}
            </Link>
            <ChevronRight size={14} />
          </>
        ) : null}
        {parentFloor ? (
          <>
            <Link to={`/dashboard/floors/${parentFloor.id}`} className="hover:text-slate-900 transition-colors">
              {parentFloor.name}
            </Link>
            <ChevronRight size={14} />
          </>
        ) : null}
        <Link to="/dashboard/rooms" className="hover:text-slate-900 transition-colors">Kamar</Link>
        <ChevronRight size={14} />
        <span className="text-slate-900 font-semibold truncate max-w-40">{data.name}</span>
      </div>

      {/* Main Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 bg-white p-6 rounded-2xl border border-slate-200 shadow-sm">
        <div className="flex items-center gap-4">
          <div className="w-14 h-14 rounded-2xl bg-orange/10 border border-orange/20 flex items-center justify-center text-orange shrink-0">
            <DoorOpen size={28} />
          </div>
          <div>
            <div className="flex items-center gap-3">
              <h1 className="text-2xl font-bold text-slate-900">{data.name}</h1>
              <span
                className={`inline-flex items-center px-3 py-1 rounded-full text-xs font-semibold ${
                  isOccupied
                    ? "bg-rose-50 text-rose-700 border border-rose-200"
                    : isReserved
                    ? "bg-amber-50 text-amber-700 border border-amber-200"
                    : "bg-emerald-50 text-emerald-700 border border-emerald-200"
                }`}
              >
                {isOccupied ? "Terisi (Occupied)" : isReserved ? "Dipesan (Reserved)" : "Tersedia (Available)"}
              </span>
            </div>
            <p className="text-sm text-slate-500 mt-1 flex items-center gap-2 flex-wrap">
              <span className="inline-flex items-center gap-1 font-medium text-slate-700">
                <Users size={14} className="text-slate-400" /> Kapasitas: {data.capacity} Orang
              </span>
              <span>•</span>
              <span className="inline-flex items-center gap-1 font-semibold text-orange">
                <DollarSign size={14} /> {formatCurrency(data.price || 0, "IDR")}/bulan
              </span>
              {parentProperty && (
                <>
                  <span>•</span>
                  <span>{parentProperty.name}</span>
                </>
              )}
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" onClick={() => navigate("/dashboard/rooms")}>
            <ArrowLeft size={16} className="mr-1" />
            Kembali
          </Button>
          <Button variant="outline" size="sm" onClick={() => navigate(`/dashboard/rooms/${id}/edit`)}>
            <Pencil size={16} className="mr-1" />
            Edit Kamar
          </Button>
          <Button variant="destructive" size="sm" onClick={handleDelete}>
            <Trash2 size={16} className="mr-1" />
            Hapus
          </Button>
        </div>
      </div>

      {/* ── SECTION 1: INFORMASI PENGHUNI (TENANT) & KONTRAK SEWA ── */}
      <div className="bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden">
        <div className="px-6 py-4 border-b border-slate-100 flex items-center justify-between bg-slate-50/50">
          <div className="flex items-center gap-2.5">
            <UserCheck className="w-5 h-5 text-orange" />
            <h2 className="text-base font-bold text-slate-900">Status Penghuni & Kontrak Sewa</h2>
          </div>
          {isOccupied ? (
            <span className="px-3 py-1 rounded-full text-xs font-semibold bg-rose-50 text-rose-700 border border-rose-200">
              Kamar Sedang Dihuni
            </span>
          ) : (
            <span className="px-3 py-1 rounded-full text-xs font-semibold bg-emerald-50 text-emerald-700 border border-emerald-200">
              Kamar Siap Huni
            </span>
          )}
        </div>

        <div className="p-6">
          {activeTenant ? (
            <div className="space-y-6">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {/* Profil Penyewa */}
                <div className="p-5 rounded-xl border border-slate-200 bg-slate-50/60 flex flex-col justify-between space-y-4">
                  <div>
                    <div className="flex items-center justify-between mb-3">
                      <span className="text-xs font-bold uppercase tracking-wider text-slate-500">
                        Profil Penyewa Aktif
                      </span>
                      <Link
                        to={`/dashboard/tenants/${activeTenant.id}`}
                        className="text-xs font-medium text-orange hover:underline inline-flex items-center gap-1"
                      >
                        Detail Profil <ExternalLink size={12} />
                      </Link>
                    </div>

                    <div className="flex items-center gap-3.5">
                      <div className="w-12 h-12 rounded-full bg-orange/15 border border-orange/30 text-orange font-bold text-base flex items-center justify-center shrink-0">
                        {activeTenant.full_name?.slice(0, 2).toUpperCase() || "TN"}
                      </div>
                      <div className="min-w-0">
                        <h3 className="font-bold text-slate-900 text-lg leading-tight truncate">
                          {activeTenant.full_name}
                        </h3>
                        <p className="text-xs text-slate-500 mt-0.5">
                          ID: <span className="font-mono text-slate-600">#{activeTenant.id.slice(0, 8)}</span>
                        </p>
                      </div>
                    </div>

                    <div className="mt-4 space-y-2 text-xs text-slate-600">
                      {activeTenant.phone && (
                        <div className="flex items-center gap-2">
                          <Phone size={14} className="text-slate-400" />
                          <span className="font-medium">{activeTenant.phone}</span>
                        </div>
                      )}
                      {activeTenant.email && (
                        <div className="flex items-center gap-2">
                          <Mail size={14} className="text-slate-400" />
                          <span>{activeTenant.email}</span>
                        </div>
                      )}
                    </div>
                  </div>

                  {activeTenant.phone && (
                    <Button
                      onClick={() => openWhatsApp(activeTenant.phone, activeTenant.full_name)}
                      className="w-full bg-emerald-600 hover:bg-emerald-700 text-white font-medium text-xs py-2 flex items-center justify-center gap-2 shadow-sm"
                    >
                      <MessageSquare size={15} />
                      Chat WhatsApp Penyewa
                    </Button>
                  )}
                </div>

                {/* Detail Kontrak Sewa */}
                {activeContract ? (
                  <div className="p-5 rounded-xl border border-slate-200 bg-slate-50/60 flex flex-col justify-between space-y-4">
                    <div>
                      <div className="flex items-center justify-between mb-3">
                        <span className="text-xs font-bold uppercase tracking-wider text-slate-500">
                          Kontrak Sewa Aktif
                        </span>
                        <Link
                          to={`/dashboard/contracts/${activeContract.id}`}
                          className="text-xs font-medium text-orange hover:underline inline-flex items-center gap-1"
                        >
                          Buka Kontrak <ExternalLink size={12} />
                        </Link>
                      </div>

                      <div className="flex items-center justify-between">
                        <span className="font-mono text-sm font-semibold text-slate-900 bg-white px-2.5 py-1 rounded-md border border-slate-200">
                          #{activeContract.id.slice(0, 8)}
                        </span>
                        <span className="px-2 py-0.5 text-xs font-semibold rounded-md bg-emerald-100 text-emerald-800">
                          {activeContract.status}
                        </span>
                      </div>

                      <div className="mt-3 grid grid-cols-2 gap-3 text-xs">
                        <div className="bg-white p-2.5 rounded-lg border border-slate-200/80">
                          <span className="text-slate-400 block text-[11px]">Mulai Sewa</span>
                          <span className="font-medium text-slate-800">{formatDate(activeContract.start_date)}</span>
                        </div>
                        <div className="bg-white p-2.5 rounded-lg border border-slate-200/80">
                          <span className="text-slate-400 block text-[11px]">Selesai Sewa</span>
                          <span className="font-medium text-slate-800">{formatDate(activeContract.end_date)}</span>
                        </div>
                        <div className="bg-white p-2.5 rounded-lg border border-slate-200/80">
                          <span className="text-slate-400 block text-[11px]">Harga Sewa</span>
                          <span className="font-bold text-orange">
                            {formatCurrency(activeContract.monthly_rent || activeContract.price_per_month || 0, activeContract.currency || "IDR")}/bln
                          </span>
                        </div>
                        <div className="bg-white p-2.5 rounded-lg border border-slate-200/80">
                          <span className="text-slate-400 block text-[11px]">Deposit Jaminan</span>
                          <span className="font-medium text-slate-800">
                            {formatCurrency(activeContract.deposit_amount || 0, activeContract.currency || "IDR")}
                          </span>
                        </div>
                      </div>

                      {contractDaysRemaining !== null && (
                        <div className="mt-3 flex items-center gap-2 text-xs">
                          <Clock size={14} className={contractDaysRemaining <= 30 ? "text-rose-500" : "text-slate-400"} />
                          <span className={contractDaysRemaining <= 30 ? "text-rose-600 font-bold" : "text-slate-600"}>
                            {contractDaysRemaining > 0
                              ? `Masa sewa tersisa ${contractDaysRemaining} hari lagi`
                              : `Kontrak berakhir sejak ${Math.abs(contractDaysRemaining)} hari lalu`}
                          </span>
                        </div>
                      )}
                    </div>

                    <div className="flex gap-2">
                      <Button
                        size="sm"
                        variant="outline"
                        onClick={() =>
                          navigate(
                            `/invoices/new?contract_id=${activeContract.id}&tenant_id=${activeTenant.id}&amount=${activeContract.monthly_rent || activeContract.price_per_month || 0}`
                          )
                        }
                        className="flex-1 text-xs text-orange hover:bg-orange/5 border-orange/30"
                      >
                        <Receipt size={14} className="mr-1" />
                        Buat Tagihan
                      </Button>
                      <Button
                        size="sm"
                        variant="outline"
                        onClick={() =>
                          navigate(
                            `/deposits/new?contract_id=${activeContract.id}&tenant_id=${activeTenant.id}&amount=${activeContract.deposit_amount || 0}`
                          )
                        }
                        className="flex-1 text-xs text-slate-700"
                      >
                        <CreditCard size={14} className="mr-1" />
                        Catat Deposit
                      </Button>
                    </div>
                  </div>
                ) : (
                  <div className="p-5 rounded-xl border border-dashed border-slate-200 bg-slate-50/50 flex flex-col items-center justify-center text-center p-6 space-y-2">
                    <FileText className="w-8 h-8 text-slate-300" />
                    <p className="text-xs font-semibold text-slate-600">Belum Ada Kontrak Sewa Tertaut</p>
                    <Button
                      size="sm"
                      onClick={() => navigate(`/dashboard/contracts/new?room_id=${id}&tenant_id=${activeTenant.id}`)}
                      className="bg-orange hover:bg-orange/90 text-slate-900 text-xs mt-2"
                    >
                      <Plus size={14} className="mr-1" />
                      Buat Kontrak Sewa
                    </Button>
                  </div>
                )}
              </div>
            </div>
          ) : (
            <div className="py-8 text-center space-y-3">
              <div className="w-12 h-12 rounded-full bg-emerald-50 border border-emerald-200 text-emerald-600 flex items-center justify-center mx-auto">
                <Sparkles size={22} />
              </div>
              <div className="max-w-md mx-auto">
                <h3 className="text-sm font-bold text-slate-900">Kamar Ini Saat Ini Tersedia (Kosong)</h3>
                <p className="text-xs text-slate-500 mt-1">
                  Tidak ada penyewa atau kontrak aktif yang terdaftar di kamar ini. Anda dapat melakukan reservasi atau check-in tenant baru langsung.
                </p>
              </div>
              <div className="flex items-center justify-center gap-3 pt-2">
                <Button
                  size="sm"
                  onClick={() => navigate(`/dashboard/contracts/new?room_id=${id}`)}
                  className="bg-orange hover:bg-orange/90 text-slate-900 text-xs font-medium"
                >
                  <FileText size={14} className="mr-1.5" />
                  Buat Kontrak Sewa
                </Button>
                <Button
                  size="sm"
                  variant="outline"
                  onClick={() => navigate(`/dashboard/reservations/new?room_id=${id}`)}
                  className="text-xs"
                >
                  <Users size={14} className="mr-1.5" />
                  Reservasi Kamar
                </Button>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* ── SECTION 2: TEMPAT TIDUR DI KAMAR INI (BEDS LIST) ── */}
      <div className="bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden">
        <div className="px-6 py-4 border-b border-slate-100 flex items-center justify-between bg-slate-50/50">
          <div className="flex items-center gap-2.5">
            <BedIcon className="w-5 h-5 text-orange" />
            <div>
              <h2 className="text-base font-bold text-slate-900">Daftar Tempat Tidur (Beds)</h2>
              <p className="text-xs text-slate-500">
                Total {roomBeds.length} tempat tidur terdaftar di kamar ini
              </p>
            </div>
          </div>
          <Button
            size="sm"
            onClick={() => navigate(`/dashboard/beds/new?room_id=${id}`)}
            className="bg-orange hover:bg-orange/90 text-slate-900 text-xs font-semibold flex items-center gap-1.5"
          >
            <Plus size={14} />
            Tambah Tempat Tidur
          </Button>
        </div>

        <div className="p-6">
          {roomBeds.length > 0 ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {roomBeds.map((bed: any) => {
                const bedOccupied = isOccupied || bed.status === "Occupied";
                const bedMaintenance = bed.status === "Maintenance";

                return (
                  <div
                    key={bed.id}
                    className="p-4 rounded-xl border border-slate-200 bg-white hover:border-orange/40 hover:shadow-md transition-all flex flex-col justify-between space-y-3"
                  >
                    <div className="flex items-start justify-between">
                      <div className="flex items-center gap-3">
                        <div className="w-10 h-10 rounded-xl bg-orange/10 flex items-center justify-center text-orange shrink-0">
                          <BedIcon size={18} />
                        </div>
                        <div>
                          <h4 className="font-bold text-slate-900 text-sm">{bed.name}</h4>
                          <span className="font-mono text-[11px] text-slate-400">ID: #{bed.id.slice(0, 8)}</span>
                        </div>
                      </div>

                      <span
                        className={`inline-flex items-center px-2 py-0.5 rounded-full text-[11px] font-semibold ${
                          bedMaintenance
                            ? "bg-amber-50 text-amber-700 border border-amber-200"
                            : bedOccupied
                            ? "bg-rose-50 text-rose-700 border border-rose-200"
                            : "bg-emerald-50 text-emerald-700 border border-emerald-200"
                        }`}
                      >
                        {bedMaintenance ? "Maintenance" : bedOccupied ? "Occupied" : "Available"}
                      </span>
                    </div>

                    {bedOccupied && activeTenant && (
                      <div className="p-2.5 rounded-lg bg-rose-50/60 border border-rose-100 flex items-center gap-2 text-xs text-rose-900">
                        <Users size={13} className="text-rose-500 shrink-0" />
                        <span className="font-medium truncate">Dihuni oleh: {activeTenant.full_name}</span>
                      </div>
                    )}

                    <div className="flex items-center justify-end gap-2 pt-2 border-t border-slate-100">
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => navigate(`/dashboard/beds/${bed.id}/edit`)}
                        className="h-8 px-2.5 text-xs text-slate-600 hover:text-slate-900"
                      >
                        <Pencil size={13} className="mr-1" /> Edit
                      </Button>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleDeleteBed(bed.id, bed.name)}
                        className="h-8 px-2.5 text-xs text-rose-600 hover:text-rose-700 hover:bg-rose-50"
                      >
                        <Trash2 size={13} className="mr-1" /> Hapus
                      </Button>
                    </div>
                  </div>
                );
              })}
            </div>
          ) : (
            <div className="py-8 text-center border-2 border-dashed border-slate-200 rounded-xl space-y-2">
              <BedIcon className="w-8 h-8 text-slate-300 mx-auto" />
              <p className="text-sm font-semibold text-slate-700">Belum Ada Tempat Tidur di Kamar Ini</p>
              <p className="text-xs text-slate-400 max-w-sm mx-auto">
                Tambahkan tempat tidur (misal Single Bed, Queen Bed, King Bed) untuk unit kamar ini.
              </p>
              <Button
                size="sm"
                onClick={() => navigate(`/dashboard/beds/new?room_id=${id}`)}
                className="bg-orange hover:bg-orange/90 text-slate-900 text-xs mt-2"
              >
                <Plus size={14} className="mr-1" /> Tambah Bed Sekarang
              </Button>
            </div>
          )}
        </div>
      </div>

      {/* ── SECTION 3: SPESIFIKASI & HIERARKI PROPERTI ── */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Spesifikasi Teknis */}
        <div className="bg-white rounded-2xl border border-slate-200 shadow-sm p-6 space-y-4">
          <h3 className="font-bold text-slate-900 text-base">Spesifikasi Kamar</h3>
          <dl className="grid grid-cols-2 gap-4">
            <div className="p-3 rounded-xl bg-slate-50 border border-slate-100">
              <dt className="text-xs font-semibold text-slate-400 uppercase tracking-wider">Kapasitas Tamu</dt>
              <dd className="text-sm font-bold text-slate-900 mt-1">{data.capacity} Orang</dd>
            </div>
            <div className="p-3 rounded-xl bg-slate-50 border border-slate-100">
              <dt className="text-xs font-semibold text-slate-400 uppercase tracking-wider">Harga Patokan</dt>
              <dd className="text-sm font-bold text-orange mt-1">
                {formatCurrency(data.price || 0, "IDR")}/bln
              </dd>
            </div>
            <div className="p-3 rounded-xl bg-slate-50 border border-slate-100">
              <dt className="text-xs font-semibold text-slate-400 uppercase tracking-wider">Status Kamar</dt>
              <dd className="text-sm font-bold text-slate-900 mt-1">
                {isOccupied ? "Terisi" : isReserved ? "Dipesan" : "Tersedia"}
              </dd>
            </div>
            <div className="p-3 rounded-xl bg-slate-50 border border-slate-100">
              <dt className="text-xs font-semibold text-slate-400 uppercase tracking-wider">Didaftarkan Pada</dt>
              <dd className="text-sm font-medium text-slate-700 mt-1">
                {formatDate(data.created_at)}
              </dd>
            </div>
          </dl>
        </div>

        {/* Hierarki Gedung & Lantai */}
        <div className="bg-white rounded-2xl border border-slate-200 shadow-sm p-6 space-y-4">
          <h3 className="font-bold text-slate-900 text-base">Lokasi & Hierarki Properti</h3>
          <div className="space-y-3">
            <div className="flex items-center justify-between p-3 rounded-xl border border-slate-200 hover:bg-slate-50 transition-colors">
              <div className="flex items-center gap-2.5">
                <Building2 size={18} className="text-orange" />
                <div>
                  <p className="text-[11px] text-slate-400 uppercase font-semibold">Properti</p>
                  <p className="font-bold text-sm text-slate-900">
                    {parentProperty ? parentProperty.name : "—"}
                  </p>
                </div>
              </div>
              {parentProperty && (
                <Button variant="outline" size="sm" onClick={() => navigate(`/dashboard/properties/${parentProperty.id}`)}>
                  Lihat
                </Button>
              )}
            </div>

            <div className="flex items-center justify-between p-3 rounded-xl border border-slate-200 hover:bg-slate-50 transition-colors">
              <div className="flex items-center gap-2.5">
                <Building2 size={18} className="text-slate-500" />
                <div>
                  <p className="text-[11px] text-slate-400 uppercase font-semibold">Gedung / Tower</p>
                  <p className="font-bold text-sm text-slate-900">
                    {parentBuilding ? parentBuilding.name : "—"}
                  </p>
                </div>
              </div>
              {parentBuilding && (
                <Button variant="outline" size="sm" onClick={() => navigate(`/dashboard/buildings/${parentBuilding.id}`)}>
                  Lihat
                </Button>
              )}
            </div>

            <div className="flex items-center justify-between p-3 rounded-xl border border-slate-200 hover:bg-slate-50 transition-colors">
              <div className="flex items-center gap-2.5">
                <Layers size={18} className="text-slate-500" />
                <div>
                  <p className="text-[11px] text-slate-400 uppercase font-semibold">Lantai</p>
                  <p className="font-bold text-sm text-slate-900">
                    {parentFloor ? `${parentFloor.name} (Lantai ${parentFloor.floor_number})` : "—"}
                  </p>
                </div>
              </div>
              {parentFloor && (
                <Button variant="outline" size="sm" onClick={() => navigate(`/dashboard/floors/${parentFloor.id}`)}>
                  Lihat
                </Button>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}


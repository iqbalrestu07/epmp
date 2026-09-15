import { useState, useEffect } from 'react';
import { Outlet, NavLink, useLocation, useNavigate } from "react-router-dom";
import {
  LayoutDashboard, Building2, DoorOpen, Users, ShieldCheck,
  UserCog, Settings, Menu, X, LogOut, ChevronDown, ChevronRight,
  CalendarCheck, FileText, Bed, Receipt, CreditCard, Wrench, Package,
  Globe2, MessageCircle, Megaphone, Layers, Check, Box, BarChart3
} from 'lucide-react';
import { useAuth } from '../features/iam/context/AuthContext';
import { useOrg } from '../features/organization/context/OrgContext';
import { NotificationCenter } from '../features/notification/components/NotificationCenter';

// ─── Menu Groups Configuration ──────────────────────────────────────────────
interface MenuItem {
  label: string;
  path: string;
  icon: React.ElementType;
  requiredPermission?: string;
}

interface MenuGroup {
  id: string;
  label: string;
  items: MenuItem[];
}

const MENU_GROUPS: MenuGroup[] = [
  {
    id: 'main',
    label: 'Utama',
    items: [
      { label: 'Overview', path: '/dashboard', icon: LayoutDashboard },
      { label: '3D Explorer', path: '/explorer', icon: Box },
      { label: 'Direktori Hunian (Roster)', path: '/roster', icon: Bed },
    ],
  },
  {
    id: 'property',
    label: 'Properti & Ruang',
    items: [
      { label: 'Properti', path: '/properties', icon: Building2 },
      { label: 'Gedung', path: '/buildings', icon: Layers },
      { label: 'Lantai', path: '/floors', icon: Layers },
      { label: 'Zona', path: '/zones', icon: Layers },
      { label: 'Kamar & Unit', path: '/rooms', icon: DoorOpen },
      { label: 'Tipe Kamar', path: '/room-types', icon: DoorOpen },
      { label: 'Tempat Tidur', path: '/beds', icon: Bed },
      { label: 'Fasilitas', path: '/facilities', icon: Building2 },
    ],
  },
  {
    id: 'tenancy',
    label: 'Penyewa & Sewa',
    items: [
      { label: 'Daftar Penyewa', path: '/tenants', icon: Users },
      { label: 'Reservasi', path: '/reservations', icon: CalendarCheck },
      { label: 'Kontrak Sewa', path: '/contracts', icon: FileText },
      { label: 'Data Hunian', path: '/occupancies', icon: Bed },
    ],
  },
  {
    id: 'finance',
    label: 'Billing & Keuangan',
    items: [
      { label: 'Tagihan (Invoices)', path: '/invoices', icon: Receipt },
      { label: 'Pembayaran', path: '/payments', icon: CreditCard },
      { label: 'Deposit / Jaminan', path: '/deposits', icon: CreditCard },
      { label: 'Biaya Tambahan', path: '/charges', icon: Receipt },
      { label: 'Pengembalian (Refund)', path: '/refunds', icon: CreditCard },
      { label: 'Penyesuaian (Adjustment)', path: '/adjustments', icon: Receipt },
      { label: 'Denda Keterlambatan', path: '/penalties', icon: Receipt },
    ],
  },
  {
    id: 'maintenance',
    label: 'Aset & Pemeliharaan',
    items: [
      { label: 'Tiket Perbaikan', path: '/work-orders', icon: Wrench },
      { label: 'Master Aset', path: '/assets', icon: Package },
      { label: 'Penugasan Aset', path: '/asset-assignments', icon: Package },
      { label: 'Inspeksi Aset', path: '/asset-inspections', icon: Package },
      { label: 'Teknisi', path: '/technicians', icon: Wrench },
      { label: 'Pemasok / Vendor', path: '/suppliers', icon: Package },
    ],
  },
  {
    id: 'communication',
    label: 'Komunikasi',
    items: [
      { label: 'Blast Message (WA)', path: '/messaging/blast', icon: Megaphone },
      { label: 'Pengaturan Gateway', path: '/messaging/devices', icon: MessageCircle },
    ],
  },
  {
    id: 'reports',
    label: 'Laporan',
    items: [
      { label: 'Laporan & Analitik', path: '/reports', icon: BarChart3 },
    ],
  },
  {
    id: 'system',
    label: 'Sistem & Pengaturan',
    items: [
      { label: 'Organisasi', path: '/organizations', icon: Globe2 },
      { label: 'Role & Akses (RBAC)', path: '/management/rbac', icon: ShieldCheck, requiredPermission: 'role:read' },
      { label: 'Akun Pengguna', path: '/management/users', icon: UserCog, requiredPermission: 'user:read' },
      { label: 'Audit Log', path: '/audit-logs', icon: FileText },
      { label: 'Pengaturan Sistem', path: '/settings', icon: Settings },
    ],
  },
];

export default function MainLayout() {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [orgMenuOpen, setOrgMenuOpen] = useState(false);
  const location = useLocation();
  const navigate = useNavigate();
  const { user, logout, hasPermission } = useAuth();
  const { currentOrg, orgs, switchOrg } = useOrg();

  const isItemActive = (itemPath: string, currentPath: string) => {
    if (itemPath === '/dashboard') {
      return currentPath === '/dashboard' || currentPath === '/overview' || currentPath === '/';
    }
    const cleanItem = itemPath.replace(/^\/dashboard/, '');
    const cleanCurrent = currentPath.replace(/^\/dashboard/, '');
    return cleanCurrent === cleanItem || cleanCurrent.startsWith(cleanItem + '/');
  };

  const [openGroups, setOpenGroups] = useState<Record<string, boolean>>(() => {
    try {
      const saved = localStorage.getItem('epmp_sidebar_groups');
      if (saved) {
        return JSON.parse(saved);
      }
    } catch (e) {}

    // Default: Buka semua grup agar pengguna langsung melihat seluruh menu tanpa harus klik satu-satu
    const initial: Record<string, boolean> = {};
    MENU_GROUPS.forEach(group => {
      initial[group.id] = true;
    });
    return initial;
  });

  useEffect(() => {
    MENU_GROUPS.forEach(group => {
      if (group.items.some(item => isItemActive(item.path, location.pathname))) {
        setOpenGroups(prev => {
          if (prev[group.id]) return prev;
          const next = { ...prev, [group.id]: true };
          try {
            localStorage.setItem('epmp_sidebar_groups', JSON.stringify(next));
          } catch (e) {}
          return next;
        });
      }
    });
  }, [location.pathname]);

  const toggleGroup = (groupId: string) => {
    setOpenGroups(prev => {
      const next = { ...prev, [groupId]: !prev[groupId] };
      try {
        localStorage.setItem('epmp_sidebar_groups', JSON.stringify(next));
      } catch (e) {}
      return next;
    });
  };

  const handleLogout = async () => {
    await logout();
    navigate('/auth/signin', { replace: true });
  };

  const getInitials = (name: string) =>
    name?.split(' ').map(w => w[0]).join('').toUpperCase().slice(0, 2) ?? 'U';

  const primaryRole = user?.roles?.[0]?.name ?? 'User';

  return (
    <div className="min-h-screen bg-[#f2efe9] text-[#0b0b0c] flex">

      {/* Mobile Sidebar Overlay */}
      {sidebarOpen && (
        <div
          className="fixed inset-0 z-40 bg-slate-900/80 backdrop-blur-sm lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      {/* Sidebar - Sticky on desktop so it never hangs on long dashboard pages */}
      <aside className={`fixed top-0 left-0 z-50 h-screen w-72 bg-[#0b0b0c] text-white flex flex-col overflow-hidden transition-transform duration-300 ease-in-out ${sidebarOpen ? 'translate-x-0' : '-translate-x-full'} lg:translate-x-0 lg:sticky lg:top-0 lg:h-screen lg:shrink-0`}>
        {/* Logo */}
        <div className="h-16 flex items-center justify-between px-6 border-b border-white/10 shrink-0">
          <div className="flex items-center gap-3 text-xs tracking-[0.2em] uppercase text-orange font-bold">
            <span className="w-4 h-px bg-orange" />
            EPMP SaaS
          </div>
          <button className="lg:hidden text-white/70 hover:text-white" onClick={() => setSidebarOpen(false)}>
            <X size={20} />
          </button>
        </div>

        {/* Navigation Accordion - min-h-0 and shrink-0 on items prevents any flex squishing or overlapping */}
        <nav className="flex-1 min-h-0 overflow-y-auto py-4 px-3 flex flex-col gap-2 custom-scrollbar">
          {MENU_GROUPS.map((group) => {
            const visibleItems = group.items.filter(item => {
              if (!item.requiredPermission) return true;
              return hasPermission(item.requiredPermission);
            });

            if (visibleItems.length === 0) return null;

            const isOpen = !!openGroups[group.id];
            const hasActiveChild = visibleItems.some(item => isItemActive(item.path, location.pathname));

            return (
              <div key={group.id} className="rounded-xl shrink-0 bg-white/[0.02] border border-white/5 transition-colors">
                <button
                  type="button"
                  onClick={() => toggleGroup(group.id)}
                  className={`w-full flex items-center shrink-0 justify-between px-3.5 py-2.5 text-xs font-semibold tracking-wide uppercase transition-colors select-none ${
                    hasActiveChild ? 'text-orange' : 'text-white/60 hover:text-white hover:bg-white/5'
                  }`}
                >
                  <span className="flex items-center gap-2">
                    {group.label}
                  </span>
                  <div className="text-white/40">
                    {isOpen ? <ChevronDown size={14} /> : <ChevronRight size={14} />}
                  </div>
                </button>

                {isOpen && (
                  <div className="flex flex-col gap-1 px-2 pb-2 pt-0.5">
                    {visibleItems.map((item, idx) => {
                      const Icon = item.icon;
                      const active = isItemActive(item.path, location.pathname);

                      return (
                        <NavLink
                          key={idx}
                          to={item.path}
                          onClick={() => setSidebarOpen(false)}
                          className={`flex items-center shrink-0 gap-2.5 px-3 py-2 rounded-lg text-sm transition-all duration-150 ${
                            active
                              ? 'bg-orange text-slate-900 font-semibold shadow-[0_2px_8px_rgba(255,102,0,0.3)]'
                              : 'text-white/70 hover:bg-white/10 hover:text-white'
                          }`}
                        >
                          <Icon size={16} className={active ? 'text-slate-900' : 'text-orange'} />
                          <span className="truncate">{item.label}</span>
                        </NavLink>
                      );
                    })}
                  </div>
                )}
              </div>
            );
          })}
        </nav>

        {/* User Footer */}
        <div className="p-4 border-t border-white/10 shrink-0 bg-[#0b0b0c]">
          <div className="flex items-center gap-3 bg-white/5 p-3 rounded-xl">
            <div className="w-10 h-10 rounded-full bg-orange flex items-center justify-center text-slate-900 font-bold text-sm flex-shrink-0">
              {user ? getInitials(user.name) : 'U'}
            </div>
            <div className="flex-1 overflow-hidden">
              <p className="text-sm font-semibold truncate">{user?.name ?? 'Loading…'}</p>
              <p className="text-xs text-white/50 truncate capitalize">{primaryRole.replace('_', ' ')}</p>
            </div>
            <button
              onClick={handleLogout}
              className="text-white/50 hover:text-red-400 transition-colors p-1"
              title="Sign out"
            >
              <LogOut size={16} />
            </button>
          </div>
        </div>
      </aside>

      {/* Main Content */}
      <div className="flex-1 flex flex-col min-w-0">

        {/* Top Header */}
        <header className="h-16 bg-white border-b border-slate-200 flex items-center justify-between px-6 sticky top-0 z-30 shadow-sm">
          <div className="flex items-center gap-4">
            <button className="lg:hidden text-slate-600 hover:text-slate-900" onClick={() => setSidebarOpen(true)}>
              <Menu size={24} />
            </button>
            <h2 className="text-lg font-semibold capitalize hidden sm:block">
              {location.pathname.split('/').pop()?.replace('-', ' ') || 'Overview'}
            </h2>
          </div>

          <div className="flex items-center gap-4">
            <div className="relative">
              <button
                onClick={() => setOrgMenuOpen(!orgMenuOpen)}
                className="hidden md:flex items-center gap-2 bg-slate-100 hover:bg-slate-200 px-4 py-2 rounded-lg text-sm font-medium transition-colors"
              >
                <Globe2 size={16} className="text-orange" />
                <span className="max-w-40 truncate">{currentOrg?.name ?? 'No Organization'}</span>
                <ChevronDown size={16} className="text-slate-500 ml-1" />
              </button>

              {orgMenuOpen && (
                <>
                  <div
                    className="fixed inset-0 z-40"
                    onClick={() => setOrgMenuOpen(false)}
                  />
                  <div className="absolute right-0 top-full mt-2 z-50 w-64 bg-white rounded-xl shadow-lg border border-slate-200 overflow-hidden">
                    <div className="px-4 py-3 border-b border-slate-200 text-xs font-bold text-slate-500 tracking-wider">
                      YOUR ORGANIZATIONS
                    </div>
                    <div className="max-h-64 overflow-y-auto">
                      {orgs.length === 0 && (
                        <div className="px-4 py-6 text-center text-sm text-slate-500">
                          No organizations yet.
                          <button
                            onClick={() => { navigate('/dashboard/organizations/new'); setOrgMenuOpen(false); }}
                            className="block w-full mt-2 text-orange font-semibold hover:underline"
                          >
                            Create one →
                          </button>
                        </div>
                      )}
                      {orgs.map(org => (
                        <button
                          key={org.id}
                          onClick={() => { switchOrg(org.id); setOrgMenuOpen(false); }}
                          className="w-full flex items-center justify-between px-4 py-3 hover:bg-slate-100 transition-colors text-left"
                        >
                          <div className="flex-1 min-w-0">
                            <p className="text-sm font-medium truncate">{org.name}</p>
                            <p className="text-xs text-slate-500 truncate">{org.domain}</p>
                          </div>
                          {currentOrg?.id === org.id && (
                            <Check size={16} className="text-orange flex-shrink-0 ml-2" />
                          )}
                        </button>
                      ))}
                    </div>
                    <div className="border-t border-slate-200">
                      <button
                        onClick={() => { navigate('/dashboard/organizations/new'); setOrgMenuOpen(false); }}
                        className="w-full flex items-center gap-2 px-4 py-3 text-sm font-medium text-orange hover:bg-orange/5 transition-colors"
                      >
                        <Globe2 size={16} />
                        Create Organization
                      </button>
                    </div>
                  </div>
                </>
              )}
            </div>

            <NotificationCenter />
          </div>
        </header>

        {/* Page Content */}
        <main className="flex-1 overflow-y-auto p-6 md:p-10 bg-[#f2efe9] dashboard-light">
          <div className="max-w-7xl mx-auto">
            <Outlet />
          </div>
        </main>
      </div>
    </div>
  );
}

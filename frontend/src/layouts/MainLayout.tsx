import { useState } from 'react';
import { Outlet, NavLink, useLocation, useNavigate } from "react-router-dom";
import {
  LayoutDashboard, Building2, DoorOpen, Users, ShieldCheck,
  UserCog, Settings, Menu, X, Bell, LogOut, ChevronDown,
  CalendarCheck, FileText, Bed, Receipt, CreditCard, Wrench, Package,
  Globe2, MessageCircle, Megaphone, Layers, Check
} from 'lucide-react';
import { useAuth } from '../features/iam/context/AuthContext';
import { useOrg } from '../features/organization/context/OrgContext';

// ─── Menu Configuration ────────────────────────────────────────────────────
const MENU_CONFIG = [
  { label: 'Overview',       path: '/dashboard',                  icon: LayoutDashboard },
  { type: 'divider' as const, label: 'CORE' },
  { label: 'Organizations',  path: '/dashboard/organizations',    icon: Globe2,     requiredPermission: 'property:read' },
  { label: 'Properties',     path: '/dashboard/properties',       icon: Building2,  requiredPermission: 'property:read' },
  { label: 'Buildings',      path: '/dashboard/buildings',        icon: Layers,     requiredPermission: 'property:read' },
  { label: 'Floors',         path: '/dashboard/floors',           icon: Layers,     requiredPermission: 'property:read' },
  { label: 'Rooms & Units',  path: '/dashboard/rooms',            icon: DoorOpen,   requiredPermission: 'room:read'     },
  { label: 'Tenants',        path: '/dashboard/tenants',          icon: Users,      requiredPermission: 'tenant:read'   },
  
  { type: 'divider' as const, label: 'OPERATIONS' },
  { label: 'Reservations',   path: '/dashboard/reservations',     icon: CalendarCheck },
  { label: 'Contracts',      path: '/dashboard/contracts',        icon: FileText },
  { label: 'Occupancies',    path: '/dashboard/occupancies',      icon: Bed },
  
  { type: 'divider' as const, label: 'FINANCE' },
  { label: 'Invoices',       path: '/dashboard/invoices',         icon: Receipt },
  { label: 'Payments',       path: '/dashboard/payments',         icon: CreditCard },
  
  { type: 'divider' as const, label: 'MAINTENANCE & ASSETS' },
  { label: 'Work Orders',    path: '/dashboard/work-orders',      icon: Wrench },
  { label: 'Assets',         path: '/dashboard/assets',           icon: Package },

  { type: 'divider' as const, label: 'COMMUNICATION' },
  { label: 'Blast Message',  path: '/dashboard/messaging/blast',  icon: Megaphone },
  { label: 'Msg Settings',   path: '/dashboard/messaging/devices',icon: MessageCircle },

  { type: 'divider' as const, label: 'SYSTEM' },
  { label: 'Roles & Perms',  path: '/dashboard/management/rbac',  icon: ShieldCheck, requiredPermission: 'role:read'   },
  { label: 'User Accounts',  path: '/dashboard/management/users', icon: UserCog,    requiredPermission: 'user:read'    },
  { label: 'Settings',       path: '/dashboard/settings',         icon: Settings    },
];

export default function MainLayout() {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [orgMenuOpen, setOrgMenuOpen] = useState(false);
  const location = useLocation();
  const navigate = useNavigate();
  const { user, logout, hasPermission } = useAuth();
  const { currentOrg, orgs, switchOrg } = useOrg();

  const handleLogout = async () => {
    await logout();
    navigate('/auth/signin', { replace: true });
  };

  const visibleMenu = MENU_CONFIG.filter(item => {
    if (item.type === 'divider') return true;
    if (!item.requiredPermission) return true;
    return hasPermission(item.requiredPermission);
  });

  const getInitials = (name: string) =>
    name?.split(' ').map(w => w[0]).join('').toUpperCase().slice(0, 2) ?? 'U';

  const primaryRole = user?.roles?.[0]?.name ?? 'User';

  return (
    <div className="min-h-screen bg-[#f2efe9] text-[#0b0b0c] flex">

      {/* Mobile Sidebar Overlay */}
      {sidebarOpen && (
        <div
          className="fixed inset-0 z-40 bg-black/50 lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      {/* Sidebar */}
      <aside className={`fixed top-0 left-0 z-50 h-screen w-72 bg-[#0b0b0c] text-white flex flex-col transition-transform duration-300 ease-in-out ${sidebarOpen ? 'translate-x-0' : '-translate-x-full'} lg:translate-x-0 lg:static`}>
        {/* Logo */}
        <div className="h-16 flex items-center justify-between px-6 border-b border-white/10">
          <div className="flex items-center gap-3 text-xs tracking-[0.2em] uppercase text-orange font-bold">
            <span className="w-4 h-px bg-orange" />
            EPMP SaaS
          </div>
          <button className="lg:hidden text-white/70 hover:text-white" onClick={() => setSidebarOpen(false)}>
            <X size={20} />
          </button>
        </div>

        {/* Navigation */}
        <nav className="flex-1 overflow-y-auto py-6 px-4 flex flex-col gap-1 custom-scrollbar">
          {visibleMenu.map((item, idx) => {
            if (item.type === 'divider') {
              return (
                <div key={idx} className="mt-6 mb-2 px-4 text-[10px] font-bold text-white/40 tracking-wider">
                  {item.label}
                </div>
              );
            }

            const Icon = item.icon as React.ElementType;
            const isActive =
              location.pathname === item.path ||
              (item.path !== '/dashboard' && location.pathname.startsWith(item.path!));

            return (
              <NavLink
                key={idx}
                to={item.path || '#'}
                onClick={() => setSidebarOpen(false)}
                className={`flex items-center gap-3 px-4 py-3 rounded-xl transition-all duration-200 ${
                  isActive
                    ? 'bg-orange text-black font-semibold shadow-[0_4px_12px_rgba(255,102,0,0.3)]'
                    : 'text-white/70 hover:bg-white/10 hover:text-white'
                }`}
              >
                <Icon size={18} className={isActive ? 'text-black' : 'text-orange'} />
                {item.label}
              </NavLink>
            );
          })}
        </nav>

        {/* User Footer */}
        <div className="p-4 border-t border-white/10">
          <div className="flex items-center gap-3 bg-white/5 p-3 rounded-xl">
            <div className="w-10 h-10 rounded-full bg-orange flex items-center justify-center text-black font-bold text-sm flex-shrink-0">
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
        <header className="h-16 bg-white border-b border-black/5 flex items-center justify-between px-6 sticky top-0 z-30 shadow-sm">
          <div className="flex items-center gap-4">
            <button className="lg:hidden text-black/70 hover:text-black" onClick={() => setSidebarOpen(true)}>
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
                className="hidden md:flex items-center gap-2 bg-black/5 hover:bg-black/10 px-4 py-2 rounded-lg text-sm font-medium transition-colors"
              >
                <Globe2 size={16} className="text-orange" />
                <span className="max-w-40 truncate">{currentOrg?.name ?? 'No Organization'}</span>
                <ChevronDown size={16} className="text-black/50 ml-1" />
              </button>

              {orgMenuOpen && (
                <>
                  <div
                    className="fixed inset-0 z-40"
                    onClick={() => setOrgMenuOpen(false)}
                  />
                  <div className="absolute right-0 top-full mt-2 z-50 w-64 bg-white rounded-xl shadow-lg border border-black/5 overflow-hidden">
                    <div className="px-4 py-3 border-b border-black/5 text-xs font-bold text-black/40 tracking-wider">
                      YOUR ORGANIZATIONS
                    </div>
                    <div className="max-h-64 overflow-y-auto">
                      {orgs.length === 0 && (
                        <div className="px-4 py-6 text-center text-sm text-black/40">
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
                          className="w-full flex items-center justify-between px-4 py-3 hover:bg-black/5 transition-colors text-left"
                        >
                          <div className="flex-1 min-w-0">
                            <p className="text-sm font-medium truncate">{org.name}</p>
                            <p className="text-xs text-black/40 truncate">{org.domain}</p>
                          </div>
                          {currentOrg?.id === org.id && (
                            <Check size={16} className="text-orange flex-shrink-0 ml-2" />
                          )}
                        </button>
                      ))}
                    </div>
                    <div className="border-t border-black/5">
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

            <button className="w-10 h-10 rounded-full bg-black/5 hover:bg-black/10 flex items-center justify-center relative transition-colors">
              <Bell size={18} className="text-black/70" />
              <span className="absolute top-2 right-2 w-2 h-2 bg-orange rounded-full border border-white" />
            </button>
          </div>
        </header>

        {/* Page Content */}
        <main className="flex-1 overflow-y-auto p-6 md:p-10 bg-[#f2efe9]">
          <div className="max-w-7xl mx-auto">
            <Outlet />
          </div>
        </main>
      </div>
    </div>
  );
}

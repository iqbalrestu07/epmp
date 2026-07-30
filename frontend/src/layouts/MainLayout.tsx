import { useState } from 'react';
import { Outlet, NavLink, useLocation } from "react-router-dom";
import { LayoutDashboard, Building2, DoorOpen, Users, ShieldCheck, UserCog, Settings, Menu, X, Bell, LogOut, ChevronDown } from 'lucide-react';

// DYNAMIC MENU CONFIGURATION
// Nanti bisa difilter berdasarkan rbac user: MENU_CONFIG.filter(item => user.permissions.includes(item.requiredPermission))
// Atau berdasarkan package langganan: MENU_CONFIG.filter(item => user.org.package.includes(item.requiredPackage))
const MENU_CONFIG = [
  { label: 'Overview', path: '/dashboard', icon: LayoutDashboard, requiredPackage: 'core' },
  { label: 'Properties', path: '/dashboard/properties', icon: Building2, requiredPackage: 'property' },
  { label: 'Rooms & Units', path: '/dashboard/rooms', icon: DoorOpen, requiredPackage: 'property' },
  { label: 'Tenants', path: '/dashboard/tenants', icon: Users, requiredPackage: 'tenant' },
  { type: 'divider' }, // Visual separator
  { 
    label: 'Management (RBAC)', 
    path: '/dashboard/management/rbac', 
    icon: ShieldCheck, 
    requiredRole: 'admin' 
  },
  { 
    label: 'User Accounts', 
    path: '/dashboard/management/users', 
    icon: UserCog, 
    requiredRole: 'admin' 
  },
  { 
    label: 'Settings', 
    path: '/dashboard/settings', 
    icon: Settings, 
    requiredPackage: 'core' 
  },
];

export default function MainLayout() {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const location = useLocation();

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
        {/* Sidebar Header */}
        <div className="h-16 flex items-center justify-between px-6 border-b border-white/10">
          <div className="flex items-center gap-3 text-xs tracking-[0.2em] uppercase text-orange font-bold">
            <span className="w-4 h-px bg-orange"></span>
            EPMP
          </div>
          <button className="lg:hidden text-white/70 hover:text-white" onClick={() => setSidebarOpen(false)}>
            <X size={20} />
          </button>
        </div>

        {/* Sidebar Navigation */}
        <nav className="flex-1 overflow-y-auto py-6 px-4 flex flex-col gap-1">
          {MENU_CONFIG.map((item, idx) => {
            if (item.type === 'divider') {
              return <hr key={idx} className="my-4 border-white/10" />;
            }
            
            const Icon = item.icon as React.ElementType;
            const isActive = location.pathname === item.path || (item.path !== '/dashboard' && location.pathname.startsWith(item.path!));
            
            return (
              <NavLink 
                key={idx} 
                to={item.path || '#'}
                onClick={() => setSidebarOpen(false)}
                className={`flex items-center gap-3 px-4 py-3 rounded-xl transition-all duration-200 ${
                  isActive 
                  ? 'bg-orange text-black font-semibold' 
                  : 'text-white/70 hover:bg-white/10 hover:text-white'
                }`}
              >
                <Icon size={18} className={isActive ? 'text-black' : 'text-orange'} />
                {item.label}
              </NavLink>
            );
          })}
        </nav>

        {/* Sidebar Footer (User Info) */}
        <div className="p-4 border-t border-white/10">
          <div className="flex items-center gap-3 bg-white/5 p-3 rounded-xl hover:bg-white/10 transition-colors cursor-pointer">
            <div className="w-10 h-10 rounded-full bg-orange flex items-center justify-center text-black font-bold">
              AD
            </div>
            <div className="flex-1 overflow-hidden">
              <p className="text-sm font-semibold truncate">Admin Acme</p>
              <p className="text-xs text-white/50 truncate">Super Administrator</p>
            </div>
            <LogOut size={16} className="text-white/50" />
          </div>
        </div>
      </aside>

      {/* Main Content Area */}
      <div className="flex-1 flex flex-col min-w-0">
        
        {/* Top Header */}
        <header className="h-16 bg-white border-b border-black/5 flex items-center justify-between px-6 sticky top-0 z-30 shadow-sm">
          <div className="flex items-center gap-4">
            <button className="lg:hidden text-black/70 hover:text-black" onClick={() => setSidebarOpen(true)}>
              <Menu size={24} />
            </button>
            <h2 className="text-lg font-semibold capitalize hidden sm:block">
              {location.pathname.split('/').pop() || 'Overview'}
            </h2>
          </div>

          <div className="flex items-center gap-4">
            {/* Context Switcher (e.g., Organization / Property selector) */}
            <button className="hidden md:flex items-center gap-2 bg-black/5 hover:bg-black/10 px-4 py-2 rounded-lg text-sm font-medium transition-colors">
              <Building2 size={16} className="text-orange" />
              Acme Properties Ltd.
              <ChevronDown size={16} className="text-black/50 ml-2" />
            </button>
            
            {/* Notifications */}
            <button className="w-10 h-10 rounded-full bg-black/5 hover:bg-black/10 flex items-center justify-center relative transition-colors">
              <Bell size={18} className="text-black/70" />
              <span className="absolute top-2 right-2 w-2 h-2 bg-orange rounded-full border border-white"></span>
            </button>
          </div>
        </header>

        {/* Page Content */}
        <main className="flex-1 overflow-y-auto p-6 md:p-10">
          <div className="max-w-7xl mx-auto">
            <Outlet />
          </div>
        </main>
      </div>

    </div>
  );
}

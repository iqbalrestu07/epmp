import { useState, useEffect, useCallback } from 'react';
import {
  Users,
  Building2,
  ShieldCheck,
  Bell,
  Plus,
  Trash2,
  Check,
  Loader2,
  AlertTriangle,
  UserPlus,
  Save,
  Shield,
  Pencil,
  X,
  Zap
} from 'lucide-react';
import { useOrg } from '@/features/organization/context/OrgContext';
import { usePropertys } from '@/features/property/hooks';
import { useInvoices } from '@/features/billing/hooks';
import { createPenalty } from '@/features/penalty/api';
import { roleService, permissionService } from '@/features/iam/services/iamService';
import type { Role, Permission, CreateRoleRequest } from '@/features/iam/types';
import { CurrencySelect } from '@/components/ui/CurrencySelect';
import { settingsApi, type OrgMember, type PropertyStaff } from '../api/settingsApi';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';

export function SettingsPage() {
  const { currentOrg } = useOrg();
  const [activeTab, setActiveTab] = useState<'members' | 'property_rbac' | 'custom_roles' | 'policies'>('members');

  // ── Tab 1: Organization Members State ──
  const [members, setMembers] = useState<OrgMember[]>([]);
  const [membersLoading, setMembersLoading] = useState(false);
  const [isAddMemberOpen, setIsAddMemberOpen] = useState(false);
  const [newMemberEmail, setNewMemberEmail] = useState('');
  const [newMemberName, setNewMemberName] = useState('');
  const [newMemberPassword, setNewMemberPassword] = useState('');
  const [newMemberRole, setNewMemberRole] = useState('member');
  const [addMemberLoading, setAddMemberLoading] = useState(false);
  const [memberError, setMemberError] = useState<string | null>(null);

  // ── Tab 2: Property Staff State ──
  const { data: propertiesData, isLoading: propertiesLoading } = usePropertys({ per_page: 100 });
  const properties = Array.isArray(propertiesData?.data) ? propertiesData.data : [];
  const [selectedPropertyId, setSelectedPropertyId] = useState<string>('');
  const [propertyStaff, setPropertyStaff] = useState<PropertyStaff[]>([]);
  const [staffLoading, setStaffLoading] = useState(false);
  const [isAssignStaffOpen, setIsAssignStaffOpen] = useState(false);
  const [assignUserId, setAssignUserId] = useState('');
  const [assignRoleId, setAssignRoleId] = useState('');
  const [assignStaffLoading, setAssignStaffLoading] = useState(false);
  const [assignError, setAssignError] = useState<string | null>(null);

  // ── Tab 3: Roles & Permissions State ──
  const [roles, setRoles] = useState<Role[]>([]);
  const [permissions, setPermissions] = useState<Permission[]>([]);
  const [rolesLoading, setRolesLoading] = useState(false);
  const [roleModalOpen, setRoleModalOpen] = useState(false);
  const [editingRole, setEditingRole] = useState<Role | null>(null);
  const [newRoleName, setNewRoleName] = useState('');
  const [newRoleDescription, setNewRoleDescription] = useState('');
  const [selectedPermIds, setSelectedPermIds] = useState<Set<string>>(new Set());
  const [savingRole, setSavingRole] = useState(false);
  const [roleError, setRoleError] = useState<string | null>(null);

  // ── Tab 4: Policies State ──
  const [defaultCurrency, setDefaultCurrency] = useState(localStorage.getItem('epmp_default_currency') || 'IDR');
  const [dueSoonDays, setDueSoonDays] = useState(localStorage.getItem('epmp_policy_due_soon') || '7');
  const [expiringSoonDays, setExpiringSoonDays] = useState(localStorage.getItem('epmp_policy_expiring_soon') || '30');
  const [waTemplate, setWaTemplate] = useState(
    localStorage.getItem('epmp_policy_wa_template') ||
    'Halo {tenant_name}, tagihan kos Anda sebesar {amount} akan jatuh tempo pada {due_date}. Mohon lakukan pembayaran melalui portal EPMP. Terima kasih!'
  );
  const [latePenaltyType, setLatePenaltyType] = useState(localStorage.getItem('epmp_late_penalty_type') || 'daily_fixed');
  const [latePenaltyAmount, setLatePenaltyAmount] = useState(localStorage.getItem('epmp_late_penalty_amount') || '25000');
  const [lateGracePeriod, setLateGracePeriod] = useState(localStorage.getItem('epmp_late_grace_period') || '3');
  const [autoPenaltyProcessing, setAutoPenaltyProcessing] = useState(false);
  const [autoPenaltyResult, setAutoPenaltyResult] = useState<{ success: number; skipped: number } | null>(null);
  const [policySaved, setPolicySaved] = useState(false);

  const { data: allInvoicesData } = useInvoices({ per_page: 200 });

  // ── Load Organization Members ──
  const loadMembers = useCallback(async () => {
    if (!currentOrg?.id) return;
    setMembersLoading(true);
    try {
      const list = await settingsApi.getOrgMembers();
      setMembers(list);
    } catch (err: any) {
      console.error('Failed to load org members', err);
    } finally {
      setMembersLoading(false);
    }
  }, [currentOrg?.id]);

  // ── Load Roles & Permissions ──
  const loadRolesAndPerms = useCallback(async () => {
    setRolesLoading(true);
    try {
      const [rRes, pRes] = await Promise.all([
        roleService.list(),
        permissionService.list(),
      ]);
      setRoles(rRes.data ?? []);
      setPermissions(pRes.data ?? []);
    } catch (err: any) {
      console.error('Failed to load roles and permissions', err);
    } finally {
      setRolesLoading(false);
    }
  }, []);

  // ── Load Property Staff ──
  const loadPropertyStaff = useCallback(async (propertyId: string) => {
    if (!propertyId) return;
    setStaffLoading(true);
    try {
      const staffList = await settingsApi.getPropertyStaff(propertyId);
      setPropertyStaff(staffList);
    } catch (err: any) {
      console.error('Failed to load property staff', err);
      setPropertyStaff([]);
    } finally {
      setStaffLoading(false);
    }
  }, []);

  useEffect(() => {
    if (currentOrg?.id) {
      loadMembers();
      loadRolesAndPerms();
    }
  }, [currentOrg?.id, loadMembers, loadRolesAndPerms]);

  // Set initial selected property
  useEffect(() => {
    if (properties.length > 0 && !selectedPropertyId) {
      setSelectedPropertyId(properties[0].id);
    }
  }, [properties, selectedPropertyId]);

  useEffect(() => {
    if (selectedPropertyId) {
      loadPropertyStaff(selectedPropertyId);
    }
  }, [selectedPropertyId, loadPropertyStaff]);

  // ── Actions: Members ──
  const handleAddMember = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newMemberEmail.trim()) {
      setMemberError('Email wajib diisi');
      return;
    }
    setAddMemberLoading(true);
    setMemberError(null);
    try {
      await settingsApi.addOrgMember({
        email: newMemberEmail.trim(),
        name: newMemberName.trim() || undefined,
        password: newMemberPassword.trim() || 'Password123!',
        role: newMemberRole,
      });
      setIsAddMemberOpen(false);
      setNewMemberEmail('');
      setNewMemberName('');
      setNewMemberPassword('');
      await loadMembers();
    } catch (err: any) {
      setMemberError(err?.message || 'Gagal menambahkan anggota organisasi');
    } finally {
      setAddMemberLoading(false);
    }
  };

  const handleRemoveMember = async (userId: string, userName: string) => {
    if (!confirm(`Hapus anggota "${userName}" dari organisasi ini?`)) return;
    try {
      await settingsApi.removeOrgMember(userId);
      await loadMembers();
    } catch (err: any) {
      alert(err?.message || 'Gagal menghapus anggota');
    }
  };

  // ── Actions: Property Staff ──
  const handleAssignStaff = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedPropertyId) {
      setAssignError('Pilih properti terlebih dahulu');
      return;
    }
    if (!assignUserId || !assignRoleId) {
      setAssignError('Pilih anggota dan peran yang akan ditugaskan');
      return;
    }
    setAssignStaffLoading(true);
    setAssignError(null);
    try {
      await settingsApi.assignPropertyStaff(selectedPropertyId, {
        user_id: assignUserId,
        role_id: assignRoleId,
      });
      setIsAssignStaffOpen(false);
      setAssignUserId('');
      setAssignRoleId('');
      await loadPropertyStaff(selectedPropertyId);
    } catch (err: any) {
      setAssignError(err?.message || 'Gagal menugaskan peran staf ke properti');
    } finally {
      setAssignStaffLoading(false);
    }
  };

  const handleRemoveStaff = async (userId: string, roleId: string, roleName: string) => {
    if (!confirm(`Cabut peran "${roleName}" dari pengguna ini pada properti?`)) return;
    try {
      await settingsApi.removePropertyStaff(selectedPropertyId, userId, roleId);
      await loadPropertyStaff(selectedPropertyId);
    } catch (err: any) {
      alert(err?.message || 'Gagal mencabut peran staf');
    }
  };

  // ── Actions: Roles ──
  const handleSaveRole = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newRoleName.trim()) {
      setRoleError('Nama peran wajib diisi');
      return;
    }
    setSavingRole(true);
    setRoleError(null);
    try {
      if (editingRole) {
        await roleService.update(editingRole.id, {
          name: newRoleName.trim(),
          description: newRoleDescription.trim(),
        });
        if (!editingRole.is_system) {
          await roleService.setPermissions(editingRole.id, {
            permission_ids: Array.from(selectedPermIds),
          });
        }
      } else {
        const req: CreateRoleRequest = {
          name: newRoleName.trim(),
          description: newRoleDescription.trim(),
          permission_ids: Array.from(selectedPermIds),
        };
        await roleService.create(req);
      }
      setRoleModalOpen(false);
      setEditingRole(null);
      setNewRoleName('');
      setNewRoleDescription('');
      setSelectedPermIds(new Set());
      await loadRolesAndPerms();
    } catch (err: any) {
      setRoleError(err?.message || 'Gagal menyimpan peran');
    } finally {
      setSavingRole(false);
    }
  };

  const handleDeleteRole = async (role: Role) => {
    if (role.is_system) {
      alert('Peran sistem bawaan tidak dapat dihapus');
      return;
    }
    if (!confirm(`Hapus peran "${role.name}"? Peran ini tidak akan bisa digunakan lagi.`)) return;
    try {
      await roleService.delete(role.id);
      await loadRolesAndPerms();
    } catch (err: any) {
      alert(err?.message || 'Gagal menghapus peran');
    }
  };

  // ── Actions: Policies ──
  const handleSavePolicies = (e: React.FormEvent) => {
    e.preventDefault();
    localStorage.setItem('epmp_default_currency', defaultCurrency);
    localStorage.setItem('epmp_policy_due_soon', dueSoonDays);
    localStorage.setItem('epmp_policy_expiring_soon', expiringSoonDays);
    localStorage.setItem('epmp_policy_wa_template', waTemplate);
    localStorage.setItem('epmp_late_penalty_type', latePenaltyType);
    localStorage.setItem('epmp_late_penalty_amount', latePenaltyAmount);
    localStorage.setItem('epmp_late_grace_period', lateGracePeriod);
    setPolicySaved(true);
    setTimeout(() => setPolicySaved(false), 3000);
  };

  const handleAutoGeneratePenalties = async () => {
    if (!confirm('Apakah Anda yakin ingin memproses dan membuat denda keterlambatan secara otomatis untuk seluruh tagihan yang belum lunas dan telah melewati masa tenggang?')) {
      return;
    }

    setAutoPenaltyProcessing(true);
    setAutoPenaltyResult(null);
    try {
      const invoices = Array.isArray(allInvoicesData?.data) ? allInvoicesData.data : [];
      const now = new Date();
      const graceDays = Number(lateGracePeriod) || 0;
      let createdCount = 0;
      let skippedCount = 0;

      for (const inv of invoices) {
        if (inv.status === 'Paid' || !inv.due_date) {
          skippedCount++;
          continue;
        }

        const due = new Date(inv.due_date);
        const overdueDays = Math.floor((now.getTime() - due.getTime()) / (1000 * 60 * 60 * 24));

        if (overdueDays > graceDays) {
          const effectiveDays = overdueDays - graceDays;
          let penaltyAmount = 0;
          const rate = Number(latePenaltyAmount) || 0;

          if (latePenaltyType === 'daily_fixed') {
            penaltyAmount = rate * effectiveDays;
          } else if (latePenaltyType === 'daily_percent') {
            penaltyAmount = ((inv.amount * rate) / 100) * effectiveDays;
          } else {
            penaltyAmount = rate;
          }

          if (penaltyAmount > 0) {
            await createPenalty({
              invoice_id: inv.id,
              amount: Math.round(penaltyAmount),
              status: 'Applied',
              penalty_date: new Date().toISOString().split('T')[0],
              description: `Denda otomatis: terlambat ${overdueDays} hari (lewat grace period ${graceDays} hari)`,
            });
            createdCount++;
          }
        } else {
          skippedCount++;
        }
      }

      setAutoPenaltyResult({ success: createdCount, skipped: skippedCount });
    } catch (err: any) {
      console.error('Failed to auto generate penalties', err);
      alert('Gagal memproses denda: ' + (err.message || 'Error'));
    } finally {
      setAutoPenaltyProcessing(false);
    }
  };

  // Group permissions by resource for role modal
  const groupedPermissions = permissions.reduce<Record<string, Permission[]>>((acc, p) => {
    if (!acc[p.resource]) acc[p.resource] = [];
    acc[p.resource].push(p);
    return acc;
  }, {});

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold text-slate-900">Pengaturan & RBAC</h1>
        <p className="text-sm text-slate-500 mt-1">
          Kelola anggota organisasi, penugasan staf per properti, peran hak akses kustom, serta kebijakan peringatan.
        </p>
      </div>

      {/* Organization Banner */}
      <div className="bg-gradient-to-r from-orange/10 via-amber-50 to-orange/5 border border-orange/20 rounded-2xl p-5 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-2">
            <span className="text-xs font-semibold px-2.5 py-0.5 rounded-full bg-orange text-slate-900">
              Organisasi Aktif
            </span>
            <span className="text-xs text-slate-400 font-mono">ID: {currentOrg?.id}</span>
          </div>
          <h2 className="text-xl font-bold text-slate-900 mt-1">{currentOrg?.name ?? 'Memuat...'}</h2>
          <p className="text-xs text-slate-500">Domain: {currentOrg?.domain || 'Belum diatur'}</p>
        </div>
        <div className="flex items-center gap-2">
          <div className="text-right hidden sm:block">
            <p className="text-xs text-slate-500">Total Anggota</p>
            <p className="text-lg font-bold text-slate-900">{members.length} Orang</p>
          </div>
        </div>
      </div>

      {/* Tabs */}
      <div className="flex border-b border-slate-200 gap-2 overflow-x-auto">
        <button
          onClick={() => setActiveTab('members')}
          className={`flex items-center gap-2 pb-3 px-3 text-sm font-semibold border-b-2 transition-colors whitespace-nowrap ${
            activeTab === 'members'
              ? 'border-orange text-slate-900'
              : 'border-transparent text-slate-500 hover:text-slate-700'
          }`}
        >
          <Users className="w-4 h-4" />
          Anggota Organisasi
          <span className="ml-1 text-xs px-2 py-0.5 rounded-full bg-slate-100 text-slate-600 font-medium">
            {members.length}
          </span>
        </button>

        <button
          onClick={() => setActiveTab('property_rbac')}
          className={`flex items-center gap-2 pb-3 px-3 text-sm font-semibold border-b-2 transition-colors whitespace-nowrap ${
            activeTab === 'property_rbac'
              ? 'border-orange text-slate-900'
              : 'border-transparent text-slate-500 hover:text-slate-700'
          }`}
        >
          <Building2 className="w-4 h-4" />
          Hak Akses Properti (Multi-Level RBAC)
          <span className="ml-1 text-xs px-2 py-0.5 rounded-full bg-slate-100 text-slate-600 font-medium">
            {properties.length} Properti
          </span>
        </button>

        <button
          onClick={() => setActiveTab('custom_roles')}
          className={`flex items-center gap-2 pb-3 px-3 text-sm font-semibold border-b-2 transition-colors whitespace-nowrap ${
            activeTab === 'custom_roles'
              ? 'border-orange text-slate-900'
              : 'border-transparent text-slate-500 hover:text-slate-700'
          }`}
        >
          <ShieldCheck className="w-4 h-4" />
          Peran Kustom & Izin
          <span className="ml-1 text-xs px-2 py-0.5 rounded-full bg-slate-100 text-slate-600 font-medium">
            {roles.length}
          </span>
        </button>

        <button
          onClick={() => setActiveTab('policies')}
          className={`flex items-center gap-2 pb-3 px-3 text-sm font-semibold border-b-2 transition-colors whitespace-nowrap ${
            activeTab === 'policies'
              ? 'border-orange text-slate-900'
              : 'border-transparent text-slate-500 hover:text-slate-700'
          }`}
        >
          <Bell className="w-4 h-4" />
          Kebijakan & Peringatan
        </button>
      </div>

      {/* ─────────────────────────────────────────────────────────────
          TAB 1: ANGGOTA ORGANISASI (Tier 1)
      ────────────────────────────────────────────────────────────── */}
      {activeTab === 'members' && (
        <div className="space-y-4">
          <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3">
            <div>
              <h3 className="text-lg font-bold text-slate-900">Daftar Anggota Organisasi</h3>
              <p className="text-xs text-slate-500">
                Pengguna yang terdaftar di bawah organisasi ini. Anggota dapat ditugaskan peran di tiap properti.
              </p>
            </div>
            <Button
              onClick={() => {
                setMemberError(null);
                setIsAddMemberOpen(true);
              }}
              className="bg-slate-900 text-white hover:bg-slate-800"
            >
              <UserPlus className="w-4 h-4 mr-2" />
              Tambah Anggota Baru
            </Button>
          </div>

          <div className="rounded-2xl border border-slate-200 bg-white overflow-hidden shadow-sm">
            {membersLoading ? (
              <div className="p-8 text-center text-slate-400">
                <Loader2 className="w-6 h-6 animate-spin mx-auto mb-2 text-orange" />
                Memuat anggota organisasi...
              </div>
            ) : members.length === 0 ? (
              <div className="p-8 text-center text-slate-500">
                Belum ada anggota dalam organisasi ini.
              </div>
            ) : (
              <table className="w-full text-sm text-left">
                <thead className="border-b border-slate-100 bg-slate-50/75 text-slate-600 font-semibold">
                  <tr>
                    <th className="px-5 py-3.5">Nama Pengguna</th>
                    <th className="px-5 py-3.5">Email</th>
                    <th className="px-5 py-3.5">Peran Organisasi</th>
                    <th className="px-5 py-3.5">Bergabung Pada</th>
                    <th className="px-5 py-3.5">Status</th>
                    <th className="px-5 py-3.5 text-right">Aksi</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100">
                  {members.map((m) => (
                    <tr key={m.id} className="hover:bg-slate-50/50 transition-colors">
                      <td className="px-5 py-3.5 font-medium text-slate-900">
                        {m.user_name || 'Tanpa Nama'}
                      </td>
                      <td className="px-5 py-3.5 text-slate-600">{m.user_email}</td>
                      <td className="px-5 py-3.5">
                        <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold ${
                          m.role === 'owner' || m.role === 'admin'
                            ? 'bg-purple-100 text-purple-700'
                            : 'bg-blue-50 text-blue-700'
                        }`}>
                          {m.role.toUpperCase()}
                        </span>
                      </td>
                      <td className="px-5 py-3.5 text-slate-500 text-xs">
                        {m.joined_at ? new Date(m.joined_at).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' }) : '—'}
                      </td>
                      <td className="px-5 py-3.5">
                        <span className="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-xs font-medium bg-emerald-50 text-emerald-700 border border-emerald-200">
                          <span className="w-1.5 h-1.5 rounded-full bg-emerald-500"></span>
                          Aktif
                        </span>
                      </td>
                      <td className="px-5 py-3.5 text-right">
                        <button
                          onClick={() => handleRemoveMember(m.user_id, m.user_name || m.user_email)}
                          title="Hapus dari Organisasi"
                          className="p-1.5 text-slate-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                        >
                          <Trash2 className="w-4 h-4" />
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}
          </div>
        </div>
      )}

      {/* ─────────────────────────────────────────────────────────────
          TAB 2: HAK AKSES PROPERTI (Multi-Level RBAC per Property)
      ────────────────────────────────────────────────────────────── */}
      {activeTab === 'property_rbac' && (
        <div className="space-y-5">
          <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3">
            <div>
              <h3 className="text-lg font-bold text-slate-900">Penugasan Staf & RBAC per Properti</h3>
              <p className="text-xs text-slate-500">
                Pilih properti untuk melihat atau menambahkan staf dengan peran default (Manajer, Keuangan, Resepsionis, Teknisi) maupun peran kustom.
              </p>
            </div>
            <Button
              onClick={() => {
                setAssignError(null);
                setIsAssignStaffOpen(true);
              }}
              disabled={!selectedPropertyId}
              className="bg-slate-900 text-white hover:bg-slate-800"
            >
              <UserPlus className="w-4 h-4 mr-2" />
              Tugaskan Staf ke Properti Ini
            </Button>
          </div>

          {/* Property Selector Bar */}
          <div className="flex items-center gap-3 bg-white p-4 rounded-2xl border border-slate-200 shadow-sm">
            <Building2 className="w-5 h-5 text-orange shrink-0" />
            <span className="text-sm font-semibold text-slate-700 shrink-0">Pilih Properti:</span>
            {propertiesLoading ? (
              <span className="text-xs text-slate-400">Memuat properti...</span>
            ) : properties.length === 0 ? (
              <span className="text-xs text-amber-600">Belum ada properti di organisasi ini.</span>
            ) : (
              <select
                value={selectedPropertyId}
                onChange={(e) => setSelectedPropertyId(e.target.value)}
                className="w-full sm:w-80 px-3 py-2 text-sm bg-slate-50 border border-slate-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-orange/30 focus:border-orange font-medium"
              >
                {properties.map((p) => (
                  <option key={p.id} value={p.id}>
                    {p.name} ({p.property_type || 'Properti'})
                  </option>
                ))}
              </select>
            )}
          </div>

          {/* Staff Table */}
          <div className="rounded-2xl border border-slate-200 bg-white overflow-hidden shadow-sm">
            <div className="px-5 py-4 border-b border-slate-100 flex items-center justify-between">
              <h4 className="text-sm font-bold text-slate-800">
                Staf Terdaftar pada Properti Ini
              </h4>
              <span className="text-xs text-slate-500">
                {propertyStaff.length} staf ditugaskan
              </span>
            </div>

            {staffLoading ? (
              <div className="p-8 text-center text-slate-400">
                <Loader2 className="w-6 h-6 animate-spin mx-auto mb-2 text-orange" />
                Memuat staf properti...
              </div>
            ) : propertyStaff.length === 0 ? (
              <div className="p-8 text-center">
                <p className="text-sm text-slate-500">Belum ada staf yang ditugaskan di properti ini.</p>
                <p className="text-xs text-slate-400 mt-1">Klik tombol &ldquo;Tugaskan Staf ke Properti Ini&rdquo; untuk menambahkan.</p>
              </div>
            ) : (
              <table className="w-full text-sm text-left">
                <thead className="border-b border-slate-100 bg-slate-50/75 text-slate-600 font-semibold">
                  <tr>
                    <th className="px-5 py-3.5">Nama Staf</th>
                    <th className="px-5 py-3.5">Email</th>
                    <th className="px-5 py-3.5">Peran Ditugaskan</th>
                    <th className="px-5 py-3.5">Tipe Peran</th>
                    <th className="px-5 py-3.5">Waktu Penugasan</th>
                    <th className="px-5 py-3.5 text-right">Aksi</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100">
                  {propertyStaff.map((staff) => (
                    <tr key={`${staff.user_id}-${staff.role_id}`} className="hover:bg-slate-50/50 transition-colors">
                      <td className="px-5 py-3.5 font-medium text-slate-900">
                        {staff.user_name || 'Tanpa Nama'}
                      </td>
                      <td className="px-5 py-3.5 text-slate-600">{staff.user_email}</td>
                      <td className="px-5 py-3.5">
                        <span className="font-semibold text-slate-800">
                          {staff.role_name}
                        </span>
                      </td>
                      <td className="px-5 py-3.5">
                        {staff.is_system ? (
                          <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-semibold bg-slate-100 text-slate-700">
                            Peran Bawaan
                          </span>
                        ) : (
                          <span className="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-semibold bg-amber-50 text-amber-700 border border-amber-200">
                            Peran Kustom
                          </span>
                        )}
                      </td>
                      <td className="px-5 py-3.5 text-slate-500 text-xs">
                        {staff.created_at ? new Date(staff.created_at).toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' }) : '—'}
                      </td>
                      <td className="px-5 py-3.5 text-right">
                        <button
                          onClick={() => handleRemoveStaff(staff.user_id, staff.role_id, staff.role_name)}
                          title="Cabut Peran Properti"
                          className="p-1.5 text-slate-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                        >
                          <Trash2 className="w-4 h-4" />
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}
          </div>
        </div>
      )}

      {/* ─────────────────────────────────────────────────────────────
          TAB 3: PERAN KUSTOM & IZIN (RBAC Role Builder)
      ────────────────────────────────────────────────────────────── */}
      {activeTab === 'custom_roles' && (
        <div className="space-y-5">
          <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3">
            <div>
              <h3 className="text-lg font-bold text-slate-900">Katalog Peran & Izin Akses</h3>
              <p className="text-xs text-slate-500">
                Lihat peran bawaan sistem dan buat peran kustom baru dengan izin granular per modul (Kamar, Tenant, Tagihan, Kontrak, dll).
              </p>
            </div>
            <Button
              onClick={() => {
                setEditingRole(null);
                setNewRoleName('');
                setNewRoleDescription('');
                setSelectedPermIds(new Set());
                setRoleError(null);
                setRoleModalOpen(true);
              }}
              className="bg-slate-900 text-white hover:bg-slate-800"
            >
              <Plus className="w-4 h-4 mr-2" />
              Buat Peran Kustom
            </Button>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {rolesLoading ? (
              <div className="col-span-full p-8 text-center text-slate-400">
                <Loader2 className="w-6 h-6 animate-spin mx-auto mb-2 text-orange" />
                Memuat peran...
              </div>
            ) : roles.map((role) => (
              <div
                key={role.id}
                className="bg-white rounded-2xl border border-slate-200 p-5 shadow-sm hover:shadow transition-shadow flex flex-col justify-between"
              >
                <div>
                  <div className="flex items-start justify-between gap-2 mb-2">
                    <div className="flex items-center gap-2">
                      <Shield className={`w-5 h-5 ${role.is_system ? 'text-slate-400' : 'text-orange'}`} />
                      <h4 className="font-bold text-slate-900">{role.name}</h4>
                    </div>
                    {role.is_system ? (
                      <span className="text-xs px-2 py-0.5 rounded-full bg-slate-100 text-slate-600 font-medium shrink-0">
                        Sistem
                      </span>
                    ) : (
                      <span className="text-xs px-2 py-0.5 rounded-full bg-orange/15 text-slate-900 font-semibold shrink-0">
                        Kustom
                      </span>
                    )}
                  </div>
                  <p className="text-xs text-slate-500 line-clamp-2 min-h-[32px]">
                    {role.description || 'Tidak ada deskripsi.'}
                  </p>
                  <div className="mt-3 text-xs text-slate-400">
                    {role.permissions?.length ?? 0} Izin Akses Terpasang
                  </div>
                </div>

                <div className="mt-4 pt-4 border-t border-slate-100 flex items-center justify-between">
                  <span className="text-xs text-slate-400 font-mono">#{role.id.slice(0, 8)}</span>
                  <div className="flex items-center gap-1">
                    <button
                      onClick={() => {
                        setEditingRole(role);
                        setNewRoleName(role.name);
                        setNewRoleDescription(role.description || '');
                        setSelectedPermIds(new Set(role.permissions?.map(p => p.id) ?? []));
                        setRoleError(null);
                        setRoleModalOpen(true);
                      }}
                      title="Lihat / Edit Izin"
                      className="p-1.5 text-slate-400 hover:text-slate-700 hover:bg-slate-100 rounded-lg transition-colors"
                    >
                      <Pencil className="w-4 h-4" />
                    </button>
                    {!role.is_system && (
                      <button
                        onClick={() => handleDeleteRole(role)}
                        title="Hapus Peran Kustom"
                        className="p-1.5 text-slate-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                      >
                        <Trash2 className="w-4 h-4" />
                      </button>
                    )}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* ─────────────────────────────────────────────────────────────
          TAB 4: KEBIJAKAN & NOTIFIKASI
      ────────────────────────────────────────────────────────────── */}
      {activeTab === 'policies' && (
        <form onSubmit={handleSavePolicies} className="space-y-6 max-w-2xl bg-white p-6 rounded-2xl border border-slate-200 shadow-sm">
          <div>
            <h3 className="text-lg font-bold text-slate-900">Kebijakan & Ambang Batas Peringatan</h3>
            <p className="text-xs text-slate-500">
              Konfigurasi ambang batas hari untuk tanda khusus dan template notifikasi WhatsApp ke tenant.
            </p>
          </div>

          <div className="space-y-4">
            <div>
              <label className="block text-sm font-semibold text-slate-700 mb-1">
                Mata Uang Utama Organisasi (Default Operating Currency)
              </label>
              <div className="flex items-center gap-3">
                <CurrencySelect
                  value={defaultCurrency}
                  onValueChange={setDefaultCurrency}
                  className="w-80"
                />
                <span className="text-xs text-slate-500">Mata uang default untuk seluruh properti, kamar, dan penagihan.</span>
              </div>
            </div>
            <div>
              <label className="block text-sm font-semibold text-slate-700 mb-1">
                Ambang Batas &ldquo;Segera Jatuh Tempo&rdquo; (Hari Sebelum)
              </label>
              <div className="flex items-center gap-2">
                <input
                  type="number"
                  min="1"
                  max="30"
                  value={dueSoonDays}
                  onChange={(e) => setDueSoonDays(e.target.value)}
                  className="w-32 px-3 py-2 text-sm bg-slate-50 border border-slate-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-orange/30 font-medium"
                />
                <span className="text-xs text-slate-500">hari sebelum tanggal jatuh tempo (default 7 hari)</span>
              </div>
            </div>

            <div>
              <label className="block text-sm font-semibold text-slate-700 mb-1">
                Ambang Batas &ldquo;Akan Habis Kontrak&rdquo; (Hari Sebelum)
              </label>
              <div className="flex items-center gap-2">
                <input
                  type="number"
                  min="1"
                  max="90"
                  value={expiringSoonDays}
                  onChange={(e) => setExpiringSoonDays(e.target.value)}
                  className="w-32 px-3 py-2 text-sm bg-slate-50 border border-slate-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-orange/30 font-medium"
                />
                <span className="text-xs text-slate-500">hari sebelum kontrak berakhir (default 30 hari)</span>
              </div>
            </div>

            {/* Konfigurasi Denda Keterlambatan Otomatis */}
            <div className="pt-4 border-t border-slate-100 space-y-4">
              <div className="flex items-center gap-2">
                <AlertTriangle className="w-5 h-5 text-amber-500" />
                <div>
                  <h4 className="text-sm font-bold text-slate-900">Aturan & Template Default Denda Keterlambatan</h4>
                  <p className="text-xs text-slate-500">
                    Atur perhitungan denda untuk tagihan sewa yang terlambat dibayar agar dapat digenerate otomatis.
                  </p>
                </div>
              </div>

              <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
                <div>
                  <label className="block text-xs font-semibold text-slate-700 mb-1">
                    Metode Perhitungan
                  </label>
                  <select
                    value={latePenaltyType}
                    onChange={(e) => setLatePenaltyType(e.target.value)}
                    className="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-orange/30 font-medium"
                  >
                    <option value="daily_fixed">Nominal Tetap / Hari</option>
                    <option value="daily_percent">Persentase / Hari (%)</option>
                    <option value="flat">Denda Sekali Dikenakan (Flat)</option>
                  </select>
                </div>

                <div>
                  <label className="block text-xs font-semibold text-slate-700 mb-1">
                    Besaran Denda {latePenaltyType === 'daily_percent' ? '(%)' : '(Nominal)'}
                  </label>
                  <input
                    type="number"
                    min="0"
                    step={latePenaltyType === 'daily_percent' ? '0.1' : '1000'}
                    value={latePenaltyAmount}
                    onChange={(e) => setLatePenaltyAmount(e.target.value)}
                    className="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-orange/30 font-medium"
                  />
                </div>

                <div>
                  <label className="block text-xs font-semibold text-slate-700 mb-1">
                    Masa Tenggang (Grace Period)
                  </label>
                  <div className="flex items-center gap-1.5">
                    <input
                      type="number"
                      min="0"
                      max="30"
                      value={lateGracePeriod}
                      onChange={(e) => setLateGracePeriod(e.target.value)}
                      className="w-full px-3 py-2 text-xs bg-slate-50 border border-slate-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-orange/30 font-medium"
                    />
                    <span className="text-xs text-slate-400 shrink-0">hari</span>
                  </div>
                </div>
              </div>

              {/* Automated Execution Card */}
              <div className="bg-amber-50/60 border border-amber-200 rounded-xl p-3.5 flex flex-col sm:flex-row sm:items-center justify-between gap-3">
                <div className="text-xs text-amber-900">
                  <p className="font-semibold flex items-center gap-1.5">
                    <Zap className="w-3.5 h-3.5 text-amber-600" /> Eksekusi Otomatisasi Denda
                  </p>
                  <p className="text-[11px] text-amber-700 mt-0.5">
                    Pindai seluruh tagihan belum lunas yang melewati {lateGracePeriod} hari masa tenggang dan terapkan denda otomatis.
                  </p>
                  {autoPenaltyResult && (
                    <p className="text-[11px] font-bold text-emerald-700 mt-1">
                      ✅ Berhasil membuat {autoPenaltyResult.success} denda baru ({autoPenaltyResult.skipped} tagihan tidak memenuhi syarat/sudah lunas).
                    </p>
                  )}
                </div>

                <Button
                  type="button"
                  onClick={handleAutoGeneratePenalties}
                  disabled={autoPenaltyProcessing}
                  className="bg-amber-600 hover:bg-amber-700 text-white font-semibold text-xs shrink-0 shadow-sm"
                >
                  {autoPenaltyProcessing ? (
                    <>
                      <Loader2 className="w-3.5 h-3.5 mr-1.5 animate-spin" /> Memproses...
                    </>
                  ) : (
                    <>
                      <Zap className="w-3.5 h-3.5 mr-1.5" /> Jalankan Denda Otomatis
                    </>
                  )}
                </Button>
              </div>
            </div>

            <div>
              <label className="block text-sm font-semibold text-slate-700 mb-1">
                Template Pesan Pengingat WhatsApp
              </label>
              <textarea
                rows={4}
                value={waTemplate}
                onChange={(e) => setWaTemplate(e.target.value)}
                className="w-full px-3 py-2 text-sm bg-slate-50 border border-slate-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-orange/30 font-sans"
              />
              <p className="text-xs text-slate-400 mt-1">
                Variabel yang tersedia: <code className="bg-slate-100 px-1 py-0.5 rounded">{'{tenant_name}'}</code>, <code className="bg-slate-100 px-1 py-0.5 rounded">{'{amount}'}</code>, <code className="bg-slate-100 px-1 py-0.5 rounded">{'{due_date}'}</code>
              </p>
            </div>
          </div>

          <div className="pt-4 border-t border-slate-100 flex items-center justify-between">
            {policySaved ? (
              <span className="text-xs font-semibold text-emerald-600 flex items-center gap-1">
                <Check className="w-4 h-4" /> Kebijakan berhasil disimpan!
              </span>
            ) : <div />}
            <Button type="submit" className="bg-slate-900 text-white hover:bg-slate-800">
              <Save className="w-4 h-4 mr-2" />
              Simpan Preferensi
            </Button>
          </div>
        </form>
      )}

      {/* ─────────────────────────────────────────────────────────────
          MODAL: TAMBAH ANGGOTA ORGANISASI
      ────────────────────────────────────────────────────────────── */}
      {isAddMemberOpen && (
        <div className="fixed inset-0 z-50 bg-black/40 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-white rounded-2xl max-w-md w-full p-6 shadow-2xl space-y-4">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <UserPlus className="w-5 h-5 text-orange" />
                <h3 className="text-lg font-bold text-slate-900">Undang / Tambah Anggota</h3>
              </div>
              <button
                onClick={() => setIsAddMemberOpen(false)}
                className="p-1 text-slate-400 hover:text-slate-600 rounded-lg"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            {memberError && (
              <div className="p-3 text-xs bg-red-50 text-red-700 border border-red-200 rounded-xl flex items-center gap-2">
                <AlertTriangle className="w-4 h-4 shrink-0" />
                <span>{memberError}</span>
              </div>
            )}

            <form onSubmit={handleAddMember} className="space-y-3">
              <div>
                <label className="block text-xs font-semibold text-slate-700 mb-1">
                  Email Anggota <span className="text-red-500">*</span>
                </label>
                <Input
                  type="email"
                  required
                  placeholder="staf@contoh.com"
                  value={newMemberEmail}
                  onChange={(e) => setNewMemberEmail(e.target.value)}
                />
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-700 mb-1">
                  Nama Lengkap (Opsional)
                </label>
                <Input
                  type="text"
                  placeholder="Budi Santoso"
                  value={newMemberName}
                  onChange={(e) => setNewMemberName(e.target.value)}
                />
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-700 mb-1">
                  Password Sementara (Opsional, bawaan: Password123!)
                </label>
                <Input
                  type="password"
                  placeholder="Minimal 8 karakter"
                  value={newMemberPassword}
                  onChange={(e) => setNewMemberPassword(e.target.value)}
                />
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-700 mb-1">
                  Peran Organisasi
                </label>
                <select
                  value={newMemberRole}
                  onChange={(e) => setNewMemberRole(e.target.value)}
                  className="w-full px-3 py-2 text-sm bg-slate-50 border border-slate-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-orange/30"
                >
                  <option value="member">Member (Staf Umum)</option>
                  <option value="admin">Admin (Administrator Organisasi)</option>
                </select>
              </div>

              <div className="pt-3 flex justify-end gap-2">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => setIsAddMemberOpen(false)}
                >
                  Batal
                </Button>
                <Button
                  type="submit"
                  disabled={addMemberLoading}
                  className="bg-slate-900 text-white hover:bg-slate-800"
                >
                  {addMemberLoading ? (
                    <>
                      <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                      Menyimpan...
                    </>
                  ) : (
                    'Tambahkan Anggota'
                  )}
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* ─────────────────────────────────────────────────────────────
          MODAL: TUGASKAN STAF KE PROPERTI
      ────────────────────────────────────────────────────────────── */}
      {isAssignStaffOpen && (
        <div className="fixed inset-0 z-50 bg-black/40 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-white rounded-2xl max-w-md w-full p-6 shadow-2xl space-y-4">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <Building2 className="w-5 h-5 text-orange" />
                <h3 className="text-lg font-bold text-slate-900">Tugaskan Staf ke Properti</h3>
              </div>
              <button
                onClick={() => setIsAssignStaffOpen(false)}
                className="p-1 text-slate-400 hover:text-slate-600 rounded-lg"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            {assignError && (
              <div className="p-3 text-xs bg-red-50 text-red-700 border border-red-200 rounded-xl flex items-center gap-2">
                <AlertTriangle className="w-4 h-4 shrink-0" />
                <span>{assignError}</span>
              </div>
            )}

            <form onSubmit={handleAssignStaff} className="space-y-3">
              <div>
                <label className="block text-xs font-semibold text-slate-700 mb-1">
                  Pilih Anggota Organisasi <span className="text-red-500">*</span>
                </label>
                <select
                  required
                  value={assignUserId}
                  onChange={(e) => setAssignUserId(e.target.value)}
                  className="w-full px-3 py-2 text-sm bg-slate-50 border border-slate-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-orange/30 font-medium"
                >
                  <option value="">-- Pilih Anggota --</option>
                  {members.map((m) => (
                    <option key={m.user_id} value={m.user_id}>
                      {m.user_name || 'Tanpa Nama'} ({m.user_email})
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-700 mb-1">
                  Pilih Peran Properti (Bawaan atau Kustom) <span className="text-red-500">*</span>
                </label>
                <select
                  required
                  value={assignRoleId}
                  onChange={(e) => setAssignRoleId(e.target.value)}
                  className="w-full px-3 py-2 text-sm bg-slate-50 border border-slate-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-orange/30 font-medium"
                >
                  <option value="">-- Pilih Peran --</option>
                  <optgroup label="Peran Bawaan Sistem">
                    {roles
                      .filter((r) => r.is_system)
                      .map((r) => (
                        <option key={r.id} value={r.id}>
                          {r.name}
                        </option>
                      ))}
                  </optgroup>
                  <optgroup label="Peran Kustom">
                    {roles
                      .filter((r) => !r.is_system)
                      .map((r) => (
                        <option key={r.id} value={r.id}>
                          {r.name}
                        </option>
                      ))}
                  </optgroup>
                </select>
              </div>

              <div className="pt-3 flex justify-end gap-2">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => setIsAssignStaffOpen(false)}
                >
                  Batal
                </Button>
                <Button
                  type="submit"
                  disabled={assignStaffLoading}
                  className="bg-slate-900 text-white hover:bg-slate-800"
                >
                  {assignStaffLoading ? (
                    <>
                      <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                      Menugaskan...
                    </>
                  ) : (
                    'Tugaskan Staf'
                  )}
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* ─────────────────────────────────────────────────────────────
          MODAL: BUAT / EDIT PERAN & IZIN (RBAC Role Builder)
      ────────────────────────────────────────────────────────────── */}
      {roleModalOpen && (
        <div className="fixed inset-0 z-50 bg-black/40 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-white rounded-2xl max-w-2xl w-full max-h-[85vh] flex flex-col shadow-2xl">
            <div className="p-6 border-b border-slate-100 flex items-center justify-between">
              <div>
                <h3 className="text-lg font-bold text-slate-900">
                  {editingRole ? `Konfigurasi Peran: ${editingRole.name}` : 'Buat Peran Kustom Baru'}
                </h3>
                <p className="text-xs text-slate-500">
                  {editingRole?.is_system
                    ? 'Peran sistem bawaan hanya dapat ditinjau izinnya.'
                    : 'Tentukan nama peran dan centang izin akses granular yang diberikan.'}
                </p>
              </div>
              <button
                onClick={() => setRoleModalOpen(false)}
                className="p-1 text-slate-400 hover:text-slate-600 rounded-lg"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            {roleError && (
              <div className="mx-6 mt-4 p-3 text-xs bg-red-50 text-red-700 border border-red-200 rounded-xl flex items-center gap-2">
                <AlertTriangle className="w-4 h-4 shrink-0" />
                <span>{roleError}</span>
              </div>
            )}

            <form onSubmit={handleSaveRole} className="p-6 flex-1 overflow-y-auto space-y-4">
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                <div>
                  <label className="block text-xs font-semibold text-slate-700 mb-1">
                    Nama Peran <span className="text-red-500">*</span>
                  </label>
                  <Input
                    required
                    disabled={editingRole?.is_system}
                    placeholder="Contoh: Supervisor Malam"
                    value={newRoleName}
                    onChange={(e) => setNewRoleName(e.target.value)}
                  />
                </div>
                <div>
                  <label className="block text-xs font-semibold text-slate-700 mb-1">
                    Deskripsi
                  </label>
                  <Input
                    disabled={editingRole?.is_system}
                    placeholder="Tugas & cakupan akses peran..."
                    value={newRoleDescription}
                    onChange={(e) => setNewRoleDescription(e.target.value)}
                  />
                </div>
              </div>

              <div>
                <label className="block text-xs font-semibold text-slate-700 mb-2">
                  Daftar Hak Akses (Permissions)
                </label>

                <div className="space-y-4 border border-slate-200 rounded-xl p-4 bg-slate-50/50">
                  {Object.entries(groupedPermissions).map(([resource, perms]) => (
                    <div key={resource} className="space-y-2">
                      <div className="flex items-center justify-between pb-1 border-b border-slate-200">
                        <span className="text-xs font-bold text-slate-800 uppercase tracking-wider">
                          Modul: {resource}
                        </span>
                        {!editingRole?.is_system && (
                          <button
                            type="button"
                            onClick={() => {
                              const allChecked = perms.every((p) => selectedPermIds.has(p.id));
                              setSelectedPermIds((prev) => {
                                const next = new Set(prev);
                                perms.forEach((p) => {
                                  if (allChecked) next.delete(p.id);
                                  else next.add(p.id);
                                });
                                return next;
                              });
                            }}
                            className="text-xs text-orange font-semibold hover:underline"
                          >
                            {perms.every((p) => selectedPermIds.has(p.id)) ? 'Hapus Semua' : 'Pilih Semua'}
                          </button>
                        )}
                      </div>

                      <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                        {perms.map((p) => {
                          const isChecked = selectedPermIds.has(p.id);
                          return (
                            <label
                              key={p.id}
                              className={`flex items-start gap-2.5 p-2 rounded-lg border text-xs cursor-pointer transition-colors ${
                                isChecked
                                  ? 'bg-orange/10 border-orange/30 text-slate-900'
                                  : 'bg-white border-slate-200 text-slate-600 hover:bg-slate-100'
                              } ${editingRole?.is_system ? 'cursor-not-allowed opacity-75' : ''}`}
                            >
                              <input
                                type="checkbox"
                                disabled={editingRole?.is_system}
                                checked={isChecked}
                                onChange={(e) => {
                                  setSelectedPermIds((prev) => {
                                    const next = new Set(prev);
                                    if (e.target.checked) next.add(p.id);
                                    else next.delete(p.id);
                                    return next;
                                  });
                                }}
                                className="mt-0.5 rounded border-slate-300 text-orange focus:ring-orange"
                              />
                              <div>
                                <p className="font-semibold">{p.key}</p>
                                <p className="text-[11px] text-slate-400">{p.description || p.action}</p>
                              </div>
                            </label>
                          );
                        })}
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              <div className="pt-3 border-t border-slate-100 flex justify-end gap-2">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => setRoleModalOpen(false)}
                >
                  Tutup
                </Button>
                {!editingRole?.is_system && (
                  <Button
                    type="submit"
                    disabled={savingRole}
                    className="bg-slate-900 text-white hover:bg-slate-800"
                  >
                    {savingRole ? (
                      <>
                        <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                        Menyimpan...
                      </>
                    ) : (
                      'Simpan Peran'
                    )}
                  </Button>
                )}
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}

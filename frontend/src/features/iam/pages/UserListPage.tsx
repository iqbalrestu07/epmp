import { useState, useEffect } from 'react';
import { Plus, Loader2, X, Pencil, Trash2, ShieldCheck } from 'lucide-react';
import { userService, roleService } from '../services/iamService';
import type { User, Role, CreateUserRequest } from '../types';
import PermissionGuard from '../../../components/PermissionGuard';

// ── Invite User Modal ────────────────────────────────────────────────────────

function InviteUserModal({
  allRoles,
  onClose,
  onSave,
}: {
  allRoles: Role[];
  onClose: () => void;
  onSave: () => void;
}) {
  const [form, setForm] = useState<CreateUserRequest>({ name: '', email: '', password: '', role_ids: [] });
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const toggleRole = (id: string) => {
    setForm(prev => ({
      ...prev,
      role_ids: prev.role_ids?.includes(id)
        ? prev.role_ids.filter(r => r !== id)
        : [...(prev.role_ids ?? []), id],
    }));
  };

  const handleSave = async () => {
    if (!form.name || !form.email || !form.password) {
      setError('Name, email, and password are required.');
      return;
    }
    setSaving(true);
    setError(null);
    try {
      await userService.create(form);
      onSave();
      onClose();
    } catch (err: any) {
      setError(err?.message ?? 'Failed to create user.');
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4">
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-md animate-in fade-in zoom-in-95 duration-200">
        <div className="flex items-center justify-between p-6 border-b">
          <h2 className="text-xl font-semibold">Invite User</h2>
          <button onClick={onClose} className="text-gray-400 hover:text-black transition-colors">
            <X size={20} />
          </button>
        </div>

        <div className="p-6 space-y-4">
          {error && (
            <div className="px-4 py-3 rounded-xl bg-red-50 border border-red-200 text-red-600 text-sm">
              {error}
            </div>
          )}

          <div className="flex flex-col gap-1.5">
            <label className="text-sm font-medium text-gray-700">Full Name</label>
            <input
              value={form.name}
              onChange={e => setForm(p => ({ ...p, name: e.target.value }))}
              placeholder="John Smith"
              className="border rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-orange/30 focus:border-orange transition-all"
            />
          </div>

          <div className="flex flex-col gap-1.5">
            <label className="text-sm font-medium text-gray-700">Email</label>
            <input
              value={form.email}
              onChange={e => setForm(p => ({ ...p, email: e.target.value }))}
              type="email"
              placeholder="john@acme.com"
              className="border rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-orange/30 focus:border-orange transition-all"
            />
          </div>

          <div className="flex flex-col gap-1.5">
            <label className="text-sm font-medium text-gray-700">Initial Password</label>
            <input
              value={form.password}
              onChange={e => setForm(p => ({ ...p, password: e.target.value }))}
              type="password"
              placeholder="Min. 8 characters"
              className="border rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-orange/30 focus:border-orange transition-all"
            />
          </div>

          <div className="flex flex-col gap-1.5">
            <label className="text-sm font-medium text-gray-700">Assign Roles</label>
            <div className="border rounded-xl overflow-hidden divide-y">
              {allRoles.map(role => (
                <label key={role.id} className="flex items-center gap-3 px-4 py-3 hover:bg-gray-50 cursor-pointer transition-colors">
                  <div
                    className={`w-5 h-5 rounded flex items-center justify-center border-2 transition-colors ${
                      form.role_ids?.includes(role.id) ? 'bg-orange border-orange' : 'border-gray-300'
                    }`}
                    onClick={() => toggleRole(role.id)}
                  >
                    {form.role_ids?.includes(role.id) && (
                      <svg className="w-3 h-3 text-black" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={3}>
                        <path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" />
                      </svg>
                    )}
                  </div>
                  <div>
                    <p className="text-sm font-medium">{role.name}</p>
                    {role.description && <p className="text-xs text-gray-500">{role.description}</p>}
                  </div>
                </label>
              ))}
            </div>
          </div>
        </div>

        <div className="p-6 border-t flex justify-end gap-3">
          <button onClick={onClose} className="px-5 py-2 rounded-xl border text-sm font-medium hover:bg-gray-50 transition-colors">
            Cancel
          </button>
          <button
            onClick={handleSave}
            disabled={saving}
            className="px-5 py-2 rounded-xl bg-black text-white text-sm font-medium hover:bg-black/80 transition-colors disabled:opacity-60 flex items-center gap-2"
          >
            {saving && <Loader2 size={14} className="animate-spin" />}
            {saving ? 'Creating…' : 'Create User'}
          </button>
        </div>
      </div>
    </div>
  );
}

// ── Main Page ──────────────────────────────────────────────────────────────

export default function UserListPage() {
  const [users, setUsers] = useState<User[]>([]);
  const [allRoles, setAllRoles] = useState<Role[]>([]);
  const [loading, setLoading] = useState(true);
  const [modalOpen, setModalOpen] = useState(false);
  const [deleteConfirm, setDeleteConfirm] = useState<string | null>(null);

  const fetchData = async () => {
    setLoading(true);
    try {
      const [usersRes, rolesRes] = await Promise.all([
        userService.list(),
        roleService.list(),
      ]);
      setUsers(usersRes.data.data ?? []);
      setAllRoles(rolesRes.data ?? []);
    } catch {
      // keep empty
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { fetchData(); }, []);

  const handleDeactivate = async (id: string) => {
    try {
      await userService.delete(id);
      await fetchData();
    } catch (err: any) {
      alert(err?.message ?? 'Failed to deactivate user.');
    } finally {
      setDeleteConfirm(null);
    }
  };

  const getInitials = (name: string) =>
    name.split(' ').map(w => w[0]).join('').toUpperCase().slice(0, 2);

  return (
    <div className="animate-in fade-in slide-in-from-bottom-4 duration-500">
      {/* Header */}
      <div className="flex justify-between items-start mb-8">
        <div>
          <h1 className="text-3xl font-display mb-2">User Accounts</h1>
          <p className="text-gray-500">Manage team members in this workspace.</p>
        </div>
        <PermissionGuard permission="user:write">
          <button
            onClick={() => setModalOpen(true)}
            className="bg-black text-white px-4 py-2 rounded-lg flex items-center gap-2 hover:bg-black/80 transition-colors"
          >
            <Plus size={18} />
            Invite User
          </button>
        </PermissionGuard>
      </div>

      {/* Table */}
      <div className="bg-white rounded-2xl border shadow-sm overflow-hidden">
        {loading ? (
          <div className="flex items-center justify-center py-20">
            <Loader2 size={28} className="animate-spin text-orange" />
          </div>
        ) : users.length === 0 ? (
          <div className="text-center py-20 text-gray-500">No users found.</div>
        ) : (
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-gray-50 border-b">
                <th className="px-6 py-4 font-semibold text-sm text-gray-600">User</th>
                <th className="px-6 py-4 font-semibold text-sm text-gray-600">Roles</th>
                <th className="px-6 py-4 font-semibold text-sm text-gray-600">Status</th>
                <th className="px-6 py-4 font-semibold text-sm text-gray-600 text-right">Actions</th>
              </tr>
            </thead>
            <tbody>
              {users.map((u) => (
                <tr key={u.id} className="border-b hover:bg-gray-50/50 transition-colors">
                  <td className="px-6 py-4">
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 rounded-full bg-orange/10 flex items-center justify-center font-bold text-sm text-orange">
                        {getInitials(u.name)}
                      </div>
                      <div>
                        <p className="font-medium">{u.name}</p>
                        <p className="text-xs text-gray-500">{u.email}</p>
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-4">
                    <div className="flex flex-wrap gap-1">
                      {u.roles.length > 0 ? u.roles.map(r => (
                        <span key={r.id} className="inline-flex items-center gap-1 px-2 py-0.5 bg-gray-100 text-gray-700 rounded-full text-xs font-medium">
                          <ShieldCheck size={10} className="text-orange" />
                          {r.name}
                        </span>
                      )) : (
                        <span className="text-xs text-gray-400 italic">No roles</span>
                      )}
                    </div>
                  </td>
                  <td className="px-6 py-4">
                    <span className={`inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-medium ${
                      u.is_active
                        ? 'bg-green-100 text-green-700'
                        : 'bg-red-100 text-red-600'
                    }`}>
                      <span className={`w-1.5 h-1.5 rounded-full ${u.is_active ? 'bg-green-500' : 'bg-red-400'}`} />
                      {u.is_active ? 'Active' : 'Inactive'}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-right">
                    <div className="flex items-center gap-2 justify-end">
                      <PermissionGuard permission="user:write">
                        <button className="p-2 rounded-lg text-gray-500 hover:bg-gray-100 hover:text-black transition-colors" title="Edit user">
                          <Pencil size={16} />
                        </button>
                      </PermissionGuard>
                      <PermissionGuard permission="user:delete">
                        {deleteConfirm === u.id ? (
                          <div className="flex items-center gap-1">
                            <button
                              onClick={() => handleDeactivate(u.id)}
                              className="px-3 py-1.5 bg-red-500 text-white text-xs rounded-lg hover:bg-red-600 transition-colors"
                            >Confirm</button>
                            <button
                              onClick={() => setDeleteConfirm(null)}
                              className="px-3 py-1.5 bg-gray-100 text-xs rounded-lg hover:bg-gray-200 transition-colors"
                            >Cancel</button>
                          </div>
                        ) : (
                          <button
                            onClick={() => setDeleteConfirm(u.id)}
                            className="p-2 rounded-lg text-gray-500 hover:bg-red-50 hover:text-red-500 transition-colors"
                            title="Deactivate user"
                          >
                            <Trash2 size={16} />
                          </button>
                        )}
                      </PermissionGuard>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      {/* Modal */}
      {modalOpen && (
        <InviteUserModal
          allRoles={allRoles}
          onClose={() => setModalOpen(false)}
          onSave={fetchData}
        />
      )}
    </div>
  );
}

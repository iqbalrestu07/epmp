import { useState, useEffect } from 'react';
import { ShieldCheck, Plus, Pencil, Trash2, X, Check, Loader2 } from 'lucide-react';
import { roleService, permissionService } from '../services/iamService';
import type { Role, Permission, CreateRoleRequest } from '../types';
import PermissionGuard from '../../../components/PermissionGuard';

// ── Permission checkbox group ───────────────────────────────────────────────

function PermissionToggle({
  permission,
  checked,
  onChange,
}: {
  permission: Permission;
  checked: boolean;
  onChange: (id: string, value: boolean) => void;
}) {
  return (
    <label className="flex items-center gap-3 p-3 rounded-xl hover:bg-black/5 cursor-pointer transition-colors">
      <div
        className={`w-5 h-5 rounded-md border-2 flex items-center justify-center transition-colors ${
          checked ? 'bg-orange border-orange' : 'border-gray-300'
        }`}
        onClick={() => onChange(permission.id, !checked)}
      >
        {checked && <Check size={12} className="text-black" strokeWidth={3} />}
      </div>
      <div>
        <p className="text-sm font-medium">{permission.key}</p>
        <p className="text-xs text-gray-500">{permission.description}</p>
      </div>
    </label>
  );
}

// ── Create / Edit Role Modal ─────────────────────────────────────────────────

function RoleModal({
  existingRole,
  allPermissions,
  onClose,
  onSave,
}: {
  existingRole?: Role | null;
  allPermissions: Permission[];
  onClose: () => void;
  onSave: () => void;
}) {
  const [name, setName] = useState(existingRole?.name ?? '');
  const [description, setDescription] = useState(existingRole?.description ?? '');
  const [selectedPerms, setSelectedPerms] = useState<Set<string>>(
    new Set(existingRole?.permissions?.map(p => p.id) ?? [])
  );
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Group permissions by resource
  const grouped = allPermissions.reduce<Record<string, Permission[]>>((acc, p) => {
    if (!acc[p.resource]) acc[p.resource] = [];
    acc[p.resource].push(p);
    return acc;
  }, {});

  const togglePerm = (id: string, value: boolean) => {
    setSelectedPerms(prev => {
      const next = new Set(prev);
      if (value) next.add(id);
      else next.delete(id);
      return next;
    });
  };

  const handleSave = async () => {
    if (!name.trim()) {
      setError('Role name is required.');
      return;
    }
    setSaving(true);
    setError(null);
    try {
      if (existingRole) {
        await roleService.update(existingRole.id, { name, description });
        if (!existingRole.is_system) {
          await roleService.setPermissions(existingRole.id, {
            permission_ids: Array.from(selectedPerms),
          });
        }
      } else {
        const req: CreateRoleRequest = {
          name,
          description,
          permission_ids: Array.from(selectedPerms),
        };
        await roleService.create(req);
      }
      onSave();
      onClose();
    } catch (err: any) {
      setError(err?.message ?? 'Failed to save role.');
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4">
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-2xl max-h-[90vh] flex flex-col overflow-hidden animate-in fade-in zoom-in-95 duration-200">
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b">
          <h2 className="text-xl font-semibold">
            {existingRole ? 'Edit Role' : 'Create New Role'}
          </h2>
          <button onClick={onClose} className="text-gray-400 hover:text-black transition-colors">
            <X size={20} />
          </button>
        </div>

        {/* Body */}
        <div className="overflow-y-auto flex-1 p-6 space-y-6">
          {error && (
            <div className="px-4 py-3 rounded-xl bg-red-50 border border-red-200 text-red-600 text-sm">
              {error}
            </div>
          )}

          <div className="space-y-4">
            <div className="flex flex-col gap-2">
              <label className="text-sm font-medium text-gray-700">Role Name</label>
              <input
                value={name}
                onChange={e => setName(e.target.value)}
                disabled={existingRole?.is_system}
                placeholder="e.g. Property Inspector"
                className="border rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-orange/30 focus:border-orange transition-all disabled:bg-gray-50 disabled:text-gray-400"
              />
            </div>
            <div className="flex flex-col gap-2">
              <label className="text-sm font-medium text-gray-700">Description</label>
              <input
                value={description}
                onChange={e => setDescription(e.target.value)}
                disabled={existingRole?.is_system}
                placeholder="What can this role do?"
                className="border rounded-xl px-4 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-orange/30 focus:border-orange transition-all disabled:bg-gray-50 disabled:text-gray-400"
              />
            </div>
          </div>

          {/* Permissions */}
          <div>
            <h3 className="text-sm font-semibold text-gray-700 mb-3">Permissions</h3>
            {existingRole?.is_system ? (
              <p className="text-sm text-gray-500 italic">System roles have all permissions by default and cannot be modified.</p>
            ) : (
              <div className="space-y-4">
                {Object.entries(grouped).map(([resource, perms]) => (
                  <div key={resource} className="border rounded-xl overflow-hidden">
                    <div className="bg-gray-50 px-4 py-2 border-b">
                      <span className="text-xs font-bold uppercase tracking-wider text-gray-600">{resource}</span>
                    </div>
                    <div className="divide-y">
                      {perms.map(p => (
                        <PermissionToggle
                          key={p.id}
                          permission={p}
                          checked={selectedPerms.has(p.id)}
                          onChange={togglePerm}
                        />
                      ))}
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>

        {/* Footer */}
        <div className="p-6 border-t flex justify-end gap-3">
          <button onClick={onClose} className="px-5 py-2 rounded-xl border text-sm font-medium hover:bg-gray-50 transition-colors">
            Cancel
          </button>
          {!existingRole?.is_system && (
            <button
              onClick={handleSave}
              disabled={saving}
              className="px-5 py-2 rounded-xl bg-black text-white text-sm font-medium hover:bg-black/80 transition-colors disabled:opacity-60 flex items-center gap-2"
            >
              {saving && <Loader2 size={14} className="animate-spin" />}
              {saving ? 'Saving…' : 'Save Role'}
            </button>
          )}
        </div>
      </div>
    </div>
  );
}

// ── Main Page ──────────────────────────────────────────────────────────────

export default function RBACPage() {
  const [roles, setRoles] = useState<Role[]>([]);
  const [allPermissions, setAllPermissions] = useState<Permission[]>([]);
  const [loading, setLoading] = useState(true);
  const [modalOpen, setModalOpen] = useState(false);
  const [editingRole, setEditingRole] = useState<Role | null>(null);
  const [deleteConfirm, setDeleteConfirm] = useState<string | null>(null);

  const fetchData = async () => {
    setLoading(true);
    try {
      const [rolesRes, permsRes] = await Promise.all([
        roleService.list(),
        permissionService.list(),
      ]);
      // Fetch permissions for each role
      const rolesWithPerms = await Promise.all(
        (rolesRes.data ?? []).map(async (role: Role) => {
          try {
            const full = await roleService.getById(role.id);
            return full.data;
          } catch {
            return role;
          }
        })
      );
      setRoles(rolesWithPerms);
      setAllPermissions(permsRes.data ?? []);
    } catch {
      // keep empty
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { fetchData(); }, []);

  const handleDelete = async (id: string) => {
    try {
      await roleService.delete(id);
      await fetchData();
    } catch (err: any) {
      alert(err?.message ?? 'Failed to delete role.');
    } finally {
      setDeleteConfirm(null);
    }
  };

  return (
    <div className="animate-in fade-in slide-in-from-bottom-4 duration-500">
      {/* Header */}
      <div className="flex justify-between items-start mb-8">
        <div>
          <h1 className="text-3xl font-display mb-2">Roles & Permissions</h1>
          <p className="text-gray-500">Manage what your team members can see and do.</p>
        </div>
        <PermissionGuard permission="role:write">
          <button
            onClick={() => { setEditingRole(null); setModalOpen(true); }}
            className="bg-black text-white px-4 py-2 rounded-lg flex items-center gap-2 hover:bg-black/80 transition-colors"
          >
            <Plus size={18} />
            Create Role
          </button>
        </PermissionGuard>
      </div>

      {/* Table */}
      <div className="bg-white rounded-2xl border shadow-sm overflow-hidden">
        {loading ? (
          <div className="flex items-center justify-center py-20">
            <Loader2 size={28} className="animate-spin text-orange" />
          </div>
        ) : roles.length === 0 ? (
          <div className="text-center py-20 text-gray-500">No roles found.</div>
        ) : (
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-gray-50 border-b">
                <th className="px-6 py-4 font-semibold text-sm text-gray-600">Role Name</th>
                <th className="px-6 py-4 font-semibold text-sm text-gray-600">Type</th>
                <th className="px-6 py-4 font-semibold text-sm text-gray-600">Permissions</th>
                <th className="px-6 py-4 font-semibold text-sm text-gray-600 text-right">Actions</th>
              </tr>
            </thead>
            <tbody>
              {roles.map((r) => (
                <tr key={r.id} className="border-b hover:bg-gray-50/50 transition-colors">
                  <td className="px-6 py-4">
                    <div className="flex items-center gap-3">
                      <ShieldCheck size={18} className="text-orange flex-shrink-0" />
                      <div>
                        <p className="font-medium">{r.name}</p>
                        {r.description && <p className="text-xs text-gray-500">{r.description}</p>}
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-4">
                    {r.is_system ? (
                      <span className="px-2 py-1 bg-orange/10 text-orange rounded-full text-xs font-medium">System</span>
                    ) : (
                      <span className="px-2 py-1 bg-gray-100 text-gray-600 rounded-full text-xs font-medium">Custom</span>
                    )}
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-600">
                    {r.is_system ? (
                      <span className="text-orange font-medium">All Access</span>
                    ) : (
                      <span>{r.permissions?.length ?? 0} permissions</span>
                    )}
                  </td>
                  <td className="px-6 py-4 text-right">
                    <div className="flex items-center gap-2 justify-end">
                      <PermissionGuard permission="role:write">
                        <button
                          onClick={() => { setEditingRole(r); setModalOpen(true); }}
                          className="p-2 rounded-lg text-gray-500 hover:bg-gray-100 hover:text-black transition-colors"
                          title="Edit role"
                        >
                          <Pencil size={16} />
                        </button>
                      </PermissionGuard>
                      {!r.is_system && (
                        <PermissionGuard permission="role:delete">
                          {deleteConfirm === r.id ? (
                            <div className="flex items-center gap-1">
                              <button
                                onClick={() => handleDelete(r.id)}
                                className="px-3 py-1.5 bg-red-500 text-white text-xs rounded-lg hover:bg-red-600 transition-colors"
                              >Confirm</button>
                              <button
                                onClick={() => setDeleteConfirm(null)}
                                className="px-3 py-1.5 bg-gray-100 text-xs rounded-lg hover:bg-gray-200 transition-colors"
                              >Cancel</button>
                            </div>
                          ) : (
                            <button
                              onClick={() => setDeleteConfirm(r.id)}
                              className="p-2 rounded-lg text-gray-500 hover:bg-red-50 hover:text-red-500 transition-colors"
                              title="Delete role"
                            >
                              <Trash2 size={16} />
                            </button>
                          )}
                        </PermissionGuard>
                      )}
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
        <RoleModal
          existingRole={editingRole}
          allPermissions={allPermissions}
          onClose={() => { setModalOpen(false); setEditingRole(null); }}
          onSave={fetchData}
        />
      )}
    </div>
  );
}

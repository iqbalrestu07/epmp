import { ShieldCheck, Plus } from 'lucide-react';

export default function RBACPage() {
  const roles = [
    { name: 'Super Administrator', users: 2, permissions: 'All Access' },
    { name: 'Property Manager', users: 5, permissions: 'Properties, Rooms, Tenants (Write)' },
    { name: 'Finance Staff', users: 3, permissions: 'Invoices, Payments, Reports (Write)' },
    { name: 'Maintenance', users: 12, permissions: 'Work Orders (Write), Assets (Read)' },
  ];

  return (
    <div className="animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div className="flex justify-between items-center mb-8">
        <div>
          <h1 className="text-3xl font-display mb-2">Roles & Permissions</h1>
          <p className="text-gray-500">Manage what your team members can see and do.</p>
        </div>
        <button className="bg-black text-white px-4 py-2 rounded-lg flex items-center gap-2 hover:bg-black/80 transition-colors">
          <Plus size={18} />
          Create Role
        </button>
      </div>

      <div className="bg-white rounded-2xl border shadow-sm overflow-hidden">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="bg-gray-50 border-b">
              <th className="px-6 py-4 font-semibold text-sm text-gray-600">Role Name</th>
              <th className="px-6 py-4 font-semibold text-sm text-gray-600">Active Users</th>
              <th className="px-6 py-4 font-semibold text-sm text-gray-600">Permission Scope</th>
              <th className="px-6 py-4 font-semibold text-sm text-gray-600 text-right">Actions</th>
            </tr>
          </thead>
          <tbody>
            {roles.map((r, i) => (
              <tr key={i} className="border-b hover:bg-gray-50/50">
                <td className="px-6 py-4">
                  <div className="flex items-center gap-3">
                    <ShieldCheck size={18} className="text-orange" />
                    <span className="font-medium">{r.name}</span>
                  </div>
                </td>
                <td className="px-6 py-4 text-sm text-gray-600">{r.users} users</td>
                <td className="px-6 py-4 text-sm text-gray-600">{r.permissions}</td>
                <td className="px-6 py-4 text-right">
                  <button className="text-orange text-sm font-medium hover:underline">Edit</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

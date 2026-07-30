import { Plus } from 'lucide-react';

export default function UserListPage() {
  const users = [
    { name: 'Admin Acme', email: 'admin@acme.com', role: 'Super Administrator', status: 'Active' },
    { name: 'John Manager', email: 'john@acme.com', role: 'Property Manager', status: 'Active' },
    { name: 'Sarah Finance', email: 'sarah@acme.com', role: 'Finance Staff', status: 'Invited' },
  ];

  return (
    <div className="animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div className="flex justify-between items-center mb-8">
        <div>
          <h1 className="text-3xl font-display mb-2">User Accounts</h1>
          <p className="text-gray-500">Manage team members in this workspace.</p>
        </div>
        <button className="bg-black text-white px-4 py-2 rounded-lg flex items-center gap-2 hover:bg-black/80 transition-colors">
          <Plus size={18} />
          Invite User
        </button>
      </div>

      <div className="bg-white rounded-2xl border shadow-sm overflow-hidden">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="bg-gray-50 border-b">
              <th className="px-6 py-4 font-semibold text-sm text-gray-600">User</th>
              <th className="px-6 py-4 font-semibold text-sm text-gray-600">Role</th>
              <th className="px-6 py-4 font-semibold text-sm text-gray-600">Status</th>
              <th className="px-6 py-4 font-semibold text-sm text-gray-600 text-right">Actions</th>
            </tr>
          </thead>
          <tbody>
            {users.map((u, i) => (
              <tr key={i} className="border-b hover:bg-gray-50/50">
                <td className="px-6 py-4">
                  <div className="flex items-center gap-3">
                    <div className="w-8 h-8 rounded-full bg-gray-200 flex items-center justify-center font-bold text-xs">
                      {u.name.charAt(0)}
                    </div>
                    <div>
                      <p className="font-medium">{u.name}</p>
                      <p className="text-xs text-gray-500">{u.email}</p>
                    </div>
                  </div>
                </td>
                <td className="px-6 py-4 text-sm text-gray-600">{u.role}</td>
                <td className="px-6 py-4">
                  <span className={`px-2 py-1 rounded-full text-xs font-medium ${u.status === 'Active' ? 'bg-green-100 text-green-700' : 'bg-yellow-100 text-yellow-700'}`}>
                    {u.status}
                  </span>
                </td>
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

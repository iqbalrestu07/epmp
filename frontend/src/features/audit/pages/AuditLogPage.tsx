import { useState } from "react";
import { Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useAuditLogs } from "../hooks";
import type { AuditLog } from "../types";
import { ChevronDown, ChevronRight } from "lucide-react";

const ACTION_STYLES: Record<string, string> = {
  CREATE: "bg-emerald-100 text-emerald-800",
  UPDATE: "bg-amber-100 text-amber-800",
  DELETE: "bg-red-100 text-red-800",
};

function AuditRow({ log }: { log: AuditLog }) {
  const [open, setOpen] = useState(false);
  return (
    <>
      <tr className="border-t hover:bg-slate-50 cursor-pointer" onClick={() => setOpen(!open)}>
        <td className="px-4 py-3">
          {open ? <ChevronDown className="w-4 h-4 text-slate-400" /> : <ChevronRight className="w-4 h-4 text-slate-400" />}
        </td>
        <td className="px-4 py-3 text-slate-600 whitespace-nowrap">
          {new Date(log.created_at).toLocaleString("id-ID")}
        </td>
        <td className="px-4 py-3 text-slate-900">{log.user_email || "-"}</td>
        <td className="px-4 py-3">
          <span className={`px-2 py-0.5 rounded text-xs font-medium ${ACTION_STYLES[log.action] ?? "bg-slate-100 text-slate-700"}`}>
            {log.action}
          </span>
        </td>
        <td className="px-4 py-3 font-medium text-slate-900">{log.module}</td>
        <td className="px-4 py-3 text-slate-500 font-mono text-xs">{log.entity_id || "-"}</td>
        <td className="px-4 py-3 text-right">
          <span className={`text-xs font-mono ${log.status_code < 400 ? "text-emerald-600" : "text-red-600"}`}>
            {log.status_code}
          </span>
        </td>
      </tr>
      {open && (
        <tr className="border-t bg-slate-50/60">
          <td colSpan={7} className="px-6 py-3">
            <div className="text-xs text-slate-500 space-y-1">
              <p><span className="font-medium">Path:</span> <span className="font-mono">{log.method} {log.path}</span></p>
              <p><span className="font-medium">IP:</span> {log.ip_address || "-"} · <span className="font-medium">UA:</span> {log.user_agent || "-"}</p>
              {log.request_body && (
                <pre className="mt-2 rounded bg-slate-900 text-slate-100 p-3 overflow-x-auto">
                  {JSON.stringify(log.request_body, null, 2)}
                </pre>
              )}
            </div>
          </td>
        </tr>
      )}
    </>
  );
}

export function AuditLogPage() {
  const [page, setPage] = useState(1);
  const [search, setSearch] = useState("");
  const [searchInput, setSearchInput] = useState("");
  const [action, setAction] = useState("");

  const { data, isLoading } = useAuditLogs({
    page,
    per_page: 20,
    search: search || undefined,
    action: action || undefined,
  });

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setSearch(searchInput);
    setPage(1);
  };

  return (
    <div className="space-y-6">
      <div>
        <div className="flex items-center gap-2 text-sm text-slate-500 mb-1">
          <Link to="/dashboard" className="hover:text-slate-900">Dashboard</Link>
          <span>/</span>
          <span className="text-slate-900 font-medium">Audit Log</span>
        </div>
        <h1 className="text-2xl font-bold text-slate-900">Audit Log</h1>
        <p className="text-sm text-slate-500 mt-1">Jejak seluruh operasi perubahan data pada organisasi ini.</p>
      </div>

      <form onSubmit={handleSearch} className="flex flex-wrap gap-2">
        <Input
          placeholder="Cari module, user, path..."
          value={searchInput}
          onChange={(e) => setSearchInput(e.target.value)}
          className="max-w-sm"
        />
        <select
          className="h-10 rounded-md border border-input bg-background px-3 text-sm"
          value={action}
          onChange={(e) => { setAction(e.target.value); setPage(1); }}
        >
          <option value="">Semua Aksi</option>
          <option value="CREATE">CREATE</option>
          <option value="UPDATE">UPDATE</option>
          <option value="DELETE">DELETE</option>
        </select>
        <Button type="submit" variant="outline">Search</Button>
      </form>

      {isLoading ? (
        <div className="text-center py-12 text-slate-400">Memuat audit log...</div>
      ) : (
        <>
          <div className="rounded-lg border bg-white shadow-sm overflow-hidden">
            <table className="w-full text-sm">
              <thead className="bg-slate-50 text-slate-600">
                <tr>
                  <th className="px-4 py-3 w-8"></th>
                  <th className="text-left px-4 py-3 font-medium">Waktu</th>
                  <th className="text-left px-4 py-3 font-medium">User</th>
                  <th className="text-left px-4 py-3 font-medium">Aksi</th>
                  <th className="text-left px-4 py-3 font-medium">Module</th>
                  <th className="text-left px-4 py-3 font-medium">Entity ID</th>
                  <th className="text-right px-4 py-3 font-medium">Status</th>
                </tr>
              </thead>
              <tbody>
                {(data?.data ?? []).length === 0 && (
                  <tr><td colSpan={7} className="px-4 py-8 text-center text-slate-400">Belum ada aktivitas tercatat</td></tr>
                )}
                {(data?.data ?? []).map((log) => <AuditRow key={log.id} log={log} />)}
              </tbody>
            </table>
          </div>

          {data && data.total_pages > 1 && (
            <div className="flex items-center justify-between">
              <p className="text-sm text-slate-500">
                Page {data.page} of {data.total_pages} ({data.total} total)
              </p>
              <div className="flex gap-2">
                <Button variant="outline" size="sm" disabled={page <= 1} onClick={() => setPage(page - 1)}>
                  Previous
                </Button>
                <Button variant="outline" size="sm" disabled={page >= data.total_pages} onClick={() => setPage(page + 1)}>
                  Next
                </Button>
              </div>
            </div>
          )}
        </>
      )}
    </div>
  );
}

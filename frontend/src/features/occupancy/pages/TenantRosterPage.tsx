import { TenantRoomRoster } from '../components/TenantRoomRoster';

export function TenantRosterPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-slate-900">Direktori Kamar & Hunian (Roster)</h1>
        <p className="text-sm text-slate-500 mt-1">
          Ikhtisar terpadu relasi antara unit kamar, penyewa aktif, periode sewa kontrak, serta status tagihan secara real-time.
        </p>
      </div>

      <TenantRoomRoster />
    </div>
  );
}

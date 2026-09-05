import { useState, useEffect } from "react";
import {
  Send,
  Users,
  ChevronRight,
  Clock,
  CheckCircle2,
  XCircle,
  Loader2,
  Eye,
  History,
  MessageSquare,
  X,
  Phone,
  AlertCircle,
} from "lucide-react";
import {
  useBlasts,
  useCreateBlast,
  useSendBlast,
  useBlastLogs,
  useTemplates,
  useDevices,
  useRecipientPreview,
} from "../hooks";
import type { BlastMessage, BlastMessageLog } from "../types";

// ─── Helpers ──────────────────────────────────────────────────────────────────

function formatDate(s?: string) {
  if (!s) return "—";
  return new Date(s).toLocaleString("id-ID", { dateStyle: "medium", timeStyle: "short" });
}

function BlastStatusBadge({ status }: { status: BlastMessage["status"] }) {
  const map = {
    draft: { label: "Draft", cls: "bg-slate-100 text-slate-500", icon: Clock },
    sending: { label: "Mengirim...", cls: "bg-blue-100 text-blue-700", icon: Loader2 },
    done: { label: "Selesai", cls: "bg-green-100 text-green-700", icon: CheckCircle2 },
    failed: { label: "Gagal", cls: "bg-red-100 text-red-600", icon: XCircle },
  };
  const { label, cls, icon: Icon } = map[status] ?? map.draft;
  return (
    <span className={`inline-flex items-center gap-1 text-xs font-medium px-2 py-0.5 rounded-full ${cls}`}>
      <Icon className={`w-3 h-3 ${status === "sending" ? "animate-spin" : ""}`} />
      {label}
    </span>
  );
}

function LogStatusBadge({ status }: { status: BlastMessageLog["status"] }) {
  const map = {
    pending: { label: "Menunggu", cls: "text-slate-400" },
    sent: { label: "Terkirim", cls: "text-green-600" },
    failed: { label: "Gagal", cls: "text-red-500" },
    read: { label: "Dibaca", cls: "text-blue-600" },
  };
  const { label, cls } = map[status] ?? map.pending;
  return <span className={`text-xs font-medium ${cls}`}>{label}</span>;
}

const TARGET_TYPES = [
  {
    value: "all_tenants",
    label: "Semua Tenant Aktif",
    desc: "Semua tenant yang sedang aktif dan memiliki nomor WA",
    icon: Users,
    color: "border-blue-200 bg-blue-50",
    iconColor: "text-blue-600",
  },
  {
    value: "overdue",
    label: "Tenant Tunggakan",
    desc: "Tenant yang memiliki tagihan belum terbayar",
    icon: AlertCircle,
    color: "border-amber-200 bg-amber-50",
    iconColor: "text-amber-600",
  },
  {
    value: "manual",
    label: "Manual (Semua Aktif)",
    desc: "Sama dengan semua tenant aktif, pilih sendiri",
    icon: MessageSquare,
    color: "border-slate-200 bg-slate-50",
    iconColor: "text-slate-500",
  },
];

const VARIABLE_HINTS = [
  "{{tenant_name}}",
  "{{room_name}}",
  "{{invoice_amount}}",
  "{{rent_amount}}",
  "{{due_date}}",
  "{{days_left}}",
];

// ─── Recipient Preview Panel ──────────────────────────────────────────────────

function RecipientPreviewPanel({ targetType }: { targetType: string }) {
  const { data, isLoading } = useRecipientPreview(targetType);

  if (isLoading) return <div className="text-xs text-slate-400">Memuat penerima...</div>;
  if (!data || data.total === 0) {
    return (
      <div className="text-xs text-slate-400 italic">
        Tidak ada tenant yang memenuhi kriteria.
      </div>
    );
  }

  return (
    <div className="space-y-1.5">
      <p className="text-sm font-medium text-slate-700">
        {data.total} penerima akan menerima pesan:
      </p>
      <div className="max-h-40 overflow-y-auto space-y-1">
        {data.recipients.map((r) => (
          <div
            key={r.id}
            className="flex items-center justify-between text-xs bg-white px-3 py-1.5 rounded-lg border border-slate-100"
          >
            <div className="flex items-center gap-2">
              <div className="w-5 h-5 rounded-full bg-slate-200 flex items-center justify-center text-[10px] font-bold text-slate-500">
                {r.name.charAt(0)}
              </div>
              <span className="font-medium">{r.name}</span>
              {r.room_name && <span className="text-slate-400">· {r.room_name}</span>}
            </div>
            <div className="flex items-center gap-1 text-slate-400">
              <Phone className="w-3 h-3" />
              {r.phone}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

// ─── Blast Detail Modal ───────────────────────────────────────────────────────

function BlastDetailModal({ blastId, onClose }: { blastId: string; onClose: () => void }) {
  const { data: blast } = useBlasts();
  const blastData = blast?.find((b) => b.id === blastId);
  const { data: logs = [], isLoading } = useBlastLogs(blastId);
  const { mutate: sendBlast, isPending: isSending } = useSendBlast();

  return (
    <div className="fixed inset-0 z-50 bg-black/50 flex items-center justify-center p-4">
      <div className="bg-white rounded-2xl w-full max-w-2xl shadow-2xl max-h-[90vh] overflow-hidden flex flex-col">
        <div className="flex items-center justify-between p-6 border-b border-slate-100 shrink-0">
          <div>
            <h2 className="text-lg font-bold">{blastData?.title || "Detail Blast"}</h2>
            {blastData && <BlastStatusBadge status={blastData.status} />}
          </div>
          <button onClick={onClose} className="p-2 hover:bg-slate-100 rounded-lg">
            <X className="w-4 h-4" />
          </button>
        </div>

        <div className="overflow-y-auto flex-1 p-6 space-y-4">
          {blastData && (
            <div className="grid grid-cols-3 gap-3">
              <div className="bg-slate-50 rounded-xl p-3 text-center">
                <p className="text-2xl font-bold text-slate-800">{blastData.total_recipients}</p>
                <p className="text-xs text-slate-500 mt-1">Total Penerima</p>
              </div>
              <div className="bg-green-50 rounded-xl p-3 text-center">
                <p className="text-2xl font-bold text-green-700">{blastData.sent_count}</p>
                <p className="text-xs text-slate-500 mt-1">Berhasil</p>
              </div>
              <div className="bg-red-50 rounded-xl p-3 text-center">
                <p className="text-2xl font-bold text-red-600">{blastData.failed_count}</p>
                <p className="text-xs text-slate-500 mt-1">Gagal</p>
              </div>
            </div>
          )}

          {blastData?.status === "draft" && (
            <button
              onClick={() => sendBlast(blastId)}
              disabled={isSending}
              className="w-full flex items-center justify-center gap-2 bg-green-600 text-white py-3 rounded-xl font-medium hover:bg-green-700 transition-colors disabled:opacity-50"
            >
              {isSending ? (
                <Loader2 className="w-4 h-4 animate-spin" />
              ) : (
                <Send className="w-4 h-4" />
              )}
              {isSending ? "Mengirim..." : "Kirim Sekarang"}
            </button>
          )}

          <div>
            <h3 className="font-semibold text-sm mb-2 text-slate-700">Log Pengiriman</h3>
            {isLoading ? (
              <div className="text-center py-8 text-slate-400 text-sm">Memuat log...</div>
            ) : logs.length === 0 ? (
              <div className="text-center py-8 text-slate-400 text-sm">
                Belum ada log. Klik "Kirim Sekarang" untuk mulai.
              </div>
            ) : (
              <div className="space-y-2">
                {logs.map((log) => (
                  <div
                    key={log.id}
                    className="flex items-start justify-between gap-3 p-3 rounded-lg border border-slate-100 text-sm"
                  >
                    <div className="flex-1">
                      <div className="flex items-center gap-2">
                        <span className="font-medium">{log.recipient_name}</span>
                        <span className="text-slate-400">·</span>
                        <span className="text-slate-500 text-xs flex items-center gap-1">
                          <Phone className="w-3 h-3" />
                          {log.phone}
                        </span>
                      </div>
                      <p className="text-xs text-slate-400 mt-0.5 line-clamp-2">{log.message}</p>
                      {log.status === "failed" && log.error_message && (
                        <p className="text-xs text-red-500 mt-1 font-medium">
                          ✗ {log.error_message}
                        </p>
                      )}
                    </div>
                    <div className="text-right shrink-0">
                      <LogStatusBadge status={log.status} />
                      {log.sent_at && (
                        <p className="text-xs text-slate-400 mt-0.5">
                          {new Date(log.sent_at).toLocaleTimeString("id-ID", {
                            hour: "2-digit",
                            minute: "2-digit",
                          })}
                        </p>
                      )}
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

// ─── Main Page ────────────────────────────────────────────────────────────────

export function BlastMessagePage() {
  const { data: blasts = [], isLoading: loadingBlasts } = useBlasts();
  const { data: templates = [] } = useTemplates();
  const { data: devices = [] } = useDevices();
  const { mutate: createBlast, isPending: creating } = useCreateBlast();
  const { mutate: sendBlast, isPending: sending } = useSendBlast();

  // Form state
  const [title, setTitle] = useState("");
  const [selectedTemplate, setSelectedTemplate] = useState("");
  const [message, setMessage] = useState("");
  const [targetType, setTargetType] = useState("all_tenants");
  const [selectedDevice, setSelectedDevice] = useState("");
  const [previewName] = useState("Budi Santoso");
  const [detailBlastId, setDetailBlastId] = useState<string | null>(null);
  const [view, setView] = useState<"compose" | "history">("compose");

  const connectedDevices = devices.filter((d) => d.status === "connected");

  // Auto-select first connected device
  useEffect(() => {
    if (!selectedDevice && connectedDevices.length > 0) {
      setSelectedDevice(connectedDevices[0].id);
    }
  }, [connectedDevices, selectedDevice]);

  // When template selected, fill in message
  useEffect(() => {
    if (selectedTemplate) {
      const tmpl = templates.find((t) => t.id === selectedTemplate);
      if (tmpl) setMessage(tmpl.content);
    }
  }, [selectedTemplate, templates]);

  const insertVar = (v: string) => setMessage((prev) => prev + v);

  const previewMessage = message
    .replace("{{tenant_name}}", previewName)
    .replace("{{room_name}}", "Kamar 101")
    .replace("{{invoice_amount}}", "2.500.000")
    .replace("{{rent_amount}}", "2.500.000")
    .replace("{{due_date}}", "20 Sep 2026")
    .replace("{{days_left}}", "7")
    .replace("{{start_date}}", "01 Sep 2026")
    .replace("{{message}}", "Harap memperhatikan kebersihan area bersama.");

  const handleCreateAndSend = () => {
    if (!message.trim()) return;
    if (connectedDevices.length === 0) {
      alert("Belum ada perangkat WhatsApp yang terhubung. Hubungkan perangkat terlebih dahulu di menu Pengaturan Gateway.");
      return;
    }

    const devId = selectedDevice || connectedDevices[0]?.id;

    createBlast(
      {
        device_id: devId,
        title: title || "Blast " + new Date().toLocaleDateString("id-ID"),
        template: message,
        target_type: targetType,
      },
      {
        onSuccess: (blast) => {
          sendBlast(blast.id, {
            onSuccess: () => {
              setView("history");
              setTitle("");
              setMessage("");
              setSelectedTemplate("");
            },
            onError: (err: any) => {
              alert(err?.response?.data?.message || err?.message || "Gagal mengirim blast");
            },
          });
        },
        onError: (err: any) => {
          alert(err?.response?.data?.message || err?.message || "Gagal membuat blast");
        },
      }
    );
  };

  const handleSaveDraft = () => {
    if (!message.trim()) return;
    createBlast(
      {
        device_id: selectedDevice || undefined,
        title: title || "Draft " + new Date().toLocaleDateString("id-ID"),
        template: message,
        target_type: targetType,
      },
      {
        onSuccess: () => {
          setView("history");
          setTitle("");
          setMessage("");
          setSelectedTemplate("");
        },
      }
    );
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Blast Message (WA)</h1>
          <p className="text-muted-foreground text-sm">
            Kirim pesan WhatsApp ke banyak tenant sekaligus.
          </p>
        </div>
        <div className="flex gap-2">
          <button
            onClick={() => setView("compose")}
            className={`flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
              view === "compose"
                ? "bg-orange text-slate-900"
                : "bg-white border border-slate-200 text-slate-600 hover:bg-slate-50"
            }`}
          >
            <MessageSquare className="w-4 h-4" />
            Compose
          </button>
          <button
            onClick={() => setView("history")}
            className={`flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
              view === "history"
                ? "bg-orange text-slate-900"
                : "bg-white border border-slate-200 text-slate-600 hover:bg-slate-50"
            }`}
          >
            <History className="w-4 h-4" />
            Riwayat
            {blasts.length > 0 && (
              <span className="ml-1 text-xs bg-slate-200 text-slate-600 px-1.5 py-0.5 rounded-full">
                {blasts.length}
              </span>
            )}
          </button>
        </div>
      </div>

      {/* Compose View */}
      {view === "compose" && (
        <div className="grid grid-cols-1 xl:grid-cols-3 gap-6">
          {/* Form */}
          <div className="xl:col-span-2 space-y-4">
            {/* Title */}
            <div className="bg-white rounded-xl border border-slate-200 p-5 shadow-sm">
              <label className="block text-sm font-semibold mb-2">Judul Blast</label>
              <input
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="Contoh: Pengingat Tagihan September"
                className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:ring-2 focus:ring-orange-300 outline-none"
              />
            </div>

            {/* Device selection */}
            <div className="bg-white rounded-xl border border-slate-200 p-5 shadow-sm">
              <label className="block text-sm font-semibold mb-2">
                Perangkat Pengirim
              </label>
              {connectedDevices.length === 0 ? (
                <div className="flex items-center gap-2 p-3 bg-amber-50 rounded-lg text-xs text-amber-700">
                  <AlertCircle className="w-4 h-4 shrink-0" />
                  <span>
                    Belum ada perangkat WA yang terhubung. Hubungkan terlebih dahulu di{" "}
                    <strong>Pengaturan Gateway</strong>.
                  </span>
                </div>
              ) : (
                <select
                  value={selectedDevice}
                  onChange={(e) => setSelectedDevice(e.target.value)}
                  className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:ring-2 focus:ring-orange-300 outline-none"
                >
                  <option value="">— Pilih Perangkat —</option>
                  {connectedDevices.map((d) => (
                    <option key={d.id} value={d.id}>
                      {d.label} {d.phone ? `(${d.phone})` : ""}
                    </option>
                  ))}
                </select>
              )}
            </div>

            {/* Target Type */}
            <div className="bg-white rounded-xl border border-slate-200 p-5 shadow-sm space-y-3">
              <label className="block text-sm font-semibold">Target Penerima</label>
              <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
                {TARGET_TYPES.map((tt) => (
                  <button
                    key={tt.value}
                    type="button"
                    onClick={() => setTargetType(tt.value)}
                    className={`p-3 rounded-xl border-2 text-left transition-all ${
                      targetType === tt.value
                        ? tt.color + " border-current"
                        : "border-slate-200 hover:border-slate-300 bg-white"
                    }`}
                  >
                    <tt.icon className={`w-5 h-5 mb-2 ${targetType === tt.value ? tt.iconColor : "text-slate-400"}`} />
                    <p className="text-sm font-medium text-slate-800">{tt.label}</p>
                    <p className="text-xs text-slate-400 mt-0.5">{tt.desc}</p>
                  </button>
                ))}
              </div>

              {/* Live recipient preview */}
              <div className="mt-3 p-3 bg-slate-50 rounded-xl">
                <RecipientPreviewPanel targetType={targetType} />
              </div>
            </div>

            {/* Message Composer */}
            <div className="bg-white rounded-xl border border-slate-200 p-5 shadow-sm space-y-3">
              <div className="flex items-center justify-between">
                <label className="block text-sm font-semibold">Pesan</label>
                {templates.length > 0 && (
                  <select
                    value={selectedTemplate}
                    onChange={(e) => setSelectedTemplate(e.target.value)}
                    className="text-xs px-2 py-1 border border-slate-200 rounded-lg focus:ring-2 focus:ring-orange-300 outline-none"
                  >
                    <option value="">Gunakan Template...</option>
                    {templates.map((t) => (
                      <option key={t.id} value={t.id}>
                        {t.name}
                      </option>
                    ))}
                  </select>
                )}
              </div>

              {/* Variable pills */}
              <div className="flex flex-wrap gap-1">
                {VARIABLE_HINTS.map((v) => (
                  <button
                    key={v}
                    type="button"
                    onClick={() => insertVar(v)}
                    className="text-xs px-2 py-0.5 bg-slate-100 hover:bg-orange-100 text-slate-500 hover:text-orange-700 rounded font-mono transition-colors"
                  >
                    {v}
                  </button>
                ))}
              </div>

              <textarea
                rows={7}
                value={message}
                onChange={(e) => setMessage(e.target.value)}
                placeholder="Halo {{tenant_name}}, ini adalah pemberitahuan dari manajemen..."
                className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:ring-2 focus:ring-orange-300 outline-none resize-none"
              />

              <div className="flex justify-end gap-3 pt-2">
                <button
                  onClick={handleSaveDraft}
                  disabled={creating || !message.trim()}
                  className="flex items-center gap-2 px-4 py-2 border border-slate-200 rounded-lg text-sm hover:bg-slate-50 disabled:opacity-50"
                >
                  Simpan Draft
                </button>
                <button
                  onClick={handleCreateAndSend}
                  disabled={creating || sending || !message.trim() || connectedDevices.length === 0}
                  className="flex items-center gap-2 bg-green-600 text-white px-5 py-2 rounded-lg text-sm font-medium hover:bg-green-700 disabled:opacity-50 transition-colors"
                  title={connectedDevices.length === 0 ? "Hubungkan perangkat WhatsApp di Pengaturan Gateway terlebih dahulu" : undefined}
                >
                  {creating || sending ? (
                    <Loader2 className="w-4 h-4 animate-spin" />
                  ) : (
                    <Send className="w-4 h-4" />
                  )}
                  {creating || sending ? "Mengirim..." : "Kirim Sekarang"}
                </button>
              </div>
              {connectedDevices.length === 0 && (
                <p className="text-xs text-amber-600 text-right">
                  Perangkat WhatsApp belum terhubung. Silakan buka Pengaturan Gateway.
                </p>
              )}
            </div>
          </div>

          {/* Preview Panel */}
          <div className="space-y-4">
            <div className="bg-[#0a1628] rounded-2xl p-4 shadow-lg">
              {/* Phone notch */}
              <div className="flex justify-center mb-3">
                <div className="w-20 h-1.5 bg-slate-600 rounded-full" />
              </div>
              {/* WA Header */}
              <div className="bg-[#1F2C34] rounded-t-xl px-4 py-3 flex items-center gap-3">
                <div className="w-9 h-9 rounded-full bg-slate-500 flex items-center justify-center text-white text-sm font-bold">
                  {previewName.charAt(0)}
                </div>
                <div>
                  <p className="text-white text-sm font-medium">{previewName}</p>
                  <p className="text-green-400 text-xs">online</p>
                </div>
              </div>
              {/* Chat area */}
              <div className="bg-[#EFEAE2] min-h-48 rounded-b-xl p-4">
                <div
                  className="bg-white rounded-lg rounded-tl-none p-3 shadow-sm text-sm text-slate-800 max-w-[85%] whitespace-pre-wrap"
                  style={{ wordBreak: "break-word" }}
                >
                  {previewMessage || (
                    <span className="text-slate-400 italic">Preview pesan akan tampil di sini...</span>
                  )}
                  <span className="block text-[10px] text-slate-400 text-right mt-2">
                    {new Date().toLocaleTimeString("id-ID", {
                      hour: "2-digit",
                      minute: "2-digit",
                    })}
                  </span>
                </div>
              </div>
            </div>

            <div className="bg-white rounded-xl border border-slate-200 p-4 shadow-sm">
              <p className="text-sm font-semibold text-slate-700 mb-2">Preview Variable</p>
              <table className="w-full text-xs">
                <tbody className="divide-y divide-slate-100">
                  {[
                    ["tenant_name", previewName],
                    ["room_name", "Kamar 101"],
                    ["invoice_amount", "2.500.000"],
                    ["due_date", "20 Sep 2026"],
                    ["days_left", "7"],
                  ].map(([key, val]) => (
                    <tr key={key}>
                      <td className="py-1 font-mono text-slate-500">{`{{${key}}}`}</td>
                      <td className="py-1 text-right text-slate-800 font-medium">{val}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      )}

      {/* History View */}
      {view === "history" && (
        <div className="space-y-4">
          {loadingBlasts ? (
            <div className="text-center py-12 text-slate-400">Memuat riwayat...</div>
          ) : blasts.length === 0 ? (
            <div className="text-center py-16 bg-white rounded-xl border border-slate-200">
              <History className="w-12 h-12 text-slate-200 mx-auto mb-3" />
              <p className="font-medium text-slate-500">Belum ada riwayat blast</p>
              <p className="text-sm text-slate-400 mt-1">
                Buat dan kirim blast pertamamu di tab Compose.
              </p>
              <button
                onClick={() => setView("compose")}
                className="mt-4 flex items-center gap-2 mx-auto bg-orange text-slate-900 px-4 py-2 rounded-lg text-sm font-medium hover:bg-orange/90"
              >
                <MessageSquare className="w-4 h-4" /> Buat Blast
              </button>
            </div>
          ) : (
            <div className="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden">
              <table className="w-full text-sm">
                <thead className="bg-slate-50 border-b border-slate-100">
                  <tr>
                    {["Judul", "Target", "Status", "Penerima", "Waktu", ""].map((h) => (
                      <th key={h} className="text-left text-xs font-semibold text-slate-500 px-4 py-3">
                        {h}
                      </th>
                    ))}
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-50">
                  {blasts.map((b) => (
                    <tr key={b.id} className="hover:bg-slate-50 transition-colors">
                      <td className="px-4 py-3 font-medium text-slate-800">
                        {b.title || "—"}
                      </td>
                      <td className="px-4 py-3 text-slate-500">
                        {TARGET_TYPES.find((t) => t.value === b.target_type)?.label ?? b.target_type}
                      </td>
                      <td className="px-4 py-3">
                        <BlastStatusBadge status={b.status} />
                      </td>
                      <td className="px-4 py-3">
                        {b.status === "done" || b.status === "sending" ? (
                          <div className="flex items-center gap-2 text-xs">
                            <span className="text-green-600 font-medium">{b.sent_count} ✓</span>
                            {b.failed_count > 0 && (
                              <span className="text-red-500">{b.failed_count} ✗</span>
                            )}
                            <span className="text-slate-400">/ {b.total_recipients}</span>
                          </div>
                        ) : (
                          <span className="text-slate-400">—</span>
                        )}
                      </td>
                      <td className="px-4 py-3 text-slate-400 text-xs">
                        {formatDate(b.created_at)}
                      </td>
                      <td className="px-4 py-3">
                        <button
                          onClick={() => setDetailBlastId(b.id)}
                          className="flex items-center gap-1 text-xs text-slate-500 hover:text-slate-800 transition-colors"
                        >
                          <Eye className="w-3.5 h-3.5" />
                          Detail
                          <ChevronRight className="w-3 h-3" />
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      )}

      {/* Blast Detail Modal */}
      {detailBlastId && (
        <BlastDetailModal blastId={detailBlastId} onClose={() => setDetailBlastId(null)} />
      )}
    </div>
  );
}

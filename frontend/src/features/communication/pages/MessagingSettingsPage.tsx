import { useState, useEffect } from "react";
import QRCode from "react-qr-code";
import { useQueryClient } from "@tanstack/react-query";
import {
  Smartphone,
  Plus,
  Trash2,
  Wifi,
  WifiOff,
  QrCode,
  RefreshCw,
  BookOpen,
  Pencil,
  X,
  Check,
  AlertCircle,
} from "lucide-react";
import {
  useDevices,
  useCreateDevice,
  useUpdateDeviceStatus,
  useDeleteDevice,
  useDeviceQR,
  useTemplates,
  useCreateTemplate,
  useUpdateTemplate,
  useDeleteTemplate,
} from "../hooks";
import type { WADevice, MessageTemplate } from "../types";

// ─── Helpers ──────────────────────────────────────────────────────────────────

function StatusBadge({ status }: { status: WADevice["status"] }) {
  const map = {
    connected: { label: "Terhubung", cls: "bg-green-100 text-green-700" },
    disconnected: { label: "Terputus", cls: "bg-slate-100 text-slate-500" },
    qr_pending: { label: "Menunggu QR", cls: "bg-amber-100 text-amber-700" },
  };
  const { label, cls } = map[status] ?? map.disconnected;
  return (
    <span className={`text-xs font-medium px-2 py-0.5 rounded-full ${cls}`}>
      {status === "connected" && (
        <span className="inline-block w-1.5 h-1.5 rounded-full bg-green-500 mr-1 animate-pulse" />
      )}
      {label}
    </span>
  );
}

function CategoryBadge({ cat }: { cat: string }) {
  const map: Record<string, string> = {
    invoice: "bg-blue-100 text-blue-700",
    reminder: "bg-amber-100 text-amber-700",
    announcement: "bg-purple-100 text-purple-700",
    general: "bg-slate-100 text-slate-600",
  };
  const labels: Record<string, string> = {
    invoice: "Tagihan",
    reminder: "Pengingat",
    announcement: "Pengumuman",
    general: "Umum",
  };
  return (
    <span className={`text-xs px-2 py-0.5 rounded-full font-medium ${map[cat] ?? map.general}`}>
      {labels[cat] ?? cat}
    </span>
  );
}

// ─── Device QR Modal ─────────────────────────────────────────────────────────

function QRModal({ device, onClose }: { device: WADevice; onClose: () => void }) {
  const queryClient = useQueryClient();
  const { data: devices = [] } = useDevices();
  const currentDevice = devices.find((d) => d.id === device.id) || device;
  const isConnected = currentDevice.status === "connected";

  const { mutate: getQR, data: qrData, isPending, isError, error } = useDeviceQR();
  const { mutate: updateStatus } = useUpdateDeviceStatus();
  const [countdown, setCountdown] = useState(0);

  // Auto-fetch QR when modal opens
  useEffect(() => {
    if (!isConnected) {
      getQR(device.id);
    }
  }, [device.id, isConnected]);

  // Poll device status every 2 seconds to detect phone scan completion
  useEffect(() => {
    if (isConnected) return;
    const interval = setInterval(() => {
      queryClient.invalidateQueries({ queryKey: ["wa-devices"] });
    }, 2000);
    return () => clearInterval(interval);
  }, [isConnected, queryClient]);

  // Countdown timer when QR is active
  useEffect(() => {
    if (qrData?.expires_in) {
      setCountdown(qrData.expires_in);
      const interval = setInterval(() => {
        setCountdown((prev) => {
          if (prev <= 1) {
            clearInterval(interval);
            return 0;
          }
          return prev - 1;
        });
      }, 1000);
      return () => clearInterval(interval);
    }
  }, [qrData]);

  const handleRefresh = () => {
    getQR(device.id);
  };

  const handleManualConnect = () => {
    updateStatus({ id: device.id, status: "connected" });
  };

  const isExpired = countdown === 0 && !!qrData?.qr_string;

  // Render Connected Success Screen
  if (isConnected) {
    return (
      <div className="fixed inset-0 z-50 bg-black/50 flex items-center justify-center p-4">
        <div className="bg-white rounded-2xl w-full max-w-md shadow-2xl p-6 text-center space-y-4">
          <div className="w-16 h-16 rounded-full bg-green-100 text-green-600 flex items-center justify-center mx-auto animate-bounce">
            <Check className="w-8 h-8" />
          </div>
          <div>
            <h2 className="text-xl font-bold text-slate-800">WhatsApp Berhasil Terhubung!</h2>
            <p className="text-sm text-slate-500 mt-1">{currentDevice.label}</p>
            {currentDevice.phone ? (
              <div className="mt-3 p-3 bg-green-50 rounded-xl border border-green-200">
                <p className="text-xs text-green-600">Nomor WhatsApp Terhubung:</p>
                <p className="text-lg font-bold text-green-800 font-mono mt-0.5">
                  +{currentDevice.phone}
                </p>
              </div>
            ) : (
              <p className="text-xs text-green-600 mt-2">Perangkat aktif dan siap digunakan.</p>
            )}
          </div>
          <button
            onClick={onClose}
            className="w-full bg-slate-900 text-white py-2.5 rounded-lg text-sm font-medium hover:bg-slate-800 transition-colors"
          >
            Selesai
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="fixed inset-0 z-50 bg-black/50 flex items-center justify-center p-4">
      <div className="bg-white rounded-2xl w-full max-w-md shadow-2xl">
        <div className="flex items-center justify-between p-6 border-b border-slate-100">
          <div>
            <h2 className="text-lg font-bold">Hubungkan WhatsApp</h2>
            <p className="text-sm text-slate-500">{device.label}</p>
          </div>
          <button onClick={onClose} className="p-2 hover:bg-slate-100 rounded-lg transition-colors">
            <X className="w-4 h-4" />
          </button>
        </div>

        <div className="p-6 text-center space-y-4">
          <p className="text-sm text-slate-600">
            Buka WhatsApp di HP → <strong>Pengaturan</strong> atau <strong>Menu (titik tiga)</strong> →{" "}
            <strong>Perangkat Tertaut</strong> → <strong>Tautkan Perangkat</strong>
          </p>

          {/* QR Code display */}
          <div className="relative flex items-center justify-center">
            <div
              className={`p-4 bg-white rounded-2xl border-2 ${
                isExpired ? "border-red-200" : isError ? "border-amber-200" : "border-slate-200"
              } shadow-sm transition-all`}
            >
              {isPending ? (
                <div className="w-52 h-52 flex flex-col items-center justify-center gap-3">
                  <RefreshCw className="w-8 h-8 text-slate-400 animate-spin" />
                  <p className="text-xs text-slate-400">Menghubungkan ke WhatsApp...</p>
                </div>
              ) : isError ? (
                <div className="w-52 h-52 flex flex-col items-center justify-center p-3 gap-2 text-center">
                  <AlertCircle className="w-8 h-8 text-amber-500" />
                  <p className="text-xs font-semibold text-slate-700">Gagal Memuat QR WhatsApp</p>
                  <p className="text-[11px] text-slate-400">
                    {(error as Error)?.message || "Koneksi ke WhatsApp terputus. Silakan coba lagi."}
                  </p>
                  <button
                    onClick={handleRefresh}
                    className="mt-1 px-3 py-1 bg-slate-800 text-white rounded-md text-xs hover:bg-slate-700"
                  >
                    Coba Lagi
                  </button>
                </div>
              ) : qrData?.qr_string ? (
                <div className={`transition-opacity ${isExpired ? "opacity-25" : "opacity-100"}`}>
                  <QRCode
                    value={qrData.qr_string}
                    size={208}
                    level="M"
                    style={{ height: "auto", maxWidth: "100%", width: "100%" }}
                  />
                </div>
              ) : (
                <div className="w-52 h-52 flex flex-col items-center justify-center gap-3">
                  <QrCode className="w-12 h-12 text-slate-300" />
                  <p className="text-xs text-slate-400">Memuat QR code...</p>
                </div>
              )}
            </div>

            {/* Expired overlay */}
            {isExpired && (
              <div className="absolute inset-0 flex flex-col items-center justify-center gap-3 bg-white/75 backdrop-blur-[1px] rounded-2xl">
                <button
                  onClick={handleRefresh}
                  disabled={isPending}
                  className="flex items-center gap-2 bg-slate-900 text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-slate-800 transition-colors shadow-lg"
                >
                  <RefreshCw className={`w-4 h-4 ${isPending ? "animate-spin" : ""}`} />
                  Perbarui QR
                </button>
                <p className="text-xs text-slate-500 font-medium">QR code kadaluarsa</p>
              </div>
            )}
          </div>

          {/* Countdown & Status indicator */}
          {qrData?.qr_string && !isExpired && (
            <div className="space-y-1">
              <p className="text-xs text-slate-500">
                Arahkan kamera ke QR code di atas
              </p>
              <p className="text-xs text-slate-400">
                Kadaluarsa dalam:{" "}
                <span className={`font-mono font-bold ${countdown <= 15 ? "text-red-500" : "text-slate-700"}`}>
                  {countdown} detik
                </span>
              </p>
            </div>
          )}

          <div className="flex gap-3 pt-1">
            <button
              onClick={handleRefresh}
              disabled={isPending}
              className="flex-1 flex items-center justify-center gap-2 border border-slate-200 text-slate-600 px-4 py-2.5 rounded-lg font-medium hover:bg-slate-50 transition-colors disabled:opacity-50 text-sm"
            >
              <RefreshCw className={`w-4 h-4 ${isPending ? "animate-spin" : ""}`} />
              Perbarui QR
            </button>
            <button
              onClick={handleManualConnect}
              className="px-4 py-2.5 border border-slate-200 text-slate-400 hover:text-slate-600 rounded-lg text-xs transition-colors"
              title="Gunakan jika ingin menandai terhubung manual"
            >
              Tandai Manual
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

// ─── Add Device Modal ─────────────────────────────────────────────────────────

function AddDeviceModal({ onClose }: { onClose: () => void }) {
  const { mutate: createDevice, isPending } = useCreateDevice();
  const [label, setLabel] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    createDevice({ label }, { onSuccess: onClose });
  };

  return (
    <div className="fixed inset-0 z-50 bg-black/50 flex items-center justify-center p-4">
      <div className="bg-white rounded-2xl w-full max-w-sm shadow-2xl">
        <div className="flex items-center justify-between p-6 border-b border-slate-100">
          <div>
            <h2 className="text-lg font-bold">Tambah Perangkat WA</h2>
            <p className="text-sm text-slate-500">Nomor akan terdeteksi setelah scan QR</p>
          </div>
          <button onClick={onClose} className="p-2 hover:bg-slate-100 rounded-lg">
            <X className="w-4 h-4" />
          </button>
        </div>
        <form onSubmit={handleSubmit} className="p-6 space-y-4">
          <div>
            <label className="block text-sm font-medium mb-1">Label Perangkat *</label>
            <input
              value={label}
              onChange={(e) => setLabel(e.target.value)}
              placeholder="Contoh: Divisi Keuangan, Operasional..."
              className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:ring-2 focus:ring-orange-300 outline-none"
              required
              autoFocus
            />
            <p className="text-xs text-slate-400 mt-1">
              Nama untuk mengidentifikasi perangkat ini di dashboard.
            </p>
          </div>
          <div className="flex gap-3 pt-2">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 px-4 py-2 border border-slate-200 rounded-lg text-sm hover:bg-slate-50"
            >
              Batal
            </button>
            <button
              type="submit"
              disabled={isPending || !label.trim()}
              className="flex-1 px-4 py-2 bg-orange text-slate-900 rounded-lg text-sm font-medium hover:bg-orange/90 disabled:opacity-50"
            >
              {isPending ? "Menyimpan..." : "Tambah"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

// ─── Template Modal ───────────────────────────────────────────────────────────

const VARIABLE_HINTS = [
  "{{tenant_name}}",
  "{{room_name}}",
  "{{invoice_amount}}",
  "{{rent_amount}}",
  "{{due_date}}",
  "{{days_left}}",
  "{{start_date}}",
  "{{message}}",
];

function TemplateModal({
  template,
  onClose,
}: {
  template?: MessageTemplate;
  onClose: () => void;
}) {
  const { mutate: create, isPending: creating } = useCreateTemplate();
  const { mutate: update, isPending: updating } = useUpdateTemplate();
  const [name, setName] = useState(template?.name ?? "");
  const [content, setContent] = useState(template?.content ?? "");
  const [category, setCategory] = useState<string>(template?.category ?? "general");

  const isPending = creating || updating;

  const extractVars = (text: string): string[] => {
    const matches = text.match(/\{\{([^}]+)\}\}/g) ?? [];
    return [...new Set(matches.map((m) => m.slice(2, -2)))];
  };

  const insertVar = (v: string) => setContent((prev) => prev + v);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const variables = extractVars(content);
    if (template) {
      update(
        {
          id: template.id,
          data: { name, content, category, variables, is_active: template.is_active },
        },
        { onSuccess: onClose }
      );
    } else {
      create({ name, content, category, variables }, { onSuccess: onClose });
    }
  };

  return (
    <div className="fixed inset-0 z-50 bg-black/50 flex items-center justify-center p-4">
      <div className="bg-white rounded-2xl w-full max-w-lg shadow-2xl max-h-[90vh] overflow-y-auto">
        <div className="flex items-center justify-between p-6 border-b border-slate-100 sticky top-0 bg-white">
          <h2 className="text-lg font-bold">
            {template ? "Edit Template" : "Buat Template"}
          </h2>
          <button onClick={onClose} className="p-2 hover:bg-slate-100 rounded-lg">
            <X className="w-4 h-4" />
          </button>
        </div>
        <form onSubmit={handleSubmit} className="p-6 space-y-4">
          <div>
            <label className="block text-sm font-medium mb-1">Nama Template *</label>
            <input
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Contoh: Tagihan Bulanan"
              className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:ring-2 focus:ring-orange-300 outline-none"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">Kategori</label>
            <select
              value={category}
              onChange={(e) => setCategory(e.target.value)}
              className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:ring-2 focus:ring-orange-300 outline-none"
            >
              <option value="general">Umum</option>
              <option value="invoice">Tagihan</option>
              <option value="reminder">Pengingat</option>
              <option value="announcement">Pengumuman</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">Konten Pesan *</label>
            <div className="flex flex-wrap gap-1 mb-2">
              {VARIABLE_HINTS.map((v) => (
                <button
                  key={v}
                  type="button"
                  onClick={() => insertVar(v)}
                  className="text-xs px-2 py-0.5 bg-slate-100 hover:bg-orange-100 text-slate-600 hover:text-orange-700 rounded font-mono transition-colors"
                >
                  {v}
                </button>
              ))}
            </div>
            <textarea
              value={content}
              onChange={(e) => setContent(e.target.value)}
              rows={6}
              placeholder="Halo {{tenant_name}}, tagihan Anda sebesar Rp {{invoice_amount}}..."
              className="w-full px-3 py-2 border border-slate-200 rounded-lg text-sm focus:ring-2 focus:ring-orange-300 outline-none resize-none"
              required
            />
            {content && (
              <p className="text-xs text-slate-500 mt-1">
                Variabel terdeteksi: {extractVars(content).map((v) => `{{${v}}}`).join(", ") || "—"}
              </p>
            )}
          </div>
          <div className="flex gap-3 pt-2">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 px-4 py-2 border border-slate-200 rounded-lg text-sm hover:bg-slate-50"
            >
              Batal
            </button>
            <button
              type="submit"
              disabled={isPending || !name.trim() || !content.trim()}
              className="flex-1 px-4 py-2 bg-orange text-slate-900 rounded-lg text-sm font-medium hover:bg-orange/90 disabled:opacity-50"
            >
              {isPending ? "Menyimpan..." : template ? "Update" : "Buat Template"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

// ─── Main Page ────────────────────────────────────────────────────────────────

export function MessagingSettingsPage() {
  const { data: devices = [], isLoading: loadingDevices } = useDevices();
  const { data: templates = [], isLoading: loadingTemplates } = useTemplates();
  const { mutate: deleteDevice } = useDeleteDevice();
  const { mutate: deleteTemplate } = useDeleteTemplate();
  const { mutate: updateStatus } = useUpdateDeviceStatus();

  const [showAddDevice, setShowAddDevice] = useState(false);
  const [qrDevice, setQrDevice] = useState<WADevice | null>(null);
  const [editTemplate, setEditTemplate] = useState<MessageTemplate | undefined>(undefined);
  const [showAddTemplate, setShowAddTemplate] = useState(false);
  const [activeTab, setActiveTab] = useState<"devices" | "templates">("devices");

  const handleDisconnect = (device: WADevice) => {
    updateStatus({ id: device.id, status: "disconnected" });
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Pengaturan Gateway</h1>
          <p className="text-muted-foreground text-sm">
            Kelola perangkat WhatsApp dan template pesan untuk komunikasi dengan tenant.
          </p>
        </div>
        <button
          onClick={() =>
            activeTab === "devices" ? setShowAddDevice(true) : setShowAddTemplate(true)
          }
          className="flex items-center gap-2 bg-orange text-slate-900 px-4 py-2 rounded-lg font-medium hover:bg-orange/90 transition-colors"
        >
          <Plus className="w-4 h-4" />
          {activeTab === "devices" ? "Tambah Perangkat" : "Buat Template"}
        </button>
      </div>

      {/* Tabs */}
      <div className="flex gap-1 p-1 bg-slate-100 rounded-xl w-fit">
        {(["devices", "templates"] as const).map((tab) => (
          <button
            key={tab}
            onClick={() => setActiveTab(tab)}
            className={`flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
              activeTab === tab
                ? "bg-white text-slate-900 shadow-sm"
                : "text-slate-500 hover:text-slate-800"
            }`}
          >
            {tab === "devices" ? (
              <>
                <Smartphone className="w-4 h-4" /> Perangkat WA
              </>
            ) : (
              <>
                <BookOpen className="w-4 h-4" /> Template Pesan
              </>
            )}
          </button>
        ))}
      </div>

      {/* Devices Tab */}
      {activeTab === "devices" && (
        <div className="space-y-4">
          {loadingDevices ? (
            <div className="text-center py-12 text-slate-400">Memuat...</div>
          ) : devices.length === 0 ? (
            <div className="text-center py-16 bg-white rounded-xl border border-slate-200">
              <Smartphone className="w-12 h-12 text-slate-200 mx-auto mb-3" />
              <p className="font-medium text-slate-500">Belum ada perangkat terhubung</p>
              <p className="text-sm text-slate-400 mt-1">
                Tambahkan perangkat WhatsApp untuk mulai mengirim pesan.
              </p>
              <button
                onClick={() => setShowAddDevice(true)}
                className="mt-4 flex items-center gap-2 mx-auto bg-orange text-slate-900 px-4 py-2 rounded-lg text-sm font-medium hover:bg-orange/90"
              >
                <Plus className="w-4 h-4" /> Tambah Perangkat
              </button>
            </div>
          ) : (
            <div className="grid gap-4">
              {devices.map((device) => (
                <div
                  key={device.id}
                  className="bg-white rounded-xl border border-slate-200 p-4 flex items-center justify-between shadow-sm"
                >
                  <div className="flex items-center gap-4">
                    <div
                      className={`w-11 h-11 rounded-full flex items-center justify-center font-bold text-sm ${
                        device.status === "connected"
                          ? "bg-green-100 text-green-700"
                          : "bg-slate-100 text-slate-500"
                      }`}
                    >
                      WA
                    </div>
                    <div>
                      <div className="flex items-center gap-2">
                        <p className="font-semibold text-slate-800">{device.label}</p>
                        <StatusBadge status={device.status} />
                      </div>
                      {device.phone && (
                        <p className="text-sm text-slate-500 mt-0.5">{device.phone}</p>
                      )}
                      {device.last_seen && (
                        <p className="text-xs text-slate-400 mt-0.5">
                          Terakhir aktif:{" "}
                          {new Date(device.last_seen).toLocaleString("id-ID", {
                            dateStyle: "medium",
                            timeStyle: "short",
                          })}
                        </p>
                      )}
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    {device.status !== "connected" ? (
                      <button
                        onClick={() => setQrDevice(device)}
                        className="flex items-center gap-1.5 text-sm px-3 py-1.5 bg-slate-800 text-white rounded-lg hover:bg-slate-700 transition-colors"
                      >
                        <Wifi className="w-3.5 h-3.5" />
                        Hubungkan
                      </button>
                    ) : (
                      <button
                        onClick={() => handleDisconnect(device)}
                        className="flex items-center gap-1.5 text-sm px-3 py-1.5 border border-slate-200 text-slate-500 rounded-lg hover:bg-slate-50 transition-colors"
                      >
                        <WifiOff className="w-3.5 h-3.5" />
                        Putuskan
                      </button>
                    )}
                    <button
                      onClick={() => {
                        if (confirm(`Hapus perangkat "${device.label}"?`)) deleteDevice(device.id);
                      }}
                      className="p-2 text-red-400 hover:bg-red-50 rounded-lg transition-colors"
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      {/* Templates Tab */}
      {activeTab === "templates" && (
        <div className="space-y-4">
          {loadingTemplates ? (
            <div className="text-center py-12 text-slate-400">Memuat...</div>
          ) : (
            <div className="grid gap-4">
              {templates.map((t) => (
                <div
                  key={t.id}
                  className="bg-white rounded-xl border border-slate-200 p-4 shadow-sm"
                >
                  <div className="flex items-start justify-between gap-4">
                    <div className="flex-1">
                      <div className="flex items-center gap-2 mb-1">
                        <p className="font-semibold text-slate-800">{t.name}</p>
                        <CategoryBadge cat={t.category} />
                      </div>
                      <p className="text-sm text-slate-500 line-clamp-2">{t.content}</p>
                      {t.variables.length > 0 && (
                        <div className="flex flex-wrap gap-1 mt-2">
                          {t.variables.map((v) => (
                            <span
                              key={v}
                              className="text-xs px-1.5 py-0.5 bg-slate-100 text-slate-500 rounded font-mono"
                            >
                              {`{{${v}}}`}
                            </span>
                          ))}
                        </div>
                      )}
                    </div>
                    <div className="flex gap-2 shrink-0">
                      <button
                        onClick={() => setEditTemplate(t)}
                        className="p-2 text-slate-400 hover:bg-slate-100 rounded-lg transition-colors"
                      >
                        <Pencil className="w-4 h-4" />
                      </button>
                      <button
                        onClick={() => {
                          if (confirm(`Hapus template "${t.name}"?`)) deleteTemplate(t.id);
                        }}
                        className="p-2 text-red-400 hover:bg-red-50 rounded-lg transition-colors"
                      >
                        <Trash2 className="w-4 h-4" />
                      </button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      )}

      {/* Modals */}
      {showAddDevice && <AddDeviceModal onClose={() => setShowAddDevice(false)} />}
      {qrDevice && <QRModal device={qrDevice} onClose={() => setQrDevice(null)} />}
      {(showAddTemplate || editTemplate) && (
        <TemplateModal
          template={editTemplate}
          onClose={() => {
            setShowAddTemplate(false);
            setEditTemplate(undefined);
          }}
        />
      )}
    </div>
  );
}

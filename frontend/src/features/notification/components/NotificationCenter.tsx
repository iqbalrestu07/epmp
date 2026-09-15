import { useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useQueryClient } from "@tanstack/react-query";
import { Bell, CheckCheck } from "lucide-react";
import {
  useNotifications,
  useMarkNotificationRead,
  useMarkAllNotificationsRead,
  NOTIFICATIONS_KEY,
} from "../hooks";
import { useOrg } from "@/features/organization/context/OrgContext";
import { getStoredToken } from "@/services/api";

const TYPE_COLORS: Record<string, string> = {
  payment: "bg-emerald-100 text-emerald-700",
  warning: "bg-amber-100 text-amber-700",
  info: "bg-blue-100 text-blue-700",
};

export function NotificationCenter() {
  const [open, setOpen] = useState(false);
  const rootRef = useRef<HTMLDivElement>(null);
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { orgId } = useOrg();

  const { data } = useNotifications();
  const markRead = useMarkNotificationRead();
  const markAll = useMarkAllNotificationsRead();

  const unread = data?.unread_count ?? 0;

  // Real-time updates via websocket; falls back to the query's poll interval.
  useEffect(() => {
    const token = getStoredToken();
    if (!token || !orgId) return;

    let ws: WebSocket | null = null;
    let closed = false;
    let retry: ReturnType<typeof setTimeout>;

    const connect = () => {
      if (closed) return;
      const proto = window.location.protocol === "https:" ? "wss" : "ws";
      ws = new WebSocket(
        `${proto}://${window.location.host}/api/v1/ws?token=${encodeURIComponent(token)}&org_id=${encodeURIComponent(orgId)}`
      );
      ws.onmessage = () => {
        queryClient.invalidateQueries({ queryKey: [NOTIFICATIONS_KEY] });
      };
      ws.onclose = () => {
        if (!closed) retry = setTimeout(connect, 5000);
      };
    };
    connect();

    return () => {
      closed = true;
      clearTimeout(retry);
      ws?.close();
    };
  }, [orgId, queryClient]);

  // Close dropdown on outside click.
  useEffect(() => {
    if (!open) return;
    const onClick = (e: MouseEvent) => {
      if (rootRef.current && !rootRef.current.contains(e.target as Node)) {
        setOpen(false);
      }
    };
    document.addEventListener("mousedown", onClick);
    return () => document.removeEventListener("mousedown", onClick);
  }, [open]);

  const handleClick = (id: string, isRead: boolean, link?: string) => {
    if (!isRead) markRead.mutate(id);
    if (link) {
      setOpen(false);
      navigate(link);
    }
  };

  return (
    <div ref={rootRef} className="relative">
      <button
        onClick={() => setOpen(!open)}
        className="w-10 h-10 rounded-full bg-slate-100 hover:bg-slate-200 flex items-center justify-center relative transition-colors"
        aria-label="Notifications"
      >
        <Bell size={18} className="text-slate-600" />
        {unread > 0 && (
          <span className="absolute -top-0.5 -right-0.5 min-w-[18px] h-[18px] px-1 rounded-full bg-red-500 text-white text-[10px] font-bold flex items-center justify-center border border-white">
            {unread > 99 ? "99+" : unread}
          </span>
        )}
      </button>

      {open && (
        <div className="absolute right-0 mt-2 w-96 max-w-[90vw] rounded-lg border bg-white shadow-lg z-50">
          <div className="flex items-center justify-between px-4 py-3 border-b">
            <h3 className="text-sm font-semibold text-slate-900">Notifikasi</h3>
            {unread > 0 && (
              <button
                onClick={() => markAll.mutate()}
                className="text-xs text-blue-600 hover:text-blue-800 flex items-center gap-1"
              >
                <CheckCheck size={14} /> Tandai semua dibaca
              </button>
            )}
          </div>
          <div className="max-h-96 overflow-y-auto">
            {(data?.data ?? []).length === 0 ? (
              <p className="px-4 py-8 text-center text-sm text-slate-400">Tidak ada notifikasi</p>
            ) : (
              (data?.data ?? []).map((n) => (
                <button
                  key={n.id}
                  onClick={() => handleClick(n.id, n.is_read, n.link)}
                  className={`w-full text-left px-4 py-3 border-b last:border-0 hover:bg-slate-50 transition-colors ${
                    n.is_read ? "opacity-60" : ""
                  }`}
                >
                  <div className="flex items-start gap-2">
                    <span className={`mt-1 w-2 h-2 rounded-full shrink-0 ${n.is_read ? "bg-transparent" : "bg-orange-500"}`} />
                    <div className="min-w-0 flex-1">
                      <div className="flex items-center gap-2">
                        <span className={`px-1.5 py-0.5 rounded text-[10px] font-medium ${TYPE_COLORS[n.type] ?? TYPE_COLORS.info}`}>
                          {n.type}
                        </span>
                        <span className="text-[10px] text-slate-400">
                          {new Date(n.created_at).toLocaleString("id-ID")}
                        </span>
                      </div>
                      <p className="text-sm font-medium text-slate-900 mt-1">{n.title}</p>
                      <p className="text-xs text-slate-500 mt-0.5 line-clamp-2">{n.message}</p>
                    </div>
                  </div>
                </button>
              ))
            )}
          </div>
        </div>
      )}
    </div>
  );
}

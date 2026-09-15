import { api } from "@/services/api";
import type { AuditLogListResponse, AuditLogQueryParams } from "../types";

const BASE_PATH = "/audit-logs";

export async function fetchAuditLogs(params?: AuditLogQueryParams): Promise<AuditLogListResponse> {
  const query = new URLSearchParams();
  if (params?.page) query.set("page", String(params.page));
  if (params?.per_page) query.set("per_page", String(params.per_page));
  if (params?.search) query.set("search", params.search);
  if (params?.module) query.set("module", params.module);
  if (params?.action) query.set("action", params.action);
  const qs = query.toString();
  const res = await api.get<{ success: boolean; data: AuditLogListResponse }>(
    `${BASE_PATH}${qs ? `?${qs}` : ""}`
  );
  return res.data;
}

import { useQuery } from "@tanstack/react-query";
import { fetchAuditLogs } from "../api";
import type { AuditLogQueryParams } from "../types";

const QUERY_KEY = "audit-logs";

export function useAuditLogs(params?: AuditLogQueryParams) {
  return useQuery({
    queryKey: [QUERY_KEY, params],
    queryFn: () => fetchAuditLogs(params),
  });
}

export interface AuditLog {
  id: string;
  organization_id: string;
  user_id: string;
  user_email: string;
  action: "CREATE" | "UPDATE" | "DELETE" | string;
  module: string;
  entity_id: string;
  method: string;
  path: string;
  status_code: number;
  ip_address: string;
  user_agent: string;
  request_body?: Record<string, unknown> | { _raw: string };
  created_at: string;
}

export interface AuditLogListResponse {
  data: AuditLog[];
  total: number;
  page: number;
  per_page: number;
  total_pages: number;
}

export interface AuditLogQueryParams {
  page?: number;
  per_page?: number;
  search?: string;
  module?: string;
  action?: string;
}

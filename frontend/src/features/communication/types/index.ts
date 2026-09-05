// types/index.ts — Communication module types

export interface WADevice {
  id: string;
  org_id: string;
  label: string;
  phone: string;
  status: "connected" | "disconnected" | "qr_pending";
  last_seen?: string;
  created_at: string;
  updated_at: string;
}

export interface MessageTemplate {
  id: string;
  org_id: string;
  name: string;
  content: string;
  variables: string[];
  category: "general" | "invoice" | "reminder" | "announcement";
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface BlastMessage {
  id: string;
  org_id: string;
  device_id?: string;
  title: string;
  template: string;
  target_type: "all_tenants" | "overdue" | "building" | "manual";
  target_filter?: Record<string, unknown>;
  status: "draft" | "sending" | "done" | "failed";
  total_recipients: number;
  sent_count: number;
  failed_count: number;
  scheduled_at?: string;
  started_at?: string;
  completed_at?: string;
  created_at: string;
  updated_at: string;
}

export interface BlastMessageLog {
  id: string;
  blast_id: string;
  tenant_id?: string;
  phone: string;
  recipient_name: string;
  message: string;
  status: "pending" | "sent" | "failed" | "read";
  error_message?: string;
  sent_at?: string;
  read_at?: string;
  created_at: string;
}

export interface RecipientPreview {
  id: string;
  name: string;
  phone: string;
  room_name: string;
}

export interface RecipientPreviewResponse {
  recipients: RecipientPreview[];
  total: number;
}

export interface CreateDeviceRequest {
  label: string;
  phone?: string;
}

export interface CreateTemplateRequest {
  name: string;
  content: string;
  category: string;
  variables: string[];
}

export interface UpdateTemplateRequest extends CreateTemplateRequest {
  is_active: boolean;
}

export interface CreateBlastRequest {
  device_id?: string;
  title: string;
  template: string;
  target_type: string;
  target_filter?: Record<string, unknown>;
  scheduled_at?: string;
}

export type TargetType = "all_tenants" | "overdue" | "manual";

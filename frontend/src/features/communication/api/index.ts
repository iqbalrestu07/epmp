// api/index.ts — Communication module API functions

import { api } from "@/services/api";
import type {
  WADevice,
  MessageTemplate,
  BlastMessage,
  BlastMessageLog,
  RecipientPreviewResponse,
  CreateDeviceRequest,
  CreateTemplateRequest,
  UpdateTemplateRequest,
  CreateBlastRequest,
} from "../types";

const BASE = "/communication";

// ── Devices ──────────────────────────────────────────────────────────────────

export async function fetchDevices(): Promise<WADevice[]> {
  const res = await api.get<{ success: boolean; data: WADevice[] }>(`${BASE}/devices`);
  return res.data;
}

export async function createDevice(data: CreateDeviceRequest): Promise<WADevice> {
  const res = await api.post<{ success: boolean; data: WADevice }>(`${BASE}/devices`, data);
  return res.data;
}

export async function updateDeviceStatus(
  id: string,
  status: string,
  phone?: string
): Promise<WADevice> {
  const res = await api.put<{ success: boolean; data: WADevice }>(
    `${BASE}/devices/${id}/status`,
    { status, phone }
  );
  return res.data;
}

export async function getDeviceQR(id: string): Promise<{ qr_string: string; expires_in: number }> {
  const res = await api.get<{ success: boolean; data: { qr_string: string; expires_in: number } }>(
    `${BASE}/devices/${id}/qr`
  );
  return res.data;
}

export async function deleteDevice(id: string): Promise<void> {
  return api.delete<void>(`${BASE}/devices/${id}`);
}

// ── Templates ─────────────────────────────────────────────────────────────────

export async function fetchTemplates(): Promise<MessageTemplate[]> {
  const res = await api.get<{ success: boolean; data: MessageTemplate[] }>(`${BASE}/templates`);
  return res.data;
}

export async function createTemplate(data: CreateTemplateRequest): Promise<MessageTemplate> {
  const res = await api.post<{ success: boolean; data: MessageTemplate }>(`${BASE}/templates`, data);
  return res.data;
}

export async function updateTemplate(id: string, data: UpdateTemplateRequest): Promise<MessageTemplate> {
  const res = await api.put<{ success: boolean; data: MessageTemplate }>(
    `${BASE}/templates/${id}`,
    data
  );
  return res.data;
}

export async function deleteTemplate(id: string): Promise<void> {
  return api.delete<void>(`${BASE}/templates/${id}`);
}

// ── Blast Messages ────────────────────────────────────────────────────────────

export async function fetchBlasts(): Promise<BlastMessage[]> {
  const res = await api.get<{ success: boolean; data: BlastMessage[] }>(`${BASE}/blast`);
  return res.data;
}

export async function fetchBlast(id: string): Promise<BlastMessage> {
  const res = await api.get<{ success: boolean; data: BlastMessage }>(`${BASE}/blast/${id}`);
  return res.data;
}

export async function createBlast(data: CreateBlastRequest): Promise<BlastMessage> {
  const res = await api.post<{ success: boolean; data: BlastMessage }>(`${BASE}/blast`, data);
  return res.data;
}

export async function fetchBlastLogs(id: string): Promise<BlastMessageLog[]> {
  const res = await api.get<{ success: boolean; data: BlastMessageLog[] }>(
    `${BASE}/blast/${id}/logs`
  );
  return res.data;
}

export async function sendBlast(
  id: string
): Promise<{ message: string; total_recipients: number; sent_count: number; failed_count: number }> {
  const res = await api.post<{
    success: boolean;
    data: { message: string; total_recipients: number; sent_count: number; failed_count: number };
  }>(`${BASE}/blast/${id}/send`, {});
  return res.data;
}

export async function fetchRecipientPreview(
  targetType: string
): Promise<RecipientPreviewResponse> {
  const res = await api.get<{ success: boolean; data: RecipientPreviewResponse }>(
    `${BASE}/recipients?target_type=${targetType}`
  );
  return res.data;
}

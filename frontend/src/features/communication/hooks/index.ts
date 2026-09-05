// hooks/index.ts — Communication module TanStack Query hooks

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  fetchDevices,
  createDevice,
  updateDeviceStatus,
  getDeviceQR,
  deleteDevice,
  fetchTemplates,
  createTemplate,
  updateTemplate,
  deleteTemplate,
  fetchBlasts,
  fetchBlast,
  createBlast,
  fetchBlastLogs,
  sendBlast,
  fetchRecipientPreview,
} from "../api";
import type {
  CreateDeviceRequest,
  CreateTemplateRequest,
  UpdateTemplateRequest,
  CreateBlastRequest,
} from "../types";

const DEVICES_KEY = "wa-devices";
const TEMPLATES_KEY = "msg-templates";
const BLASTS_KEY = "blast-messages";

// ── Devices ───────────────────────────────────────────────────────────────────

export function useDevices() {
  return useQuery({ queryKey: [DEVICES_KEY], queryFn: fetchDevices });
}

export function useCreateDevice() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateDeviceRequest) => createDevice(data),
    onSuccess: () => qc.invalidateQueries({ queryKey: [DEVICES_KEY] }),
  });
}

export function useUpdateDeviceStatus() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, status, phone }: { id: string; status: string; phone?: string }) =>
      updateDeviceStatus(id, status, phone),
    onSuccess: () => qc.invalidateQueries({ queryKey: [DEVICES_KEY] }),
  });
}

export function useDeviceQR() {
  return useMutation({
    mutationFn: (id: string) => getDeviceQR(id),
  });
}

export function useDeleteDevice() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => deleteDevice(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: [DEVICES_KEY] }),
  });
}

// ── Templates ─────────────────────────────────────────────────────────────────

export function useTemplates() {
  return useQuery({ queryKey: [TEMPLATES_KEY], queryFn: fetchTemplates });
}

export function useCreateTemplate() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateTemplateRequest) => createTemplate(data),
    onSuccess: () => qc.invalidateQueries({ queryKey: [TEMPLATES_KEY] }),
  });
}

export function useUpdateTemplate() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateTemplateRequest }) =>
      updateTemplate(id, data),
    onSuccess: () => qc.invalidateQueries({ queryKey: [TEMPLATES_KEY] }),
  });
}

export function useDeleteTemplate() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => deleteTemplate(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: [TEMPLATES_KEY] }),
  });
}

// ── Blast Messages ────────────────────────────────────────────────────────────

export function useBlasts() {
  return useQuery({ queryKey: [BLASTS_KEY], queryFn: fetchBlasts });
}

export function useBlast(id: string | undefined) {
  return useQuery({
    queryKey: [BLASTS_KEY, id],
    queryFn: () => fetchBlast(id!),
    enabled: !!id,
  });
}

export function useCreateBlast() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateBlastRequest) => createBlast(data),
    onSuccess: () => qc.invalidateQueries({ queryKey: [BLASTS_KEY] }),
  });
}

export function useBlastLogs(id: string | undefined) {
  return useQuery({
    queryKey: [BLASTS_KEY, id, "logs"],
    queryFn: () => fetchBlastLogs(id!),
    enabled: !!id,
  });
}

export function useSendBlast() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => sendBlast(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: [BLASTS_KEY] }),
  });
}

export function useRecipientPreview(targetType: string) {
  return useQuery({
    queryKey: ["recipient-preview", targetType],
    queryFn: () => fetchRecipientPreview(targetType),
    enabled: !!targetType,
  });
}

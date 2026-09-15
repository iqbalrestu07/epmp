import { api } from "@/services/api";
import type { NotificationListResponse } from "../types";

const BASE_PATH = "/notifications";

export async function fetchNotifications(page = 1, perPage = 20): Promise<NotificationListResponse> {
  const res = await api.get<{ success: boolean; data: NotificationListResponse }>(
    `${BASE_PATH}?page=${page}&per_page=${perPage}`
  );
  return res.data;
}

export async function markNotificationRead(id: string): Promise<void> {
  await api.post(`${BASE_PATH}/${id}/read`, {});
}

export async function markAllNotificationsRead(): Promise<void> {
  await api.post(`${BASE_PATH}/read-all`, {});
}

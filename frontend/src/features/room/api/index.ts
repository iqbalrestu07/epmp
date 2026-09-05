import { api, ApiError } from "@/services/api";
import type {
  Room,
  RoomListResponse,
  CreateRoomRequest,
  UpdateRoomRequest,
  RoomQueryParams,
} from "../types";

const BASE_PATH = "/rooms";
export async function fetchRooms(params?: RoomQueryParams): Promise<RoomListResponse> {
  const query = new URLSearchParams();
  if (params?.page) query.set("page", String(params.page));
  if (params?.per_page) query.set("per_page", String(params.per_page));
  if (params?.search) query.set("search", params.search);
  if (params?.floor_id) query.set("floor_id", params.floor_id);
  const qs = query.toString();
  const res = await api.get<{ success: boolean; data: RoomListResponse }>(`${BASE_PATH}${qs ? `?${qs}` : ""}`);
  return res.data;
}
export async function fetchRoomById(id: string): Promise<Room> {
  const res = await api.get<{ success: boolean; data: Room }>(`${BASE_PATH}/${id}`);
  return res.data;
}
export async function createRoom(data: CreateRoomRequest): Promise<Room> {
  const res = await api.post<{ success: boolean; data: Room }>(BASE_PATH, data);
  return res.data;
}
export async function updateRoom(id: string, data: UpdateRoomRequest): Promise<Room> {
  const res = await api.put<{ success: boolean; data: Room }>(`${BASE_PATH}/${id}`, data);
  return res.data;
}
export async function deleteRoom(id: string): Promise<void> {
  return api.delete<void>(`${BASE_PATH}/${id}`);
}

export { ApiError };

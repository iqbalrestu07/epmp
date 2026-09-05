import { api, ApiError } from "@/services/api";
import type {
  Floor,
  FloorListResponse,
  CreateFloorRequest,
  UpdateFloorRequest,
  FloorQueryParams,
} from "../types";

const BASE_PATH = "/floors";
export async function fetchFloors(params?: FloorQueryParams): Promise<FloorListResponse> {
  const query = new URLSearchParams();
  if (params?.page) query.set("page", String(params.page));
  if (params?.per_page) query.set("per_page", String(params.per_page));
  if (params?.search) query.set("search", params.search);
  if (params?.building_id) query.set("building_id", params.building_id);
  const qs = query.toString();
  const res = await api.get<{ success: boolean; data: FloorListResponse }>(`${BASE_PATH}${qs ? `?${qs}` : ""}`);
  return res.data;
}
export async function fetchFloorById(id: string): Promise<Floor> {
  const res = await api.get<{ success: boolean; data: Floor }>(`${BASE_PATH}/${id}`);
  return res.data;
}
export async function createFloor(data: CreateFloorRequest): Promise<Floor> {
  const res = await api.post<{ success: boolean; data: Floor }>(BASE_PATH, data);
  return res.data;
}
export async function updateFloor(id: string, data: UpdateFloorRequest): Promise<Floor> {
  const res = await api.put<{ success: boolean; data: Floor }>(`${BASE_PATH}/${id}`, data);
  return res.data;
}
export async function deleteFloor(id: string): Promise<void> {
  return api.delete<void>(`${BASE_PATH}/${id}`);
}

export { ApiError };

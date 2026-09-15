import { api } from "@/services/api";
import type {
  OccupancyReportResponse,
  RevenueReportResponse,
  ArAgingReportResponse,
} from "../types";

const BASE_PATH = "/reports";

export async function fetchOccupancyReport(params?: {
  property_id?: string;
  building_id?: string;
}): Promise<OccupancyReportResponse> {
  const query = new URLSearchParams();
  if (params?.property_id) query.set("property_id", params.property_id);
  if (params?.building_id) query.set("building_id", params.building_id);
  const qs = query.toString();
  const res = await api.get<{ success: boolean; data: OccupancyReportResponse }>(
    `${BASE_PATH}/occupancy${qs ? `?${qs}` : ""}`
  );
  return res.data;
}

export async function fetchRevenueReport(params?: {
  from?: string;
  to?: string;
}): Promise<RevenueReportResponse> {
  const query = new URLSearchParams();
  if (params?.from) query.set("from", params.from);
  if (params?.to) query.set("to", params.to);
  const qs = query.toString();
  const res = await api.get<{ success: boolean; data: RevenueReportResponse }>(
    `${BASE_PATH}/revenue${qs ? `?${qs}` : ""}`
  );
  return res.data;
}

export async function fetchArAgingReport(): Promise<ArAgingReportResponse> {
  const res = await api.get<{ success: boolean; data: ArAgingReportResponse }>(
    `${BASE_PATH}/ar-aging`
  );
  return res.data;
}

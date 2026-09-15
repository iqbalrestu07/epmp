import { useQuery } from "@tanstack/react-query";
import {
  fetchOccupancyReport,
  fetchRevenueReport,
  fetchArAgingReport,
} from "../api";

const QUERY_KEY = "reports";

export function useOccupancyReport(params?: {
  property_id?: string;
  building_id?: string;
}) {
  return useQuery({
    queryKey: [QUERY_KEY, "occupancy", params],
    queryFn: () => fetchOccupancyReport(params),
  });
}

export function useRevenueReport(params?: { from?: string; to?: string }) {
  return useQuery({
    queryKey: [QUERY_KEY, "revenue", params],
    queryFn: () => fetchRevenueReport(params),
  });
}

export function useArAgingReport() {
  return useQuery({
    queryKey: [QUERY_KEY, "ar-aging"],
    queryFn: fetchArAgingReport,
  });
}

export interface OccupancyBuildingRow {
  building_id: string;
  building_name: string;
  total_rooms: number;
  occupied: number;
  reserved: number;
  maintenance: number;
  available: number;
  occupancy_rate: number;
}

export interface OccupancyPropertyRow {
  property_id: string;
  property_name: string;
  total_rooms: number;
  occupied: number;
  reserved: number;
  maintenance: number;
  available: number;
  occupancy_rate: number;
  buildings: OccupancyBuildingRow[];
}

export interface OccupancyTotals {
  total_rooms: number;
  occupied: number;
  reserved: number;
  maintenance: number;
  available: number;
  occupancy_rate: number;
}

export interface OccupancyReportResponse {
  totals: OccupancyTotals;
  properties: OccupancyPropertyRow[];
}

export interface RevenueMonthlyRow {
  period: string;
  currency: string;
  invoice_count: number;
  invoiced_amount: number;
  paid_invoice_count: number;
  payment_count: number;
  received_amount: number;
}

export interface RevenueReportResponse {
  from: string;
  to: string;
  rows: RevenueMonthlyRow[];
}

export interface ArAgingInvoice {
  invoice_id: string;
  tenant_id: string;
  tenant_name: string;
  amount: number;
  paid_amount: number;
  outstanding: number;
  currency: string;
  status: string;
  due_date: string;
  days_overdue: number;
  bucket: string;
}

export interface ArAgingBucketRow {
  bucket: string;
  currency: string;
  invoice_count: number;
  outstanding: number;
}

export interface ArAgingReportResponse {
  buckets: ArAgingBucketRow[];
  invoices: ArAgingInvoice[];
}

export interface Floor {
  id: string;
  organization_id: string;
  building_id: string;
  name: string;
  floor_number: number;
  is_active: boolean;
  deleted_at?: string;
  created_at: string;
  updated_at: string;
}

export interface FloorListResponse {
  data: Floor[];
  total: number;
  page: number;
  per_page: number;
  total_pages: number;
}

export interface CreateFloorRequest {
  building_id: string;
  name: string;
  floor_number: number;
  is_active: boolean;
}

export interface UpdateFloorRequest {
  building_id: string;
  name: string;
  floor_number: number;
  is_active: boolean;
}

export interface FloorQueryParams {
  page?: number;
  per_page?: number;
  sort?: string;
  order?: "asc" | "desc";
  search?: string;
  building_id?: string;
}

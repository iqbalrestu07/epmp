export interface Building {
  id: string;
  property_id: string;
  name: string;
  total_floors: number;
  deleted_at?: string;
  created_at: string;
  updated_at: string;
}

export interface BuildingListResponse {
  data: Building[];
  total: number;
  page: number;
  per_page: number;
  total_pages: number;
}

export interface CreateBuildingRequest {
  property_id: string;
  name: string;
  total_floors: number;
}

export interface UpdateBuildingRequest {
  property_id: string;
  name: string;
  total_floors: number;
}

export interface BuildingQueryParams {
  page?: number;
  per_page?: number;
  sort?: string;
  order?: "asc" | "desc";
  search?: string;
  property_id?: string;
  name?: string;
}

export interface Asset {
  organization_id: string;
  id: string;
  property_id: string;
  name: string;
  category: string;
  status: string;
  purchase_price: number;
  deleted_at?: string;
  created_at: string;
  updated_at: string;
}

export interface AssetListResponse {
  data: Asset[];
  total: number;
  page: number;
  per_page: number;
  total_pages: number;
}

export interface CreateAssetRequest {
  property_id: string;
  name: string;
  category: string;
  status: string;
  purchase_price: number;
}

export interface UpdateAssetRequest {
  property_id?: string;
  name?: string;
  category?: string;
  status?: string;
  purchase_price?: number;
}

export interface AssetQueryParams {
  page?: number;
  per_page?: number;
  search?: string;
  property_id?: string;
  name?: string;
}

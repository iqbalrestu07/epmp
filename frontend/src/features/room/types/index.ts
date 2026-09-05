export interface Room {
  organization_id: string;
  id: string;
  name: string;
  floor_id: string;
  room_type_id?: string;
  capacity: number;
  price: number;
  is_available: boolean;
  property_id: string;
  deleted_at?: string;
  created_at: string;
  updated_at: string;
}

export interface RoomListResponse {
  data: Room[];
  total: number;
  page: number;
  per_page: number;
  total_pages: number;
}

export interface CreateRoomRequest {
  name: string;
  property_id: string;
  floor_id?: string | null;
  room_type_id?: string | null;
  capacity: number;
  price: number;
  is_available: boolean;
}

export interface UpdateRoomRequest {
  name?: string;
  property_id?: string;
  floor_id?: string | null;
  room_type_id?: string | null;
  capacity?: number;
  price?: number;
  is_available?: boolean;
}

export interface RoomQueryParams {
  page?: number;
  per_page?: number;
  search?: string;
  floor_id?: string;
}

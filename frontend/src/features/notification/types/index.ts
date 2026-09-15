export interface AppNotification {
  id: string;
  organization_id: string;
  user_id?: string;
  type: string;
  title: string;
  message: string;
  link?: string;
  is_read: boolean;
  created_at: string;
}

export interface NotificationListResponse {
  data: AppNotification[];
  total: number;
  unread_count: number;
  page: number;
  per_page: number;
  total_pages: number;
}

// IAM Types — Auth, User, Role, Permission

export interface Permission {
  id: string;
  resource: string;
  action: string;
  description: string;
  key: string; // "resource:action"
}

export interface Role {
  id: string;
  name: string;
  description: string;
  is_system: boolean;
  permissions?: Permission[];
  user_count?: number;
}

export interface User {
  id: string;
  email: string;
  name: string;
  organization_id?: string;
  is_active: boolean;
  roles: Role[];
  permissions: string[]; // ["property:read", "property:write", ...]
}

export interface TokenResponse {
  access_token: string;
  refresh_token: string;
  token_type: string;
  expires_in: number;
  user: User;
}

export interface UserListResponse {
  data: User[];
  total: number;
  page: number;
  per_page: number;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  name: string;
  email: string;
  password: string;
}

export interface CreateUserRequest {
  name: string;
  email: string;
  password: string;
  role_ids?: string[];
}

export interface UpdateUserRequest {
  name: string;
  is_active?: boolean;
}

export interface AssignRolesRequest {
  role_ids: string[];
}

export interface CreateRoleRequest {
  name: string;
  description: string;
  permission_ids?: string[];
}

export interface UpdateRoleRequest {
  name: string;
  description: string;
}

export interface SetPermissionsRequest {
  permission_ids: string[];
}

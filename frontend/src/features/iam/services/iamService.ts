import { api } from '../../../services/api';
import type {
  TokenResponse,
  LoginRequest,
  RegisterRequest,
  User,
  UserListResponse,
  CreateUserRequest,
  UpdateUserRequest,
  AssignRolesRequest,
  Role,
  CreateRoleRequest,
  UpdateRoleRequest,
  SetPermissionsRequest,
  Permission,
} from '../types';

const AUTH_BASE = '/v1/auth';
const USERS_BASE = '/v1/users';
const ROLES_BASE = '/v1/roles';
const PERMISSIONS_BASE = '/v1/permissions';

// ─── Auth ──────────────────────────────────────────────────────────────────

export const authService = {
  login: (data: LoginRequest) =>
    api.post<{ success: boolean; data: TokenResponse }>(`${AUTH_BASE}/login`, data),

  register: (data: RegisterRequest) =>
    api.post<{ success: boolean; data: TokenResponse }>(`${AUTH_BASE}/register`, data),

  logout: (refreshToken: string) =>
    api.post<void>(`${AUTH_BASE}/logout`, { refresh_token: refreshToken }),

  refresh: (refreshToken: string) =>
    api.post<{ success: boolean; data: TokenResponse }>(`${AUTH_BASE}/refresh`, {
      refresh_token: refreshToken,
    }),

  me: () =>
    api.get<{ success: boolean; data: User }>(`${AUTH_BASE}/me`),
};

// ─── Users ────────────────────────────────────────────────────────────────

export const userService = {
  list: (page = 1, perPage = 20) =>
    api.get<{ success: boolean; data: UserListResponse }>(
      `${USERS_BASE}?page=${page}&per_page=${perPage}`
    ),

  getById: (id: string) =>
    api.get<{ success: boolean; data: User }>(`${USERS_BASE}/${id}`),

  create: (data: CreateUserRequest) =>
    api.post<{ success: boolean; data: User }>(USERS_BASE, data),

  update: (id: string, data: UpdateUserRequest) =>
    api.put<{ success: boolean; data: User }>(`${USERS_BASE}/${id}`, data),

  delete: (id: string) =>
    api.delete<void>(`${USERS_BASE}/${id}`),

  assignRoles: (id: string, data: AssignRolesRequest) =>
    api.put<{ success: boolean; data: User }>(`${USERS_BASE}/${id}/roles`, data),
};

// ─── Roles ────────────────────────────────────────────────────────────────

export const roleService = {
  list: () =>
    api.get<{ success: boolean; data: Role[] }>(ROLES_BASE),

  getById: (id: string) =>
    api.get<{ success: boolean; data: Role }>(`${ROLES_BASE}/${id}`),

  create: (data: CreateRoleRequest) =>
    api.post<{ success: boolean; data: Role }>(ROLES_BASE, data),

  update: (id: string, data: UpdateRoleRequest) =>
    api.put<{ success: boolean; data: Role }>(`${ROLES_BASE}/${id}`, data),

  delete: (id: string) =>
    api.delete<void>(`${ROLES_BASE}/${id}`),

  setPermissions: (id: string, data: SetPermissionsRequest) =>
    api.put<{ success: boolean; data: Role }>(`${ROLES_BASE}/${id}/permissions`, data),
};

// ─── Permissions ──────────────────────────────────────────────────────────

export const permissionService = {
  list: () =>
    api.get<{ success: boolean; data: Permission[] }>(PERMISSIONS_BASE),
};

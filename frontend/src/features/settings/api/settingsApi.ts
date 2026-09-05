import { api } from '@/services/api';

export interface OrgMember {
  id: string;
  organization_id: string;
  user_id: string;
  user_name: string;
  user_email: string;
  role: string;
  joined_at: string;
  is_active: boolean;
}

export interface AddOrgMemberRequest {
  email: string;
  name?: string;
  password?: string;
  role?: string;
}

export interface PropertyStaff {
  id: string;
  property_id: string;
  user_id: string;
  user_name: string;
  user_email: string;
  role_id: string;
  role_name: string;
  is_system: boolean;
  created_at: string;
}

export interface AssignPropertyStaffRequest {
  user_id: string;
  role_id: string;
}

export const settingsApi = {
  // Organization Members
  getOrgMembers: async (): Promise<OrgMember[]> => {
    const res = await api.get<{ success: boolean; data: OrgMember[] }>('/organizations/members');
    return res.data ?? (Array.isArray(res) ? res : []);
  },

  addOrgMember: async (data: AddOrgMemberRequest): Promise<OrgMember> => {
    const res = await api.post<{ success: boolean; data: OrgMember }>('/organizations/members', data);
    return res.data;
  },

  removeOrgMember: async (userId: string): Promise<void> => {
    await api.delete(`/organizations/members/${userId}`);
  },

  // Property Staff RBAC
  getPropertyStaff: async (propertyId: string): Promise<PropertyStaff[]> => {
    const res = await api.get<{ success: boolean; data: PropertyStaff[] }>(`/properties/${propertyId}/staff`);
    return res.data ?? (Array.isArray(res) ? res : []);
  },

  assignPropertyStaff: async (propertyId: string, data: AssignPropertyStaffRequest): Promise<PropertyStaff> => {
    const res = await api.post<{ success: boolean; data: PropertyStaff }>(`/properties/${propertyId}/staff`, data);
    return res.data;
  },

  removePropertyStaff: async (propertyId: string, userId: string, roleId: string): Promise<void> => {
    await api.delete(`/properties/${propertyId}/staff/${userId}/roles/${roleId}`);
  },
};

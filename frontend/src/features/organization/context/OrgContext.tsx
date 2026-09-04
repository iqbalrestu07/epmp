import React, { createContext, useContext, useEffect, useState, useCallback } from 'react';
import { api, getStoredOrgId, setStoredOrgId, clearStoredOrgId } from '../../../services/api';
import type { Organization } from '../types';

interface OrgContextValue {
  currentOrg: Organization | null;
  orgs: Organization[];
  orgId: string | null;
  isLoading: boolean;
  switchOrg: (orgId: string) => void;
  refreshOrgs: () => Promise<void>;
}

const OrgContext = createContext<OrgContextValue | null>(null);

export function OrgProvider({ children }: { children: React.ReactNode }) {
  const [orgs, setOrgs] = useState<Organization[]>([]);
  const [currentOrg, setCurrentOrg] = useState<Organization | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [orgId, setOrgId] = useState<string | null>(getStoredOrgId());

  const fetchOrgs = useCallback(async () => {
    try {
      const res = await api.get<{ data: Organization[] }>('/organizations/mine');
      const list = res.data ?? [];
      setOrgs(list);

      const stored = getStoredOrgId();
      if (stored && list.length > 0) {
        const found = list.find(o => o.id === stored);
        if (found) {
          setCurrentOrg(found);
          setOrgId(found.id);
        } else {
          setCurrentOrg(list[0]);
          setOrgId(list[0].id);
          setStoredOrgId(list[0].id);
        }
      } else if (list.length > 0) {
        setCurrentOrg(list[0]);
        setOrgId(list[0].id);
        setStoredOrgId(list[0].id);
      } else {
        setCurrentOrg(null);
        setOrgId(null);
        clearStoredOrgId();
      }
    } catch {
      setOrgs([]);
      setCurrentOrg(null);
      setOrgId(null);
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchOrgs();
  }, [fetchOrgs]);

  const switchOrg = useCallback((id: string) => {
    const found = orgs.find(o => o.id === id);
    if (found) {
      setCurrentOrg(found);
      setOrgId(found.id);
      setStoredOrgId(found.id);
    }
  }, [orgs]);

  return (
    <OrgContext.Provider
      value={{
        currentOrg,
        orgs,
        orgId,
        isLoading,
        switchOrg,
        refreshOrgs: fetchOrgs,
      }}
    >
      {children}
    </OrgContext.Provider>
  );
}

export function useOrg() {
  const ctx = useContext(OrgContext);
  if (!ctx) throw new Error('useOrg must be used within OrgProvider');
  return ctx;
}

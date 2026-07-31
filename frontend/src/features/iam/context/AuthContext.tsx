import React, { createContext, useContext, useEffect, useState, useCallback } from 'react';
import type { User } from '../types';
import { authService } from '../services/iamService';
import {
  getStoredToken,
  setStoredTokens,
  clearStoredTokens,
  REFRESH_TOKEN_KEY,
} from '../../../services/api';

// ─── Types ───────────────────────────────────────────────────────────────────

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
}

interface AuthContextValue extends AuthState {
  login: (email: string, password: string) => Promise<void>;
  register: (name: string, email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  hasPermission: (permission: string) => boolean;
  hasRole: (roleName: string) => boolean;
}

// ─── Context ─────────────────────────────────────────────────────────────────

const AuthContext = createContext<AuthContextValue | null>(null);

// ─── Provider ─────────────────────────────────────────────────────────────────

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [state, setState] = useState<AuthState>({
    user: null,
    isAuthenticated: false,
    isLoading: true, // true on mount until we verify the stored token
  });

  // On mount: if there's a stored access token, fetch the current user profile.
  useEffect(() => {
    const token = getStoredToken();
    if (!token) {
      setState(prev => ({ ...prev, isLoading: false }));
      return;
    }

    authService
      .me()
      .then(res => {
        setState({ user: res.data, isAuthenticated: true, isLoading: false });
      })
      .catch(() => {
        // Token might be expired — try refresh
        const refreshToken = localStorage.getItem(REFRESH_TOKEN_KEY);
        if (!refreshToken) {
          clearStoredTokens();
          setState({ user: null, isAuthenticated: false, isLoading: false });
          return;
        }

        authService
          .refresh(refreshToken)
          .then(res => {
            setStoredTokens(res.data.access_token, res.data.refresh_token);
            setState({ user: res.data.user, isAuthenticated: true, isLoading: false });
          })
          .catch(() => {
            clearStoredTokens();
            setState({ user: null, isAuthenticated: false, isLoading: false });
          });
      });
  }, []);

  const login = useCallback(async (email: string, password: string) => {
    const res = await authService.login({ email, password });
    setStoredTokens(res.data.access_token, res.data.refresh_token);
    setState({ user: res.data.user, isAuthenticated: true, isLoading: false });
  }, []);

  const register = useCallback(async (name: string, email: string, password: string) => {
    const res = await authService.register({ name, email, password });
    setStoredTokens(res.data.access_token, res.data.refresh_token);
    setState({ user: res.data.user, isAuthenticated: true, isLoading: false });
  }, []);

  const logout = useCallback(async () => {
    const refreshToken = localStorage.getItem(REFRESH_TOKEN_KEY) ?? '';
    try {
      await authService.logout(refreshToken);
    } catch {
      // ignore errors on logout
    }
    clearStoredTokens();
    setState({ user: null, isAuthenticated: false, isLoading: false });
  }, []);

  const hasPermission = useCallback(
    (permission: string) => {
      if (!state.user) return false;
      return state.user.permissions.includes(permission);
    },
    [state.user]
  );

  const hasRole = useCallback(
    (roleName: string) => {
      if (!state.user) return false;
      return state.user.roles.some(r => r.name === roleName);
    },
    [state.user]
  );

  return (
    <AuthContext.Provider
      value={{
        ...state,
        login,
        register,
        logout,
        hasPermission,
        hasRole,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

// ─── Hook ─────────────────────────────────────────────────────────────────────

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error('useAuth must be used inside <AuthProvider>');
  }
  return ctx;
}

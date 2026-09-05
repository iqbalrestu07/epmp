const BASE_URL = "/api/v1";

export class ApiError extends Error {
  constructor(
    public status: number,
    message: string,
    public data?: unknown
  ) {
    super(message);
    this.name = "ApiError";
  }
}

// Token storage keys
export const TOKEN_KEY = "epmp_access_token";
export const REFRESH_TOKEN_KEY = "epmp_refresh_token";
export const ORG_ID_KEY = "epmp_org_id";

export function getStoredToken(): string | null {
  return localStorage.getItem(TOKEN_KEY);
}

export function setStoredTokens(accessToken: string, refreshToken: string): void {
  localStorage.setItem(TOKEN_KEY, accessToken);
  localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken);
}

export function clearStoredTokens(): void {
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(REFRESH_TOKEN_KEY);
  localStorage.removeItem(ORG_ID_KEY);
}

export function getStoredOrgId(): string | null {
  return localStorage.getItem(ORG_ID_KEY);
}

export function setStoredOrgId(orgId: string): void {
  localStorage.setItem(ORG_ID_KEY, orgId);
}

export function clearStoredOrgId(): void {
  localStorage.removeItem(ORG_ID_KEY);
}

async function request<T>(
  path: string,
  options?: RequestInit
): Promise<T> {
  const token = getStoredToken();

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(options?.headers as Record<string, string>),
  };

  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const orgId = getStoredOrgId();
  if (orgId) {
    headers["X-Organization-ID"] = orgId;
  }

  const res = await fetch(`${BASE_URL}${path}`, {
    ...options,
    headers,
  });

  if (!res.ok) {
    const body = await res.json().catch(() => ({ message: res.statusText }));
    const message = body?.error?.message ?? body?.message ?? res.statusText;
    throw new ApiError(res.status, message, body);
  }

  if (res.status === 204) {
    return undefined as T;
  }

  const json = await res.json();
  return json as T;
}

export const api = {
  get: <T>(path: string) => request<T>(path),
  post: <T>(path: string, body: unknown) =>
    request<T>(path, { method: "POST", body: JSON.stringify(body) }),
  put: <T>(path: string, body: unknown) =>
    request<T>(path, { method: "PUT", body: JSON.stringify(body) }),
  patch: <T>(path: string, body: unknown) =>
    request<T>(path, { method: "PATCH", body: JSON.stringify(body) }),
  delete: <T>(path: string) => request<T>(path, { method: "DELETE" }),
};

import { Navigate, useLocation } from 'react-router-dom';
import { useAuth } from '../features/iam/context/AuthContext';
import { useOrg } from '../features/organization/context/OrgContext';

interface ProtectedRouteProps {
  children: React.ReactNode;
  /** If provided, user must have this permission to access the route. */
  permission?: string;
  /** If provided, user must have this role to access the route. */
  role?: string;
  /** Where to redirect unauthorized users. Defaults to /auth/signin. */
  redirectTo?: string;
}

export default function ProtectedRoute({
  children,
  permission,
  role,
  redirectTo = '/auth/signin',
}: ProtectedRouteProps) {
  const { isAuthenticated, isLoading, hasPermission, hasRole } = useAuth();
  const { orgs, isLoading: orgLoading } = useOrg();
  const location = useLocation();

  if (isLoading || orgLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-[#f2efe9]">
        <div className="flex flex-col items-center gap-3">
          <div className="w-8 h-8 rounded-full border-2 border-orange border-t-transparent animate-spin" />
          <p className="text-sm text-black/50">Loading…</p>
        </div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to={redirectTo} state={{ from: location }} replace />;
  }

  // If user has no organization, force them to create one first
  // (unless they're already on the organization creation page)
  const isOnOrgCreate = location.pathname === '/dashboard/organizations/new';
  if (orgs.length === 0 && !isOnOrgCreate) {
    return <Navigate to="/dashboard/organizations/new" replace />;
  }

  if (permission && !hasPermission(permission)) {
    return <Navigate to="/dashboard" replace />;
  }

  if (role && !hasRole(role)) {
    return <Navigate to="/dashboard" replace />;
  }

  return <>{children}</>;
}

import { Navigate, useLocation } from 'react-router-dom';
import { useAuth } from '../features/iam/context/AuthContext';

interface ProtectedRouteProps {
  children: React.ReactNode;
  /** If provided, user must have this permission to access the route. */
  permission?: string;
  /** If provided, user must have this role to access the route. */
  role?: string;
  /** Where to redirect unauthorized users. Defaults to /auth/signin. */
  redirectTo?: string;
}

/**
 * ProtectedRoute wraps a route that requires authentication.
 * Optionally also enforces a specific permission or role.
 *
 * Usage:
 *   <ProtectedRoute>
 *     <Dashboard />
 *   </ProtectedRoute>
 *
 *   <ProtectedRoute permission="role:read">
 *     <RBACPage />
 *   </ProtectedRoute>
 */
export default function ProtectedRoute({
  children,
  permission,
  role,
  redirectTo = '/auth/signin',
}: ProtectedRouteProps) {
  const { isAuthenticated, isLoading, hasPermission, hasRole } = useAuth();
  const location = useLocation();

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-[#f2efe9]">
        <div className="flex flex-col items-center gap-3">
          <div className="w-8 h-8 rounded-full border-2 border-orange border-t-transparent animate-spin" />
          <p className="text-sm text-black/50">Authenticating…</p>
        </div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to={redirectTo} state={{ from: location }} replace />;
  }

  if (permission && !hasPermission(permission)) {
    return <Navigate to="/dashboard" replace />;
  }

  if (role && !hasRole(role)) {
    return <Navigate to="/dashboard" replace />;
  }

  return <>{children}</>;
}

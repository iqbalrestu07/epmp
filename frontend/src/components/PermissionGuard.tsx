import { useAuth } from '../features/iam/context/AuthContext';

interface PermissionGuardProps {
  /** The permission key required, e.g. "role:write" */
  permission?: string;
  /** The role name required, e.g. "super_admin" */
  role?: string;
  /** Rendered when the user has the required permission/role. */
  children: React.ReactNode;
  /** Optional fallback content when permission is denied. */
  fallback?: React.ReactNode;
}

/**
 * PermissionGuard conditionally renders children based on the
 * authenticated user's permissions or roles.
 *
 * Unlike ProtectedRoute (which redirects), PermissionGuard simply
 * hides or shows UI elements in-place.
 *
 * Usage:
 *   <PermissionGuard permission="user:write">
 *     <button>Create User</button>
 *   </PermissionGuard>
 */
export default function PermissionGuard({
  permission,
  role,
  children,
  fallback = null,
}: PermissionGuardProps) {
  const { hasPermission, hasRole } = useAuth();

  if (permission && !hasPermission(permission)) {
    return <>{fallback}</>;
  }

  if (role && !hasRole(role)) {
    return <>{fallback}</>;
  }

  return <>{children}</>;
}

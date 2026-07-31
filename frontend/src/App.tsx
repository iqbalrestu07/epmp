import { Routes, Route, Navigate } from "react-router-dom";
import { AuthProvider } from "./features/iam/context/AuthContext";
import ProtectedRoute from "./components/ProtectedRoute";

import MainLayout from "./layouts/MainLayout";
import DashboardPage from "./pages/DashboardPage";
import ImmersiveLanding from "./features/immersive-view/pages/ImmersiveLanding";

import AuthLayout from "./layouts/AuthLayout";
import SignInPage from "./features/iam/pages/SignInPage";
import SignUpPage from "./features/iam/pages/SignUpPage";

import { PropertyListPage } from "./features/property/pages/PropertyListPage";
import PropertyInteractiveView from "./features/property/pages/PropertyInteractiveView";
import { RoomListPage } from "./features/room/pages/RoomListPage";
import { TenantListPage } from "./features/tenant/pages/TenantListPage";
import RBACPage from "./features/iam/pages/RBACPage";
import UserListPage from "./features/iam/pages/UserListPage";

export default function App() {
  return (
    <AuthProvider>
      <Routes>
        {/* Public Routes */}
        <Route path="/" element={<ImmersiveLanding />} />

        {/* Auth Routes */}
        <Route path="/auth" element={<AuthLayout />}>
          <Route path="signin" element={<SignInPage />} />
          <Route path="signup" element={<SignUpPage />} />
        </Route>

        {/* Protected App Routes */}
        <Route
          path="/dashboard"
          element={
            <ProtectedRoute>
              <MainLayout />
            </ProtectedRoute>
          }
        >
          <Route index element={<DashboardPage />} />
          <Route path="properties" element={<PropertyListPage />} />
          <Route path="properties/interactive" element={<PropertyInteractiveView />} />
          <Route path="rooms" element={<RoomListPage />} />
          <Route path="tenants" element={<TenantListPage />} />

          {/* Management — require specific permissions */}
          <Route
            path="management/rbac"
            element={
              <ProtectedRoute permission="role:read">
                <RBACPage />
              </ProtectedRoute>
            }
          />
          <Route
            path="management/users"
            element={
              <ProtectedRoute permission="user:read">
                <UserListPage />
              </ProtectedRoute>
            }
          />

          <Route path="settings" element={<div className="p-8">Settings Page (Coming Soon)</div>} />
          <Route path="*" element={<Navigate to="/dashboard" replace />} />
        </Route>

        {/* Catch-all */}
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </AuthProvider>
  );
}

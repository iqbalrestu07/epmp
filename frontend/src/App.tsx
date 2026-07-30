import { Routes, Route, Navigate } from "react-router-dom";
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
    <Routes>
      {/* Public Routes */}
      <Route path="/" element={<ImmersiveLanding />} />
      
      {/* Auth Routes */}
      <Route path="/auth" element={<AuthLayout />}>
        <Route path="signin" element={<SignInPage />} />
        <Route path="signup" element={<SignUpPage />} />
      </Route>

      {/* Protected/App Routes */}
      <Route path="/dashboard" element={<MainLayout />}>
        <Route index element={<DashboardPage />} />
        <Route path="properties" element={<PropertyListPage />} />
        <Route path="properties/interactive" element={<PropertyInteractiveView />} />
        <Route path="rooms" element={<RoomListPage />} />
        <Route path="tenants" element={<TenantListPage />} />
        
        {/* Management & RBAC */}
        <Route path="management/rbac" element={<RBACPage />} />
        <Route path="management/users" element={<UserListPage />} />
        
        {/* Fallback for unconfigured routes */}
        <Route path="settings" element={<div className="p-8">Settings Page (Coming Soon)</div>} />
        <Route path="*" element={<Navigate to="/dashboard" replace />} />
      </Route>
    </Routes>
  );
}

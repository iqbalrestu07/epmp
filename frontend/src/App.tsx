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
import { ReservationListPage } from "./features/reservation/pages/ReservationListPage";
import { ContractListPage } from "./features/contract/pages/ContractListPage";
import { OccupancyListPage } from "./features/occupancy/pages/OccupancyListPage";
import { InvoiceListPage } from "./features/billing/pages/InvoiceListPage";
import { PaymentListPage } from "./features/payment/pages/PaymentListPage";
import { DepositListPage } from "./features/deposit/pages/DepositListPage";
import { ChargeListPage } from "./features/charge/pages/ChargeListPage";
import { RefundListPage } from "./features/refund/pages/RefundListPage";
import { AdjustmentListPage } from "./features/adjustment/pages/AdjustmentListPage";
import { PenaltyListPage } from "./features/penalty/pages/PenaltyListPage";

import { BuildingListPage } from "./features/building/pages/BuildingListPage";
import { ZoneListPage } from "./features/zone/pages/ZoneListPage";
import { BedListPage } from "./features/bed/pages/BedListPage";
import { FacilityListPage } from "./features/facility/pages/FacilityListPage";
import { RoomTypeListPage } from "./features/roomtype/pages/RoomTypeListPage";
import { AssetListPage } from "./features/asset/pages/AssetListPage";
import { AssetAssignmentListPage } from "./features/assetassignment/pages/AssetAssignmentListPage";
import { AssetInspectionListPage } from "./features/assetinspection/pages/AssetInspectionListPage";
import { WorkOrderListPage } from "./features/workorder/pages/WorkOrderListPage";
import { TechnicianListPage } from "./features/technician/pages/TechnicianListPage";
import { SupplierListPage } from "./features/supplier/pages/SupplierListPage";
import { OrganizationListPage } from "./features/organization/pages/OrganizationListPage";

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
          
          {/* Core Property & Tenant */}
          <Route path="organizations" element={<OrganizationListPage />} />
          <Route path="properties" element={<PropertyListPage />} />
          <Route path="properties/interactive" element={<PropertyInteractiveView />} />
          <Route path="buildings" element={<BuildingListPage />} />
          <Route path="zones" element={<ZoneListPage />} />
          <Route path="rooms" element={<RoomListPage />} />
          <Route path="room-types" element={<RoomTypeListPage />} />
          <Route path="beds" element={<BedListPage />} />
          <Route path="facilities" element={<FacilityListPage />} />
          <Route path="tenants" element={<TenantListPage />} />

          {/* Operations */}
          <Route path="reservations" element={<ReservationListPage />} />
          <Route path="contracts" element={<ContractListPage />} />
          <Route path="occupancies" element={<OccupancyListPage />} />

          {/* Finance */}
          <Route path="invoices" element={<InvoiceListPage />} />
          <Route path="payments" element={<PaymentListPage />} />
          <Route path="deposits" element={<DepositListPage />} />
          <Route path="charges" element={<ChargeListPage />} />
          <Route path="refunds" element={<RefundListPage />} />
          <Route path="adjustments" element={<AdjustmentListPage />} />
          <Route path="penalties" element={<PenaltyListPage />} />

          {/* Assets & Maintenance */}
          <Route path="assets" element={<AssetListPage />} />
          <Route path="asset-assignments" element={<AssetAssignmentListPage />} />
          <Route path="asset-inspections" element={<AssetInspectionListPage />} />
          <Route path="work-orders" element={<WorkOrderListPage />} />
          <Route path="technicians" element={<TechnicianListPage />} />
          <Route path="suppliers" element={<SupplierListPage />} />

          {/* Management */}
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

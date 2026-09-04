import { Routes, Route, Navigate } from "react-router-dom";
import { AuthProvider } from "./features/iam/context/AuthContext";
import { OrgProvider } from "./features/organization/context/OrgContext";
import ProtectedRoute from "./components/ProtectedRoute";

import MainLayout from "./layouts/MainLayout";
import DashboardPage from "./pages/DashboardPage";
import ImmersiveLanding from "./features/immersive-view/pages/ImmersiveLanding";

import AuthLayout from "./layouts/AuthLayout";
import SignInPage from "./features/iam/pages/SignInPage";
import SignUpPage from "./features/iam/pages/SignUpPage";

import { PropertyListPage } from "./features/property/pages/PropertyListPage";
import { PropertyCreatePage } from "./features/property/pages/PropertyCreatePage";
import { PropertyDetailPage } from "./features/property/pages/PropertyDetailPage";
import { PropertyEditPage } from "./features/property/pages/PropertyEditPage";
import PropertyInteractiveView from "./features/property/pages/PropertyInteractiveView";
import { RoomListPage } from "./features/room/pages/RoomListPage";
import { RoomCreatePage } from "./features/room/pages/RoomCreatePage";
import { RoomDetailPage } from "./features/room/pages/RoomDetailPage";
import { RoomEditPage } from "./features/room/pages/RoomEditPage";
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
import { BuildingCreatePage } from "./features/building/pages/BuildingCreatePage";
import { BuildingDetailPage } from "./features/building/pages/BuildingDetailPage";
import { BuildingEditPage } from "./features/building/pages/BuildingEditPage";
import { FloorListPage } from "./features/floor/pages/FloorListPage";
import { FloorCreatePage } from "./features/floor/pages/FloorCreatePage";
import { FloorDetailPage } from "./features/floor/pages/FloorDetailPage";
import { FloorEditPage } from "./features/floor/pages/FloorEditPage";
import { ZoneListPage } from "./features/zone/pages/ZoneListPage";
import { BedListPage } from "./features/bed/pages/BedListPage";
import { FacilityListPage } from "./features/facility/pages/FacilityListPage";
import { RoomTypeListPage } from "./features/roomtype/pages/RoomTypeListPage";
import { AssetListPage } from "./features/asset/pages/AssetListPage";
import { AssetCreatePage } from "./features/asset/pages/AssetCreatePage";
import { AssetEditPage } from "./features/asset/pages/AssetEditPage";
import { AssetDetailPage } from "./features/asset/pages/AssetDetailPage";
import { AssetAssignmentListPage } from "./features/assetassignment/pages/AssetAssignmentListPage";
import { AssetInspectionListPage } from "./features/assetinspection/pages/AssetInspectionListPage";
import { WorkOrderListPage } from "./features/workorder/pages/WorkOrderListPage";
import { TechnicianListPage } from "./features/technician/pages/TechnicianListPage";
import { SupplierListPage } from "./features/supplier/pages/SupplierListPage";
import { OrganizationListPage } from "./features/organization/pages/OrganizationListPage";
import { OrganizationCreatePage } from "./features/organization/pages/OrganizationCreatePage";
import { OrganizationDetailPage } from "./features/organization/pages/OrganizationDetailPage";
import { OrganizationEditPage } from "./features/organization/pages/OrganizationEditPage";

import { MessagingSettingsPage } from "./features/communication/pages/MessagingSettingsPage";
import { BlastMessagePage } from "./features/communication/pages/BlastMessagePage";

import RBACPage from "./features/iam/pages/RBACPage";
import UserListPage from "./features/iam/pages/UserListPage";

export default function App() {
  return (
    <AuthProvider>
      <OrgProvider>
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
          
          {/* Core: Organization → Property → Building */}
          <Route path="organizations" element={<OrganizationListPage />} />
          <Route path="organizations/new" element={<OrganizationCreatePage />} />
          <Route path="organizations/:id" element={<OrganizationDetailPage />} />
          <Route path="organizations/:id/edit" element={<OrganizationEditPage />} />
          <Route path="properties" element={<PropertyListPage />} />
          <Route path="properties/new" element={<PropertyCreatePage />} />
          <Route path="properties/interactive" element={<PropertyInteractiveView />} />
          <Route path="properties/:id" element={<PropertyDetailPage />} />
          <Route path="properties/:id/edit" element={<PropertyEditPage />} />
          <Route path="buildings" element={<BuildingListPage />} />
          <Route path="buildings/new" element={<BuildingCreatePage />} />
          <Route path="buildings/:id" element={<BuildingDetailPage />} />
          <Route path="buildings/:id/edit" element={<BuildingEditPage />} />
          <Route path="floors" element={<FloorListPage />} />
          <Route path="floors/new" element={<FloorCreatePage />} />
          <Route path="floors/:id" element={<FloorDetailPage />} />
          <Route path="floors/:id/edit" element={<FloorEditPage />} />
          <Route path="zones" element={<ZoneListPage />} />
          <Route path="rooms" element={<RoomListPage />} />
          <Route path="rooms/new" element={<RoomCreatePage />} />
          <Route path="rooms/:id" element={<RoomDetailPage />} />
          <Route path="rooms/:id/edit" element={<RoomEditPage />} />
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
          <Route path="assets/new" element={<AssetCreatePage />} />
          <Route path="assets/:id" element={<AssetDetailPage />} />
          <Route path="assets/:id/edit" element={<AssetEditPage />} />
          <Route path="asset-assignments" element={<AssetAssignmentListPage />} />
          <Route path="asset-inspections" element={<AssetInspectionListPage />} />
          <Route path="work-orders" element={<WorkOrderListPage />} />
          <Route path="technicians" element={<TechnicianListPage />} />
          <Route path="suppliers" element={<SupplierListPage />} />

          {/* Communication */}
          <Route path="messaging/devices" element={<MessagingSettingsPage />} />
          <Route path="messaging/blast" element={<BlastMessagePage />} />

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
      </OrgProvider>
    </AuthProvider>
  );
}

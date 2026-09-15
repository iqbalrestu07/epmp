import { InteractiveExplorerPage } from "./features/immersive-view/pages/InteractiveExplorerPage";
import { Routes, Route, Navigate, useLocation } from "react-router-dom";
import { AuthProvider } from "./features/iam/context/AuthContext";
import { OrgProvider } from "./features/organization/context/OrgContext";
import ProtectedRoute from "./components/ProtectedRoute";
import { ErrorBoundary } from "./components/ErrorBoundary";

import MainLayout from "./layouts/MainLayout";
import DashboardPage from "./pages/DashboardPage";
import ImmersiveLanding from "./features/immersive-view/pages/ImmersiveLanding";

import AuthLayout from "./layouts/AuthLayout";
import SignInPage from "./features/iam/pages/SignInPage";
import SignUpPage from "./features/iam/pages/SignUpPage";

// Organization
import { OrganizationListPage } from "./features/organization/pages/OrganizationListPage";
import { OrganizationCreatePage } from "./features/organization/pages/OrganizationCreatePage";
import { OrganizationDetailPage } from "./features/organization/pages/OrganizationDetailPage";
import { OrganizationEditPage } from "./features/organization/pages/OrganizationEditPage";

// Property
import { PropertyListPage } from "./features/property/pages/PropertyListPage";
import { PropertyCreatePage } from "./features/property/pages/PropertyCreatePage";
import { PropertyDetailPage } from "./features/property/pages/PropertyDetailPage";
import { PropertyEditPage } from "./features/property/pages/PropertyEditPage";
import PropertyInteractiveView from "./features/property/pages/PropertyInteractiveView";

// Building
import { BuildingListPage } from "./features/building/pages/BuildingListPage";
import { BuildingCreatePage } from "./features/building/pages/BuildingCreatePage";
import { BuildingDetailPage } from "./features/building/pages/BuildingDetailPage";
import { BuildingEditPage } from "./features/building/pages/BuildingEditPage";

// Floor
import { FloorListPage } from "./features/floor/pages/FloorListPage";
import { FloorCreatePage } from "./features/floor/pages/FloorCreatePage";
import { FloorDetailPage } from "./features/floor/pages/FloorDetailPage";
import { FloorEditPage } from "./features/floor/pages/FloorEditPage";

// Zone / Room / Bed / Facility / RoomType
import { ZoneListPage } from "./features/zone/pages/ZoneListPage";
import { ZoneCreatePage } from "./features/zone/pages/ZoneCreatePage";
import { ZoneDetailPage } from "./features/zone/pages/ZoneDetailPage";
import { ZoneEditPage } from "./features/zone/pages/ZoneEditPage";
import { RoomListPage } from "./features/room/pages/RoomListPage";
import { RoomCreatePage } from "./features/room/pages/RoomCreatePage";
import { RoomDetailPage } from "./features/room/pages/RoomDetailPage";
import { RoomEditPage } from "./features/room/pages/RoomEditPage";
import { RoomTypeListPage } from "./features/roomtype/pages/RoomTypeListPage";
import { RoomTypeCreatePage } from "./features/roomtype/pages/RoomTypeCreatePage";
import { RoomTypeDetailPage } from "./features/roomtype/pages/RoomTypeDetailPage";
import { RoomTypeEditPage } from "./features/roomtype/pages/RoomTypeEditPage";
import { BedListPage } from "./features/bed/pages/BedListPage";
import { BedCreatePage } from "./features/bed/pages/BedCreatePage";
import { BedDetailPage } from "./features/bed/pages/BedDetailPage";
import { BedEditPage } from "./features/bed/pages/BedEditPage";
import { FacilityListPage } from "./features/facility/pages/FacilityListPage";
import { FacilityCreatePage } from "./features/facility/pages/FacilityCreatePage";
import { FacilityDetailPage } from "./features/facility/pages/FacilityDetailPage";
import { FacilityEditPage } from "./features/facility/pages/FacilityEditPage";

// Tenant
import { TenantListPage } from "./features/tenant/pages/TenantListPage";
import { TenantCreatePage } from "./features/tenant/pages/TenantCreatePage";
import { TenantDetailPage } from "./features/tenant/pages/TenantDetailPage";
import { TenantEditPage } from "./features/tenant/pages/TenantEditPage";

// Operations
import { ReservationListPage } from "./features/reservation/pages/ReservationListPage";
import { ReservationCreatePage } from "./features/reservation/pages/ReservationCreatePage";
import { ReservationDetailPage } from "./features/reservation/pages/ReservationDetailPage";
import { ReservationEditPage } from "./features/reservation/pages/ReservationEditPage";
import { ContractListPage } from "./features/contract/pages/ContractListPage";
import { ContractCreatePage } from "./features/contract/pages/ContractCreatePage";
import { ContractDetailPage } from "./features/contract/pages/ContractDetailPage";
import { ContractEditPage } from "./features/contract/pages/ContractEditPage";
import { OccupancyListPage } from "./features/occupancy/pages/OccupancyListPage";
import { OccupancyCreatePage } from "./features/occupancy/pages/OccupancyCreatePage";
import { OccupancyDetailPage } from "./features/occupancy/pages/OccupancyDetailPage";
import { OccupancyEditPage } from "./features/occupancy/pages/OccupancyEditPage";
import { TenantRosterPage } from "./features/occupancy/pages/TenantRosterPage";

// Finance
import { InvoiceListPage } from "./features/billing/pages/InvoiceListPage";
import { InvoiceCreatePage } from "./features/billing/pages/InvoiceCreatePage";
import { InvoiceDetailPage } from "./features/billing/pages/InvoiceDetailPage";
import { InvoiceEditPage } from "./features/billing/pages/InvoiceEditPage";
import { PaymentListPage } from "./features/payment/pages/PaymentListPage";
import { PaymentCreatePage } from "./features/payment/pages/PaymentCreatePage";
import { PaymentDetailPage } from "./features/payment/pages/PaymentDetailPage";
import { PaymentEditPage } from "./features/payment/pages/PaymentEditPage";
import { DepositListPage } from "./features/deposit/pages/DepositListPage";
import { DepositCreatePage } from "./features/deposit/pages/DepositCreatePage";
import { DepositDetailPage } from "./features/deposit/pages/DepositDetailPage";
import { DepositEditPage } from "./features/deposit/pages/DepositEditPage";
import { ChargeListPage } from "./features/charge/pages/ChargeListPage";
import { ChargeCreatePage } from "./features/charge/pages/ChargeCreatePage";
import { ChargeDetailPage } from "./features/charge/pages/ChargeDetailPage";
import { ChargeEditPage } from "./features/charge/pages/ChargeEditPage";
import { RefundListPage } from "./features/refund/pages/RefundListPage";
import { RefundCreatePage } from "./features/refund/pages/RefundCreatePage";
import { RefundDetailPage } from "./features/refund/pages/RefundDetailPage";
import { RefundEditPage } from "./features/refund/pages/RefundEditPage";
import { AdjustmentListPage } from "./features/adjustment/pages/AdjustmentListPage";
import { AdjustmentCreatePage } from "./features/adjustment/pages/AdjustmentCreatePage";
import { AdjustmentDetailPage } from "./features/adjustment/pages/AdjustmentDetailPage";
import { AdjustmentEditPage } from "./features/adjustment/pages/AdjustmentEditPage";
import { PenaltyListPage } from "./features/penalty/pages/PenaltyListPage";
import { PenaltyCreatePage } from "./features/penalty/pages/PenaltyCreatePage";
import { PenaltyDetailPage } from "./features/penalty/pages/PenaltyDetailPage";
import { PenaltyEditPage } from "./features/penalty/pages/PenaltyEditPage";

// Assets & Maintenance
import { AssetListPage } from "./features/asset/pages/AssetListPage";
import { AssetCreatePage } from "./features/asset/pages/AssetCreatePage";
import { AssetEditPage } from "./features/asset/pages/AssetEditPage";
import { AssetDetailPage } from "./features/asset/pages/AssetDetailPage";
import { AssetAssignmentListPage } from "./features/assetassignment/pages/AssetAssignmentListPage";
import { AssetAssignmentCreatePage } from "./features/assetassignment/pages/AssetAssignmentCreatePage";
import { AssetAssignmentDetailPage } from "./features/assetassignment/pages/AssetAssignmentDetailPage";
import { AssetAssignmentEditPage } from "./features/assetassignment/pages/AssetAssignmentEditPage";
import { AssetInspectionListPage } from "./features/assetinspection/pages/AssetInspectionListPage";
import { AssetInspectionCreatePage } from "./features/assetinspection/pages/AssetInspectionCreatePage";
import { AssetInspectionDetailPage } from "./features/assetinspection/pages/AssetInspectionDetailPage";
import { AssetInspectionEditPage } from "./features/assetinspection/pages/AssetInspectionEditPage";
import { WorkOrderListPage } from "./features/workorder/pages/WorkOrderListPage";
import { WorkOrderCreatePage } from "./features/workorder/pages/WorkOrderCreatePage";
import { WorkOrderDetailPage } from "./features/workorder/pages/WorkOrderDetailPage";
import { WorkOrderEditPage } from "./features/workorder/pages/WorkOrderEditPage";
import { TechnicianListPage } from "./features/technician/pages/TechnicianListPage";
import { TechnicianCreatePage } from "./features/technician/pages/TechnicianCreatePage";
import { TechnicianDetailPage } from "./features/technician/pages/TechnicianDetailPage";
import { TechnicianEditPage } from "./features/technician/pages/TechnicianEditPage";
import { SupplierListPage } from "./features/supplier/pages/SupplierListPage";
import { SupplierCreatePage } from "./features/supplier/pages/SupplierCreatePage";
import { SupplierDetailPage } from "./features/supplier/pages/SupplierDetailPage";
import { SupplierEditPage } from "./features/supplier/pages/SupplierEditPage";

// Communication
import { MessagingSettingsPage } from "./features/communication/pages/MessagingSettingsPage";
import { BlastMessagePage } from "./features/communication/pages/BlastMessagePage";

// IAM
import RBACPage from "./features/iam/pages/RBACPage";
import UserListPage from "./features/iam/pages/UserListPage";

// Settings
import { SettingsPage } from "./features/settings/pages/SettingsPage";
import { ReportsPage } from "./features/report/pages/ReportsPage";
import { AuditLogPage } from "./features/audit/pages/AuditLogPage";

function getModuleRoutes() {
  return (
    <>
      <Route path="overview" element={<DashboardPage />} />
      <Route path="explorer" element={<InteractiveExplorerPage />} />
      <Route path="roster" element={<TenantRosterPage />} />

      {/* Organizations */}
      <Route path="organizations" element={<OrganizationListPage />} />
      <Route path="organizations/new" element={<OrganizationCreatePage />} />
      <Route path="organizations/:id" element={<OrganizationDetailPage />} />
      <Route path="organizations/:id/edit" element={<OrganizationEditPage />} />

      {/* Properties */}
      <Route path="properties" element={<PropertyListPage />} />
      <Route path="properties/new" element={<PropertyCreatePage />} />
      <Route path="properties/interactive" element={<PropertyInteractiveView />} />
      <Route path="properties/:id" element={<PropertyDetailPage />} />
      <Route path="properties/:id/edit" element={<PropertyEditPage />} />

      {/* Buildings */}
      <Route path="buildings" element={<BuildingListPage />} />
      <Route path="buildings/new" element={<BuildingCreatePage />} />
      <Route path="buildings/:id" element={<BuildingDetailPage />} />
      <Route path="buildings/:id/edit" element={<BuildingEditPage />} />

      {/* Floors */}
      <Route path="floors" element={<FloorListPage />} />
      <Route path="floors/new" element={<FloorCreatePage />} />
      <Route path="floors/:id" element={<FloorDetailPage />} />
      <Route path="floors/:id/edit" element={<FloorEditPage />} />

      {/* Zones */}
      <Route path="zones" element={<ZoneListPage />} />
      <Route path="zones/new" element={<ZoneCreatePage />} />
      <Route path="zones/:id" element={<ZoneDetailPage />} />
      <Route path="zones/:id/edit" element={<ZoneEditPage />} />

      {/* Rooms */}
      <Route path="rooms" element={<RoomListPage />} />
      <Route path="rooms/new" element={<RoomCreatePage />} />
      <Route path="rooms/:id" element={<RoomDetailPage />} />
      <Route path="rooms/:id/edit" element={<RoomEditPage />} />

      {/* Room Types */}
      <Route path="room-types" element={<RoomTypeListPage />} />
      <Route path="room-types/new" element={<RoomTypeCreatePage />} />
      <Route path="room-types/:id" element={<RoomTypeDetailPage />} />
      <Route path="room-types/:id/edit" element={<RoomTypeEditPage />} />

      {/* Beds */}
      <Route path="beds" element={<BedListPage />} />
      <Route path="beds/new" element={<BedCreatePage />} />
      <Route path="beds/:id" element={<BedDetailPage />} />
      <Route path="beds/:id/edit" element={<BedEditPage />} />

      {/* Facilities */}
      <Route path="facilities" element={<FacilityListPage />} />
      <Route path="facilities/new" element={<FacilityCreatePage />} />
      <Route path="facilities/:id" element={<FacilityDetailPage />} />
      <Route path="facilities/:id/edit" element={<FacilityEditPage />} />

      {/* Tenants */}
      <Route path="tenants" element={<TenantListPage />} />
      <Route path="tenants/new" element={<TenantCreatePage />} />
      <Route path="tenants/:id" element={<TenantDetailPage />} />
      <Route path="tenants/:id/edit" element={<TenantEditPage />} />

      {/* Operations */}
      <Route path="reservations" element={<ReservationListPage />} />
      <Route path="reservations/new" element={<ReservationCreatePage />} />
      <Route path="reservations/:id" element={<ReservationDetailPage />} />
      <Route path="reservations/:id/edit" element={<ReservationEditPage />} />

      <Route path="contracts" element={<ContractListPage />} />
      <Route path="contracts/new" element={<ContractCreatePage />} />
      <Route path="contracts/:id" element={<ContractDetailPage />} />
      <Route path="contracts/:id/edit" element={<ContractEditPage />} />

      <Route path="occupancies" element={<OccupancyListPage />} />
      <Route path="occupancies/new" element={<OccupancyCreatePage />} />
      <Route path="occupancies/:id" element={<OccupancyDetailPage />} />
      <Route path="occupancies/:id/edit" element={<OccupancyEditPage />} />

      {/* Finance */}
      <Route path="invoices" element={<InvoiceListPage />} />
      <Route path="invoices/new" element={<InvoiceCreatePage />} />
      <Route path="invoices/:id" element={<InvoiceDetailPage />} />
      <Route path="invoices/:id/edit" element={<InvoiceEditPage />} />

      <Route path="payments" element={<PaymentListPage />} />
      <Route path="payments/new" element={<PaymentCreatePage />} />
      <Route path="payments/:id" element={<PaymentDetailPage />} />
      <Route path="payments/:id/edit" element={<PaymentEditPage />} />

      <Route path="deposits" element={<DepositListPage />} />
      <Route path="deposits/new" element={<DepositCreatePage />} />
      <Route path="deposits/:id" element={<DepositDetailPage />} />
      <Route path="deposits/:id/edit" element={<DepositEditPage />} />

      <Route path="charges" element={<ChargeListPage />} />
      <Route path="charges/new" element={<ChargeCreatePage />} />
      <Route path="charges/:id" element={<ChargeDetailPage />} />
      <Route path="charges/:id/edit" element={<ChargeEditPage />} />

      <Route path="refunds" element={<RefundListPage />} />
      <Route path="refunds/new" element={<RefundCreatePage />} />
      <Route path="refunds/:id" element={<RefundDetailPage />} />
      <Route path="refunds/:id/edit" element={<RefundEditPage />} />

      <Route path="adjustments" element={<AdjustmentListPage />} />
      <Route path="adjustments/new" element={<AdjustmentCreatePage />} />
      <Route path="adjustments/:id" element={<AdjustmentDetailPage />} />
      <Route path="adjustments/:id/edit" element={<AdjustmentEditPage />} />

      <Route path="penalties" element={<PenaltyListPage />} />
      <Route path="penalties/new" element={<PenaltyCreatePage />} />
      <Route path="penalties/:id" element={<PenaltyDetailPage />} />
      <Route path="penalties/:id/edit" element={<PenaltyEditPage />} />

      {/* Assets & Maintenance */}
      <Route path="assets" element={<AssetListPage />} />
      <Route path="assets/new" element={<AssetCreatePage />} />
      <Route path="assets/:id" element={<AssetDetailPage />} />
      <Route path="assets/:id/edit" element={<AssetEditPage />} />

      <Route path="asset-assignments" element={<AssetAssignmentListPage />} />
      <Route path="asset-assignments/new" element={<AssetAssignmentCreatePage />} />
      <Route path="asset-assignments/:id" element={<AssetAssignmentDetailPage />} />
      <Route path="asset-assignments/:id/edit" element={<AssetAssignmentEditPage />} />

      <Route path="asset-inspections" element={<AssetInspectionListPage />} />
      <Route path="asset-inspections/new" element={<AssetInspectionCreatePage />} />
      <Route path="asset-inspections/:id" element={<AssetInspectionDetailPage />} />
      <Route path="asset-inspections/:id/edit" element={<AssetInspectionEditPage />} />

      <Route path="work-orders" element={<WorkOrderListPage />} />
      <Route path="work-orders/new" element={<WorkOrderCreatePage />} />
      <Route path="work-orders/:id" element={<WorkOrderDetailPage />} />
      <Route path="work-orders/:id/edit" element={<WorkOrderEditPage />} />

      <Route path="technicians" element={<TechnicianListPage />} />
      <Route path="technicians/new" element={<TechnicianCreatePage />} />
      <Route path="technicians/:id" element={<TechnicianDetailPage />} />
      <Route path="technicians/:id/edit" element={<TechnicianEditPage />} />

      <Route path="suppliers" element={<SupplierListPage />} />
      <Route path="suppliers/new" element={<SupplierCreatePage />} />
      <Route path="suppliers/:id" element={<SupplierDetailPage />} />
      <Route path="suppliers/:id/edit" element={<SupplierEditPage />} />

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

      <Route path="reports" element={<ReportsPage />} />
      <Route path="audit-logs" element={<AuditLogPage />} />
      <Route path="settings" element={<SettingsPage />} />
    </>
  );
}

function AppRoutes() {
  const location = useLocation();
  return (
    <ErrorBoundary key={location.pathname}>
      <Routes>
          {/* Public Routes */}
          <Route path="/" element={<ImmersiveLanding />} />

          {/* Auth Routes */}
          <Route path="/auth" element={<AuthLayout />}>
            <Route path="signin" element={<SignInPage />} />
            <Route path="signup" element={<SignUpPage />} />
          </Route>

          {/* Clean App Routes: direct access via /properties, /rooms, /contracts, /roster, etc. */}
          <Route
            element={
              <ProtectedRoute>
                <MainLayout />
              </ProtectedRoute>
            }
          >
            <Route path="dashboard" element={<DashboardPage />} />
            {getModuleRoutes()}
          </Route>

          {/* Dashboard Compatibility Routes: /dashboard/properties, /dashboard/rooms, etc. */}
          <Route
            path="/dashboard"
            element={
              <ProtectedRoute>
                <MainLayout />
              </ProtectedRoute>
            }
          >
            <Route index element={<DashboardPage />} />
            {getModuleRoutes()}
            <Route path="*" element={<Navigate to="/dashboard" replace />} />
          </Route>

        {/* Catch-all */}
        <Route path="*" element={<Navigate to="/dashboard" replace />} />
      </Routes>
    </ErrorBoundary>
  );
}

export default function App() {
  return (
    <AuthProvider>
      <OrgProvider>
        <AppRoutes />
      </OrgProvider>
    </AuthProvider>
  );
}

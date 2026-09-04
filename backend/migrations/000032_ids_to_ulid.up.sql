-- Migration: 000033_ids_to_ulid
-- Purpose  : Change all primary key and foreign key columns from UUID to TEXT
--            so they can store ULID strings (26-char, e.g. "01ARZ3NDEKTSV4RRFFQ69G5FAV").
--            ULIDs are lexicographically sortable and time-ordered.
--            IDs are now generated in application code (Go) using uid.New().
-- Created  : 2026-09-04
--
-- NOTE: This uses ALTER COLUMN ... TYPE TEXT USING id::TEXT
--       to safely cast existing UUID values to their text representation.
--       All FK constraints must be dropped and recreated after type change.

-- ── Drop all FK constraints first (required before type change) ────────────

-- refresh_tokens
ALTER TABLE refresh_tokens DROP CONSTRAINT IF EXISTS refresh_tokens_user_id_fkey CASCADE;

-- user_roles
ALTER TABLE user_roles DROP CONSTRAINT IF EXISTS user_roles_pkey CASCADE;
ALTER TABLE user_roles DROP CONSTRAINT IF EXISTS user_roles_user_id_fkey CASCADE;
ALTER TABLE user_roles DROP CONSTRAINT IF EXISTS user_roles_role_id_fkey CASCADE;

-- role_permissions
ALTER TABLE role_permissions DROP CONSTRAINT IF EXISTS role_permissions_pkey CASCADE;
ALTER TABLE role_permissions DROP CONSTRAINT IF EXISTS role_permissions_role_id_fkey CASCADE;
ALTER TABLE role_permissions DROP CONSTRAINT IF EXISTS role_permissions_permission_id_fkey CASCADE;

-- organizations
ALTER TABLE organizations DROP CONSTRAINT IF EXISTS organizations_created_by_fkey CASCADE;

-- properties
ALTER TABLE properties DROP CONSTRAINT IF EXISTS properties_pkey CASCADE;
ALTER TABLE properties DROP CONSTRAINT IF EXISTS properties_organization_id_fkey CASCADE;

-- buildings
ALTER TABLE buildings DROP CONSTRAINT IF EXISTS buildings_pkey CASCADE;
ALTER TABLE buildings DROP CONSTRAINT IF EXISTS buildings_property_id_fkey CASCADE;
ALTER TABLE buildings DROP CONSTRAINT IF EXISTS buildings_organization_id_fkey CASCADE;

-- floors
ALTER TABLE floors DROP CONSTRAINT IF EXISTS floors_pkey CASCADE;
ALTER TABLE floors DROP CONSTRAINT IF EXISTS floors_building_id_fkey CASCADE;
ALTER TABLE floors DROP CONSTRAINT IF EXISTS floors_organization_id_fkey CASCADE;

-- zones
ALTER TABLE zones DROP CONSTRAINT IF EXISTS zones_pkey CASCADE;
ALTER TABLE zones DROP CONSTRAINT IF EXISTS zones_building_id_fkey CASCADE;
ALTER TABLE zones DROP CONSTRAINT IF EXISTS zones_organization_id_fkey CASCADE;

-- rooms
ALTER TABLE rooms DROP CONSTRAINT IF EXISTS rooms_pkey CASCADE;
ALTER TABLE rooms DROP CONSTRAINT IF EXISTS rooms_property_id_fkey CASCADE;
ALTER TABLE rooms DROP CONSTRAINT IF EXISTS rooms_floor_id_fkey CASCADE;
ALTER TABLE rooms DROP CONSTRAINT IF EXISTS rooms_organization_id_fkey CASCADE;

-- beds
ALTER TABLE beds DROP CONSTRAINT IF EXISTS beds_pkey CASCADE;
ALTER TABLE beds DROP CONSTRAINT IF EXISTS beds_room_id_fkey CASCADE;
ALTER TABLE beds DROP CONSTRAINT IF EXISTS beds_organization_id_fkey CASCADE;

-- room_types
ALTER TABLE room_types DROP CONSTRAINT IF EXISTS room_types_pkey CASCADE;
ALTER TABLE room_types DROP CONSTRAINT IF EXISTS room_types_organization_id_fkey CASCADE;

-- facilities
ALTER TABLE facilities DROP CONSTRAINT IF EXISTS facilities_pkey CASCADE;
ALTER TABLE facilities DROP CONSTRAINT IF EXISTS facilities_organization_id_fkey CASCADE;

-- tenants
ALTER TABLE tenants DROP CONSTRAINT IF EXISTS tenants_pkey CASCADE;
ALTER TABLE tenants DROP CONSTRAINT IF EXISTS tenants_organization_id_fkey CASCADE;

-- tenant_identities
ALTER TABLE tenant_identities DROP CONSTRAINT IF EXISTS tenant_identities_pkey CASCADE;
ALTER TABLE tenant_identities DROP CONSTRAINT IF EXISTS tenant_identities_tenant_id_fkey CASCADE;
ALTER TABLE tenant_identities DROP CONSTRAINT IF EXISTS tenant_identities_organization_id_fkey CASCADE;

-- tenant_contacts
ALTER TABLE tenant_contacts DROP CONSTRAINT IF EXISTS tenant_contacts_pkey CASCADE;
ALTER TABLE tenant_contacts DROP CONSTRAINT IF EXISTS tenant_contacts_tenant_id_fkey CASCADE;
ALTER TABLE tenant_contacts DROP CONSTRAINT IF EXISTS tenant_contacts_organization_id_fkey CASCADE;

-- tenant_documents
ALTER TABLE tenant_documents DROP CONSTRAINT IF EXISTS tenant_documents_pkey CASCADE;
ALTER TABLE tenant_documents DROP CONSTRAINT IF EXISTS tenant_documents_tenant_id_fkey CASCADE;
ALTER TABLE tenant_documents DROP CONSTRAINT IF EXISTS tenant_documents_organization_id_fkey CASCADE;

-- reservations
ALTER TABLE reservations DROP CONSTRAINT IF EXISTS reservations_pkey CASCADE;
ALTER TABLE reservations DROP CONSTRAINT IF EXISTS reservations_tenant_id_fkey CASCADE;
ALTER TABLE reservations DROP CONSTRAINT IF EXISTS reservations_property_id_fkey CASCADE;
ALTER TABLE reservations DROP CONSTRAINT IF EXISTS reservations_room_id_fkey CASCADE;
ALTER TABLE reservations DROP CONSTRAINT IF EXISTS reservations_organization_id_fkey CASCADE;

-- contracts
ALTER TABLE contracts DROP CONSTRAINT IF EXISTS contracts_pkey CASCADE;
ALTER TABLE contracts DROP CONSTRAINT IF EXISTS contracts_reservation_id_fkey CASCADE;
ALTER TABLE contracts DROP CONSTRAINT IF EXISTS contracts_tenant_id_fkey CASCADE;
ALTER TABLE contracts DROP CONSTRAINT IF EXISTS contracts_property_id_fkey CASCADE;
ALTER TABLE contracts DROP CONSTRAINT IF EXISTS contracts_room_id_fkey CASCADE;
ALTER TABLE contracts DROP CONSTRAINT IF EXISTS contracts_organization_id_fkey CASCADE;

-- occupancies
ALTER TABLE occupancies DROP CONSTRAINT IF EXISTS occupancies_pkey CASCADE;
ALTER TABLE occupancies DROP CONSTRAINT IF EXISTS occupancies_contract_id_fkey CASCADE;
ALTER TABLE occupancies DROP CONSTRAINT IF EXISTS occupancies_room_id_fkey CASCADE;
ALTER TABLE occupancies DROP CONSTRAINT IF EXISTS occupancies_tenant_id_fkey CASCADE;
ALTER TABLE occupancies DROP CONSTRAINT IF EXISTS occupancies_organization_id_fkey CASCADE;

-- invoices
ALTER TABLE invoices DROP CONSTRAINT IF EXISTS invoices_pkey CASCADE;
ALTER TABLE invoices DROP CONSTRAINT IF EXISTS invoices_contract_id_fkey CASCADE;
ALTER TABLE invoices DROP CONSTRAINT IF EXISTS invoices_tenant_id_fkey CASCADE;
ALTER TABLE invoices DROP CONSTRAINT IF EXISTS invoices_organization_id_fkey CASCADE;

-- payments
ALTER TABLE payments DROP CONSTRAINT IF EXISTS payments_pkey CASCADE;
ALTER TABLE payments DROP CONSTRAINT IF EXISTS payments_organization_id_fkey CASCADE;

-- deposits
ALTER TABLE deposits DROP CONSTRAINT IF EXISTS deposits_pkey CASCADE;
ALTER TABLE deposits DROP CONSTRAINT IF EXISTS deposits_organization_id_fkey CASCADE;

-- charges
ALTER TABLE charges DROP CONSTRAINT IF EXISTS charges_pkey CASCADE;
ALTER TABLE charges DROP CONSTRAINT IF EXISTS charges_organization_id_fkey CASCADE;

-- refunds
ALTER TABLE refunds DROP CONSTRAINT IF EXISTS refunds_pkey CASCADE;
ALTER TABLE refunds DROP CONSTRAINT IF EXISTS refunds_organization_id_fkey CASCADE;

-- adjustments
ALTER TABLE adjustments DROP CONSTRAINT IF EXISTS adjustments_pkey CASCADE;
ALTER TABLE adjustments DROP CONSTRAINT IF EXISTS adjustments_organization_id_fkey CASCADE;

-- penalties
ALTER TABLE penalties DROP CONSTRAINT IF EXISTS penalties_pkey CASCADE;
ALTER TABLE penalties DROP CONSTRAINT IF EXISTS penalties_organization_id_fkey CASCADE;

-- assets
ALTER TABLE assets DROP CONSTRAINT IF EXISTS assets_pkey CASCADE;
ALTER TABLE assets DROP CONSTRAINT IF EXISTS assets_property_id_fkey CASCADE;
ALTER TABLE assets DROP CONSTRAINT IF EXISTS assets_organization_id_fkey CASCADE;

-- asset_assignments
ALTER TABLE asset_assignments DROP CONSTRAINT IF EXISTS asset_assignments_pkey CASCADE;
ALTER TABLE asset_assignments DROP CONSTRAINT IF EXISTS asset_assignments_organization_id_fkey CASCADE;

-- asset_inspections
ALTER TABLE asset_inspections DROP CONSTRAINT IF EXISTS asset_inspections_pkey CASCADE;
ALTER TABLE asset_inspections DROP CONSTRAINT IF EXISTS asset_inspections_organization_id_fkey CASCADE;

-- work_orders
ALTER TABLE work_orders DROP CONSTRAINT IF EXISTS work_orders_pkey CASCADE;
ALTER TABLE work_orders DROP CONSTRAINT IF EXISTS work_orders_organization_id_fkey CASCADE;

-- technicians
ALTER TABLE technicians DROP CONSTRAINT IF EXISTS technicians_pkey CASCADE;
ALTER TABLE technicians DROP CONSTRAINT IF EXISTS technicians_organization_id_fkey CASCADE;

-- vendors
ALTER TABLE vendors DROP CONSTRAINT IF EXISTS vendors_pkey CASCADE;
ALTER TABLE vendors DROP CONSTRAINT IF EXISTS vendors_organization_id_fkey CASCADE;

-- users (last, others reference it)
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_pkey CASCADE;

-- roles
ALTER TABLE roles DROP CONSTRAINT IF EXISTS roles_pkey CASCADE;

-- permissions
ALTER TABLE permissions DROP CONSTRAINT IF EXISTS permissions_pkey CASCADE;

-- refresh_tokens
ALTER TABLE refresh_tokens DROP CONSTRAINT IF EXISTS refresh_tokens_pkey CASCADE;

-- organizations
ALTER TABLE organizations DROP CONSTRAINT IF EXISTS organizations_pkey CASCADE;

-- ── Change all PK and FK columns from UUID to TEXT ─────────────────────────

-- users
ALTER TABLE users ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE users ALTER COLUMN organization_id TYPE TEXT USING organization_id::TEXT;

-- roles
ALTER TABLE roles ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE roles ALTER COLUMN organization_id TYPE TEXT USING organization_id::TEXT;

-- permissions
ALTER TABLE permissions ALTER COLUMN id TYPE TEXT USING id::TEXT;

-- user_roles
ALTER TABLE user_roles ALTER COLUMN user_id TYPE TEXT USING user_id::TEXT;
ALTER TABLE user_roles ALTER COLUMN role_id TYPE TEXT USING role_id::TEXT;

-- role_permissions
ALTER TABLE role_permissions ALTER COLUMN role_id TYPE TEXT USING role_id::TEXT;
ALTER TABLE role_permissions ALTER COLUMN permission_id TYPE TEXT USING permission_id::TEXT;

-- refresh_tokens
ALTER TABLE refresh_tokens ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE refresh_tokens ALTER COLUMN user_id TYPE TEXT USING user_id::TEXT;

-- organizations
ALTER TABLE organizations ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE organizations ADD COLUMN IF NOT EXISTS created_by TEXT;
ALTER TABLE organizations ALTER COLUMN created_by TYPE TEXT USING created_by::TEXT;

-- properties
ALTER TABLE properties ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE properties ALTER COLUMN organization_id TYPE TEXT USING organization_id::TEXT;

-- buildings
ALTER TABLE buildings ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE buildings ALTER COLUMN property_id TYPE TEXT USING property_id::TEXT;
ALTER TABLE buildings ADD COLUMN IF NOT EXISTS organization_id TEXT;
ALTER TABLE buildings ALTER COLUMN organization_id TYPE TEXT USING organization_id::TEXT;

-- floors
ALTER TABLE floors ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE floors ALTER COLUMN building_id TYPE TEXT USING building_id::TEXT;
ALTER TABLE floors ALTER COLUMN organization_id TYPE TEXT USING organization_id::TEXT;

-- zones
ALTER TABLE zones ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE zones ALTER COLUMN building_id TYPE TEXT USING building_id::TEXT;
ALTER TABLE zones ALTER COLUMN organization_id TYPE TEXT USING organization_id::TEXT;

-- rooms
ALTER TABLE rooms ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE rooms ALTER COLUMN property_id TYPE TEXT USING property_id::TEXT;
ALTER TABLE rooms ADD COLUMN IF NOT EXISTS floor_id TEXT;
ALTER TABLE rooms ALTER COLUMN floor_id TYPE TEXT USING floor_id::TEXT;
ALTER TABLE rooms ALTER COLUMN organization_id TYPE TEXT USING organization_id::TEXT;

-- beds
ALTER TABLE beds ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE beds ALTER COLUMN room_id TYPE TEXT USING room_id::TEXT;
ALTER TABLE beds ALTER COLUMN organization_id TYPE TEXT USING organization_id::TEXT;

-- room_types
ALTER TABLE room_types ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE room_types ALTER COLUMN organization_id TYPE TEXT USING organization_id::TEXT;

-- facilities
ALTER TABLE facilities ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE facilities ALTER COLUMN organization_id TYPE TEXT USING organization_id::TEXT;

-- tenants
ALTER TABLE tenants ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE tenants ALTER COLUMN organization_id TYPE TEXT USING organization_id::TEXT;

-- tenant_identities
ALTER TABLE tenant_identities ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE tenant_identities ALTER COLUMN tenant_id TYPE TEXT USING tenant_id::TEXT;
ALTER TABLE tenant_identities ALTER COLUMN organization_id TYPE TEXT USING organization_id::TEXT;

-- tenant_contacts
ALTER TABLE tenant_contacts ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE tenant_contacts ALTER COLUMN tenant_id TYPE TEXT USING tenant_id::TEXT;
ALTER TABLE tenant_contacts ALTER COLUMN organization_id TYPE TEXT USING organization_id::TEXT;

-- tenant_documents
ALTER TABLE tenant_documents ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE tenant_documents ALTER COLUMN tenant_id TYPE TEXT USING tenant_id::TEXT;
ALTER TABLE tenant_documents ALTER COLUMN organization_id TYPE TEXT USING organization_id::TEXT;

-- reservations
ALTER TABLE reservations ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE reservations ALTER COLUMN tenant_id TYPE TEXT USING tenant_id::TEXT;
ALTER TABLE reservations ALTER COLUMN property_id TYPE TEXT USING property_id::TEXT;
ALTER TABLE reservations ALTER COLUMN room_id TYPE TEXT USING room_id::TEXT;
ALTER TABLE reservations ALTER COLUMN organization_id TYPE TEXT USING organization_id::TEXT;

-- contracts
ALTER TABLE contracts ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE contracts ALTER COLUMN reservation_id TYPE TEXT USING reservation_id::TEXT;
ALTER TABLE contracts ALTER COLUMN tenant_id TYPE TEXT USING tenant_id::TEXT;
ALTER TABLE contracts ALTER COLUMN property_id TYPE TEXT USING property_id::TEXT;
ALTER TABLE contracts ALTER COLUMN room_id TYPE TEXT USING room_id::TEXT;
ALTER TABLE contracts ALTER COLUMN organization_id TYPE TEXT USING organization_id::TEXT;

-- occupancies
ALTER TABLE occupancies ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE occupancies ALTER COLUMN contract_id TYPE TEXT USING contract_id::TEXT;
ALTER TABLE occupancies ALTER COLUMN room_id TYPE TEXT USING room_id::TEXT;
ALTER TABLE occupancies ALTER COLUMN tenant_id TYPE TEXT USING tenant_id::TEXT;
ALTER TABLE occupancies ALTER COLUMN organization_id TYPE TEXT USING organization_id::TEXT;

-- invoices
ALTER TABLE invoices ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE invoices ALTER COLUMN contract_id TYPE TEXT USING contract_id::TEXT;
ALTER TABLE invoices ALTER COLUMN tenant_id TYPE TEXT USING tenant_id::TEXT;
ALTER TABLE invoices ALTER COLUMN organization_id TYPE TEXT USING organization_id::TEXT;

-- payments
ALTER TABLE payments ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE payments ALTER COLUMN organization_id TYPE TEXT USING organization_id::TEXT;

-- deposits
ALTER TABLE deposits ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE deposits ALTER COLUMN organization_id TYPE TEXT USING organization_id::TEXT;

-- charges
ALTER TABLE charges ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE charges ALTER COLUMN organization_id TYPE TEXT USING organization_id::TEXT;

-- refunds
ALTER TABLE refunds ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE refunds ALTER COLUMN organization_id TYPE TEXT USING organization_id::TEXT;

-- adjustments
ALTER TABLE adjustments ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE adjustments ALTER COLUMN organization_id TYPE TEXT USING organization_id::TEXT;

-- penalties
ALTER TABLE penalties ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE penalties ALTER COLUMN organization_id TYPE TEXT USING organization_id::TEXT;

-- assets
ALTER TABLE assets ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE assets ALTER COLUMN property_id TYPE TEXT USING property_id::TEXT;
ALTER TABLE assets ALTER COLUMN organization_id TYPE TEXT USING organization_id::TEXT;

-- asset_assignments
ALTER TABLE asset_assignments ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE asset_assignments ALTER COLUMN organization_id TYPE TEXT USING organization_id::TEXT;

-- asset_inspections
ALTER TABLE asset_inspections ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE asset_inspections ALTER COLUMN organization_id TYPE TEXT USING organization_id::TEXT;

-- work_orders
ALTER TABLE work_orders ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE work_orders ALTER COLUMN organization_id TYPE TEXT USING organization_id::TEXT;

-- technicians
ALTER TABLE technicians ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE technicians ALTER COLUMN organization_id TYPE TEXT USING organization_id::TEXT;

-- vendors
ALTER TABLE vendors ALTER COLUMN id TYPE TEXT USING id::TEXT;
ALTER TABLE vendors ALTER COLUMN organization_id TYPE TEXT USING organization_id::TEXT;

-- ── Re-add PRIMARY KEY constraints ─────────────────────────────────────────
ALTER TABLE users              ADD PRIMARY KEY (id);
ALTER TABLE roles              ADD PRIMARY KEY (id);
ALTER TABLE permissions        ADD PRIMARY KEY (id);
ALTER TABLE refresh_tokens     ADD PRIMARY KEY (id);
ALTER TABLE organizations      ADD PRIMARY KEY (id);
ALTER TABLE properties         ADD PRIMARY KEY (id);
ALTER TABLE buildings          ADD PRIMARY KEY (id);
ALTER TABLE floors             ADD PRIMARY KEY (id);
ALTER TABLE zones              ADD PRIMARY KEY (id);
ALTER TABLE rooms              ADD PRIMARY KEY (id);
ALTER TABLE beds               ADD PRIMARY KEY (id);
ALTER TABLE room_types         ADD PRIMARY KEY (id);
ALTER TABLE facilities         ADD PRIMARY KEY (id);
ALTER TABLE tenants            ADD PRIMARY KEY (id);
ALTER TABLE tenant_identities  ADD PRIMARY KEY (id);
ALTER TABLE tenant_contacts    ADD PRIMARY KEY (id);
ALTER TABLE tenant_documents   ADD PRIMARY KEY (id);
ALTER TABLE reservations       ADD PRIMARY KEY (id);
ALTER TABLE contracts          ADD PRIMARY KEY (id);
ALTER TABLE occupancies        ADD PRIMARY KEY (id);
ALTER TABLE invoices           ADD PRIMARY KEY (id);
ALTER TABLE payments           ADD PRIMARY KEY (id);
ALTER TABLE deposits           ADD PRIMARY KEY (id);
ALTER TABLE charges            ADD PRIMARY KEY (id);
ALTER TABLE refunds            ADD PRIMARY KEY (id);
ALTER TABLE adjustments        ADD PRIMARY KEY (id);
ALTER TABLE penalties          ADD PRIMARY KEY (id);
ALTER TABLE assets             ADD PRIMARY KEY (id);
ALTER TABLE asset_assignments  ADD PRIMARY KEY (id);
ALTER TABLE asset_inspections  ADD PRIMARY KEY (id);
ALTER TABLE work_orders        ADD PRIMARY KEY (id);
ALTER TABLE technicians        ADD PRIMARY KEY (id);
ALTER TABLE vendors            ADD PRIMARY KEY (id);

-- user_roles composite PK
ALTER TABLE user_roles         ADD PRIMARY KEY (user_id, role_id);
-- role_permissions composite PK
ALTER TABLE role_permissions   ADD PRIMARY KEY (role_id, permission_id);

-- ── Re-add FOREIGN KEY constraints ─────────────────────────────────────────
ALTER TABLE refresh_tokens       ADD CONSTRAINT fk_rt_user        FOREIGN KEY (user_id)         REFERENCES users(id)          ON DELETE CASCADE;
ALTER TABLE user_roles           ADD CONSTRAINT fk_ur_user        FOREIGN KEY (user_id)         REFERENCES users(id)          ON DELETE CASCADE;
ALTER TABLE user_roles           ADD CONSTRAINT fk_ur_role        FOREIGN KEY (role_id)         REFERENCES roles(id)          ON DELETE CASCADE;
ALTER TABLE role_permissions     ADD CONSTRAINT fk_rp_role        FOREIGN KEY (role_id)         REFERENCES roles(id)          ON DELETE CASCADE;
ALTER TABLE role_permissions     ADD CONSTRAINT fk_rp_perm        FOREIGN KEY (permission_id)   REFERENCES permissions(id)    ON DELETE CASCADE;
ALTER TABLE organizations        ADD CONSTRAINT fk_org_created_by FOREIGN KEY (created_by)      REFERENCES users(id);
ALTER TABLE properties           ADD CONSTRAINT fk_prop_org       FOREIGN KEY (organization_id) REFERENCES organizations(id);
ALTER TABLE buildings            ADD CONSTRAINT fk_bldg_prop      FOREIGN KEY (property_id)     REFERENCES properties(id);
ALTER TABLE buildings            ADD CONSTRAINT fk_bldg_org       FOREIGN KEY (organization_id) REFERENCES organizations(id);
ALTER TABLE floors               ADD CONSTRAINT fk_floor_bldg     FOREIGN KEY (building_id)     REFERENCES buildings(id);
ALTER TABLE floors               ADD CONSTRAINT fk_floor_org      FOREIGN KEY (organization_id) REFERENCES organizations(id);
ALTER TABLE zones                ADD CONSTRAINT fk_zone_bldg      FOREIGN KEY (building_id)     REFERENCES buildings(id);
ALTER TABLE zones                ADD CONSTRAINT fk_zone_org       FOREIGN KEY (organization_id) REFERENCES organizations(id);
ALTER TABLE rooms                ADD CONSTRAINT fk_room_prop      FOREIGN KEY (property_id)     REFERENCES properties(id);
ALTER TABLE rooms                ADD CONSTRAINT fk_room_floor     FOREIGN KEY (floor_id)        REFERENCES floors(id);
ALTER TABLE rooms                ADD CONSTRAINT fk_room_org       FOREIGN KEY (organization_id) REFERENCES organizations(id);
ALTER TABLE beds                 ADD CONSTRAINT fk_bed_room       FOREIGN KEY (room_id)         REFERENCES rooms(id);
ALTER TABLE beds                 ADD CONSTRAINT fk_bed_org        FOREIGN KEY (organization_id) REFERENCES organizations(id);
ALTER TABLE room_types           ADD CONSTRAINT fk_rt_org         FOREIGN KEY (organization_id) REFERENCES organizations(id);
ALTER TABLE facilities           ADD CONSTRAINT fk_fac_org        FOREIGN KEY (organization_id) REFERENCES organizations(id);
ALTER TABLE tenants              ADD CONSTRAINT fk_ten_org        FOREIGN KEY (organization_id) REFERENCES organizations(id);
ALTER TABLE tenant_identities    ADD CONSTRAINT fk_ti_ten         FOREIGN KEY (tenant_id)       REFERENCES tenants(id);
ALTER TABLE tenant_identities    ADD CONSTRAINT fk_ti_org         FOREIGN KEY (organization_id) REFERENCES organizations(id);
ALTER TABLE tenant_contacts      ADD CONSTRAINT fk_tc_ten         FOREIGN KEY (tenant_id)       REFERENCES tenants(id);
ALTER TABLE tenant_contacts      ADD CONSTRAINT fk_tc_org         FOREIGN KEY (organization_id) REFERENCES organizations(id);
ALTER TABLE tenant_documents     ADD CONSTRAINT fk_td_ten         FOREIGN KEY (tenant_id)       REFERENCES tenants(id);
ALTER TABLE tenant_documents     ADD CONSTRAINT fk_td_org         FOREIGN KEY (organization_id) REFERENCES organizations(id);
ALTER TABLE reservations         ADD CONSTRAINT fk_res_ten        FOREIGN KEY (tenant_id)       REFERENCES tenants(id);
ALTER TABLE reservations         ADD CONSTRAINT fk_res_prop       FOREIGN KEY (property_id)     REFERENCES properties(id);
ALTER TABLE reservations         ADD CONSTRAINT fk_res_room       FOREIGN KEY (room_id)         REFERENCES rooms(id);
ALTER TABLE reservations         ADD CONSTRAINT fk_res_org        FOREIGN KEY (organization_id) REFERENCES organizations(id);
ALTER TABLE contracts            ADD CONSTRAINT fk_con_res        FOREIGN KEY (reservation_id)  REFERENCES reservations(id);
ALTER TABLE contracts            ADD CONSTRAINT fk_con_ten        FOREIGN KEY (tenant_id)       REFERENCES tenants(id);
ALTER TABLE contracts            ADD CONSTRAINT fk_con_prop       FOREIGN KEY (property_id)     REFERENCES properties(id);
ALTER TABLE contracts            ADD CONSTRAINT fk_con_room       FOREIGN KEY (room_id)         REFERENCES rooms(id);
ALTER TABLE contracts            ADD CONSTRAINT fk_con_org        FOREIGN KEY (organization_id) REFERENCES organizations(id);
ALTER TABLE occupancies          ADD CONSTRAINT fk_occ_con        FOREIGN KEY (contract_id)     REFERENCES contracts(id);
ALTER TABLE occupancies          ADD CONSTRAINT fk_occ_room       FOREIGN KEY (room_id)         REFERENCES rooms(id);
ALTER TABLE occupancies          ADD CONSTRAINT fk_occ_ten        FOREIGN KEY (tenant_id)       REFERENCES tenants(id);
ALTER TABLE occupancies          ADD CONSTRAINT fk_occ_org        FOREIGN KEY (organization_id) REFERENCES organizations(id);
ALTER TABLE invoices             ADD CONSTRAINT fk_inv_con        FOREIGN KEY (contract_id)     REFERENCES contracts(id);
ALTER TABLE invoices             ADD CONSTRAINT fk_inv_ten        FOREIGN KEY (tenant_id)       REFERENCES tenants(id);
ALTER TABLE invoices             ADD CONSTRAINT fk_inv_org        FOREIGN KEY (organization_id) REFERENCES organizations(id);
ALTER TABLE payments             ADD CONSTRAINT fk_pay_org        FOREIGN KEY (organization_id) REFERENCES organizations(id);
ALTER TABLE deposits             ADD CONSTRAINT fk_dep_org        FOREIGN KEY (organization_id) REFERENCES organizations(id);
ALTER TABLE charges              ADD CONSTRAINT fk_chg_org        FOREIGN KEY (organization_id) REFERENCES organizations(id);
ALTER TABLE refunds              ADD CONSTRAINT fk_ref_org        FOREIGN KEY (organization_id) REFERENCES organizations(id);
ALTER TABLE adjustments          ADD CONSTRAINT fk_adj_org        FOREIGN KEY (organization_id) REFERENCES organizations(id);
ALTER TABLE penalties            ADD CONSTRAINT fk_pen_org        FOREIGN KEY (organization_id) REFERENCES organizations(id);
ALTER TABLE assets               ADD CONSTRAINT fk_ast_prop       FOREIGN KEY (property_id)     REFERENCES properties(id);
ALTER TABLE assets               ADD CONSTRAINT fk_ast_org        FOREIGN KEY (organization_id) REFERENCES organizations(id);
ALTER TABLE asset_assignments    ADD CONSTRAINT fk_aa_org         FOREIGN KEY (organization_id) REFERENCES organizations(id);
ALTER TABLE asset_inspections    ADD CONSTRAINT fk_ai_org         FOREIGN KEY (organization_id) REFERENCES organizations(id);
ALTER TABLE work_orders          ADD CONSTRAINT fk_wo_org         FOREIGN KEY (organization_id) REFERENCES organizations(id);
ALTER TABLE technicians          ADD CONSTRAINT fk_tech_org       FOREIGN KEY (organization_id) REFERENCES organizations(id);
ALTER TABLE vendors              ADD CONSTRAINT fk_vend_org       FOREIGN KEY (organization_id) REFERENCES organizations(id);

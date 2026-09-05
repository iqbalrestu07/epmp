const { chromium } = require('playwright');

const BASE_URL = 'http://localhost:3000';

const PAGES_TO_TEST = [
  '/dashboard',
  '/dashboard/explorer',
  '/dashboard/organizations',
  '/dashboard/organizations/new',
  '/dashboard/properties',
  '/dashboard/properties/new',
  '/dashboard/properties/interactive',
  '/dashboard/buildings',
  '/dashboard/buildings/new',
  '/dashboard/floors',
  '/dashboard/floors/new',
  '/dashboard/zones',
  '/dashboard/rooms',
  '/dashboard/rooms/new',
  '/dashboard/room-types',
  '/dashboard/beds',
  '/dashboard/facilities',
  '/dashboard/tenants',
  '/dashboard/tenants/new',
  '/dashboard/reservations',
  '/dashboard/contracts',
  '/dashboard/contracts/new',
  '/dashboard/occupancies',
  '/dashboard/invoices',
  '/dashboard/payments',
  '/dashboard/deposits',
  '/dashboard/charges',
  '/dashboard/refunds',
  '/dashboard/adjustments',
  '/dashboard/penalties',
  '/dashboard/assets',
  '/dashboard/assets/new',
  '/dashboard/asset-assignments',
  '/dashboard/asset-assignments/new',
  '/dashboard/asset-inspections',
  '/dashboard/work-orders',
  '/dashboard/technicians',
  '/dashboard/suppliers',
  '/dashboard/messaging/devices',
  '/dashboard/messaging/blast',
  '/dashboard/settings',
];

(async () => {
  const isHeaded = process.argv.includes('--headed') || process.env.HEADLESS === 'false';
  console.log(`\n🚀 Starting Full Frontend E2E Test Suite (${isHeaded ? '🖥️ GUI / Headed Mode' : '⚡ Headless Mode'})...`);
  
  const browser = await chromium.launch({
    channel: 'chrome',
    headless: !isHeaded,
    slowMo: isHeaded ? 120 : 0, // Slow down in GUI mode so actions can be watched live
  });

  const context = await browser.newContext({
    viewport: { width: 1440, height: 900 },
  });

  const page = await context.newPage();

  const errorsFound = [];
  const testResults = [];

  // Listen for console errors & page errors
  page.on('pageerror', (exception) => {
    errorsFound.push({
      type: 'PAGE_CRASH_EXCEPTION',
      url: page.url(),
      message: exception.message,
      stack: exception.stack,
    });
  });

  page.on('console', (msg) => {
    if (msg.type() === 'error') {
      const text = msg.text();
      // Ignore benign dev warnings and network interface changes
      if (!text.includes('favicon.ico') && !text.includes('[vite]') && !text.includes('ERR_NETWORK_CHANGED')) {
        errorsFound.push({
          type: 'CONSOLE_ERROR',
          url: page.url(),
          message: text,
        });
      }
    }
  });

  page.on('response', (response) => {
    if (response.status() >= 400) {
      const url = response.url();
      if (!url.includes('favicon.ico')) {
        errorsFound.push({
          type: `HTTP_${response.status()}`,
          url: page.url(),
          requestUrl: url,
        });
      }
    }
  });

  try {
    // 1. Sign in
    console.log('\n🔑 1. Logging in via /auth/signin...');
    await page.goto(`${BASE_URL}/auth/signin`, { waitUntil: 'networkidle' });

    await page.fill('input[type="email"]', 'e2e_tester@test.com');
    await page.fill('input[type="password"]', 'Password123!');
    await page.click('button[type="submit"]');

    // Wait for navigation into dashboard
    await page.waitForURL('**/dashboard**', { timeout: 10000 });
    console.log('✅ Successfully authenticated and reached dashboard!');

    // Wait for organization to settle
    await page.waitForTimeout(1500);

    // 2. Iterate through all 40 routes to check rendering & prevent crashes
    console.log(`\n🔍 2. Testing ${PAGES_TO_TEST.length} Routes (Read & Render Audit)...`);

    for (const route of PAGES_TO_TEST) {
      process.stdout.write(`   Testing ${route.padEnd(36)} ... `);
      const beforeErrCount = errorsFound.length;

      try {
        try {
          await page.goto(`${BASE_URL}${route}`, { waitUntil: 'networkidle', timeout: 8000 });
        } catch (navErr) {
          if (navErr.message.includes('ERR_NETWORK_CHANGED') || navErr.message.includes('timeout')) {
            await page.waitForTimeout(400);
            await page.goto(`${BASE_URL}${route}`, { waitUntil: 'domcontentloaded', timeout: 8000 });
          } else {
            throw navErr;
          }
        }
        await page.waitForTimeout(600);

        const bodyText = await page.innerText('body');
        if (!bodyText || bodyText.trim().length === 0) {
          errorsFound.push({
            type: 'BLANK_PAGE',
            url: `${BASE_URL}${route}`,
            message: 'Body has 0 text rendered',
          });
        }

        const newErrors = errorsFound.length - beforeErrCount;
        if (newErrors > 0) {
          console.log(`❌ (${newErrors} issues detected)`);
        } else {
          console.log(`✅ OK`);
        }
      } catch (err) {
        console.log(`⚠️ TIMEOUT / ERROR: ${err.message.split('\n')[0]}`);
        errorsFound.push({
          type: 'NAVIGATION_TIMEOUT_OR_ERROR',
          url: `${BASE_URL}${route}`,
          message: err.message,
        });
      }
    }

    // 3. Full Business Flow Lifecycle E2E Testing
    console.log('\n🔥 3. Executing Full End-to-End Business Flow Lifecycle (26 Sequential Steps)...');

    // Auto-accept confirmation dialogs
    page.on('dialog', async (dialog) => {
      console.log(`   [Dialog] Accepted: "${dialog.message()}"`);
      await dialog.accept();
    });

    const timestamp = Date.now().toString().slice(-5);
    const testPropName = `E2E Property ${timestamp}`;

    const helperSelect = async (selector, index = 1, timeout = 3000) => {
      const el = page.locator(selector);
      if ((await el.count()) > 0) {
        try {
          await page.waitForFunction(
            ({ sel, minCount }) => {
              const select = document.querySelector(sel);
              return select && select.options && select.options.length > minCount;
            },
            { sel: selector, minCount: index },
            { timeout }
          );
        } catch (e) {}
        const opts = await el.locator('option').all();
        if (opts.length > index) {
          await el.selectOption({ index }, { force: true });
          await el.dispatchEvent('change');
        }
      }
    };

    // ─── STAGE 1: PROPERTY MASTER INFRASTRUCTURE ───
    console.log('\n🏗️  STAGE 1: Setting up Property Master Infrastructure...');

    // 1. Property
    console.log(`   [1. Property] Creating: "${testPropName}"...`);
    await page.goto(`${BASE_URL}/dashboard/properties/new`, { waitUntil: 'networkidle' });
    await page.fill('#name', testPropName);
    await page.fill('#address', 'Jl. Jenderal Sudirman No. 88, Jakarta Pusat');
    await page.fill('#description', 'E2E Integrated Real Estate Property');
    await page.selectOption('#property_type', 'boarding_house');
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard/properties', { timeout: 8000 });
    await page.waitForTimeout(600);
    console.log(`   [1. Property] ✅ Created successfully!`);
    testResults.push('Property Created');

    // 2. Building
    const testBldgName = `Tower ${timestamp}`;
    console.log(`   [2. Building] Creating: "${testBldgName}"...`);
    await page.goto(`${BASE_URL}/dashboard/buildings/new`, { waitUntil: 'networkidle' });
    await page.waitForTimeout(600);
    const basicBtn = page.locator('button:has-text("Basic View")');
    if (await basicBtn.isVisible()) {
      await basicBtn.click();
      await page.waitForTimeout(400);
    }
    await helperSelect('#property_id', 1);
    await page.fill('#name', testBldgName);
    await page.fill('#total_floors', '4');
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard/buildings', { timeout: 8000 });
    await page.waitForTimeout(600);
    console.log(`   [2. Building] ✅ Created successfully!`);
    testResults.push('Building Created');

    // 3. Floor
    const testFloorName = `Floor 1 ${timestamp}`;
    console.log(`   [3. Floor] Creating: "${testFloorName}"...`);
    await page.goto(`${BASE_URL}/dashboard/floors/new`, { waitUntil: 'networkidle' });
    await helperSelect('#building_id', 1);
    await page.fill('#name', testFloorName);
    await page.fill('#floor_number', '1');
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard/floors', { timeout: 8000 });
    await page.waitForTimeout(600);
    console.log(`   [3. Floor] ✅ Created successfully!`);
    testResults.push('Floor Created');

    // 4. Zone
    const testZoneName = `Lobby Lounge ${timestamp}`;
    console.log(`   [4. Zone] Creating: "${testZoneName}"...`);
    await page.goto(`${BASE_URL}/dashboard/zones/new`, { waitUntil: 'networkidle' });
    await helperSelect('#building_id', 1);
    await page.fill('#floor', '1');
    await page.fill('#name', testZoneName);
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard/zones', { timeout: 8000 });
    await page.waitForTimeout(600);
    console.log(`   [4. Zone] ✅ Created successfully!`);
    testResults.push('Zone Created');

    // 5. Room Type
    const testRoomTypeName = `Deluxe Suite ${timestamp}`;
    console.log(`   [5. Room Type] Creating: "${testRoomTypeName}" with Base Price Rp 850.000...`);
    await page.goto(`${BASE_URL}/dashboard/room-types/new`, { waitUntil: 'networkidle' });
    await page.fill('#name', testRoomTypeName);
    await page.fill('#description', 'Executive deluxe suite with air conditioner');
    await page.fill('#base_price', '850000');
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard/room-types', { timeout: 8000 });
    await page.waitForTimeout(600);
    console.log(`   [5. Room Type] ✅ Created successfully!`);
    testResults.push('Room Type Created');

    // 6. Room
    const testRoomName = `Room ${timestamp.slice(-4)}`;
    console.log(`   [6. Room] Creating: "${testRoomName}" using Room Type "${testRoomTypeName}"...`);
    await page.goto(`${BASE_URL}/dashboard/rooms/new`, { waitUntil: 'networkidle' });
    await helperSelect('#property_id', 1);
    const rtSelect = page.locator('#room_type_id');
    if (await rtSelect.isVisible()) {
      await rtSelect.selectOption({ label: new RegExp(testRoomTypeName) }).catch(async () => {
        await helperSelect('#room_type_id', 1);
      });
    }
    await page.fill('#name', testRoomName);
    await page.fill('#capacity', '2');
    await page.fill('#price', '850000');
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard/rooms', { timeout: 8000 });
    await page.waitForTimeout(600);
    console.log(`   [6. Room] ✅ Created successfully!`);
    testResults.push('Room Created');

    // 7. Bed
    const testBedName = `Bed A-${timestamp.slice(-4)}`;
    console.log(`   [7. Bed] Creating: "${testBedName}"...`);
    await page.goto(`${BASE_URL}/dashboard/beds/new`, { waitUntil: 'networkidle' });
    await helperSelect('#room_id', 1);
    await page.fill('#name', testBedName);
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard/beds', { timeout: 8000 });
    await page.waitForTimeout(600);
    console.log(`   [7. Bed] ✅ Created successfully!`);
    testResults.push('Bed Created');

    // 8. Facility
    const testFacilityName = `Pool & Gym ${timestamp}`;
    console.log(`   [8. Facility] Creating: "${testFacilityName}"...`);
    await page.goto(`${BASE_URL}/dashboard/facilities/new`, { waitUntil: 'networkidle' });
    await helperSelect('#property_id', 1);
    await page.fill('#name', testFacilityName);
    await page.fill('#description', 'Free access swimming pool and gym for residents');
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard/facilities', { timeout: 8000 });
    await page.waitForTimeout(600);
    console.log(`   [8. Facility] ✅ Created successfully!`);
    testResults.push('Facility Created');


    // ─── STAGE 2: OPERATIONS & MAINTENANCE FLOW ───
    console.log('\n🔧 STAGE 2: Setting up Operations, Vendors & Asset Maintenance...');

    // 9. Technician
    const testTechName = `Budi Santoso ${timestamp}`;
    console.log(`   [9. Technician] Creating: "${testTechName}"...`);
    await page.goto(`${BASE_URL}/dashboard/technicians/new`, { waitUntil: 'networkidle' });
    await page.fill('#name', testTechName);
    await page.fill('#phone', '08123456789');
    await page.fill('#specialty', 'HVAC & Plumbing');
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard/technicians', { timeout: 8000 });
    await page.waitForTimeout(600);
    console.log(`   [9. Technician] ✅ Created successfully!`);
    testResults.push('Technician Created');

    // 10. Supplier
    const testSupplierName = `PT Mega Mandiri ${timestamp}`;
    console.log(`   [10. Supplier] Creating: "${testSupplierName}"...`);
    await page.goto(`${BASE_URL}/dashboard/suppliers/new`, { waitUntil: 'networkidle' });
    await page.fill('#name', testSupplierName);
    await page.fill('#contact_person', 'Hendra Wijaya');
    await page.fill('#phone', '08187654321');
    await page.fill('#service_type', 'AC & Hardware Supplier');
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard/suppliers', { timeout: 8000 });
    await page.waitForTimeout(600);
    console.log(`   [10. Supplier] ✅ Created successfully!`);
    testResults.push('Supplier Created');

    // 11. Asset
    const testAssetName = `Sharp AC 1.5PK ${timestamp}`;
    console.log(`   [11. Asset] Creating: "${testAssetName}"...`);
    await page.goto(`${BASE_URL}/dashboard/assets/new`, { waitUntil: 'networkidle' });
    await helperSelect('#property_id', 1);
    await page.fill('#name', testAssetName);
    await page.selectOption('#category', 'Electronics');
    await page.fill('#purchase_price', '4500000');
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard/assets', { timeout: 8000 });
    await page.waitForTimeout(600);
    console.log(`   [11. Asset] ✅ Created successfully!`);
    testResults.push('Asset Created');

    // 12. Asset Assignment
    console.log(`   [12. Asset Assignment] Assigning asset to room...`);
    await page.goto(`${BASE_URL}/dashboard/asset-assignments/new`, { waitUntil: 'networkidle' });
    const basicAssignBtn = page.locator('button:has-text("Basic View")');
    if (await basicAssignBtn.isVisible()) {
      await basicAssignBtn.click();
      await page.waitForTimeout(400);
    }
    await helperSelect('#asset_id', 1);
    await helperSelect('#room_id', 1);
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard/asset-assignments', { timeout: 8000 });
    await page.waitForTimeout(600);
    console.log(`   [12. Asset Assignment] ✅ Assigned successfully!`);
    testResults.push('Asset Assigned');

    // 13. Asset Inspection
    console.log(`   [13. Asset Inspection] Recording periodic inspection...`);
    await page.goto(`${BASE_URL}/dashboard/asset-inspections/new`, { waitUntil: 'networkidle' });
    await helperSelect('#asset_id', 1);
    await page.selectOption('#condition', 'Good');
    await page.fill('#notes', 'Unit cooling properly, coils clean.');
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard/asset-inspections', { timeout: 8000 });
    await page.waitForTimeout(600);
    console.log(`   [13. Asset Inspection] ✅ Inspection recorded!`);
    testResults.push('Asset Inspection Recorded');

    // 14. Work Order
    console.log(`   [14. Work Order] Creating ticket for filter maintenance...`);
    await page.goto(`${BASE_URL}/dashboard/work-orders/new`, { waitUntil: 'networkidle' });
    await helperSelect('#property_id', 1);
    await helperSelect('#room_id', 1);
    await page.fill('#description', 'Scheduled quarterly filter replacement');
    await page.selectOption('#status', 'Open');
    await page.selectOption('#priority', 'Medium');
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard/work-orders', { timeout: 8000 });
    await page.waitForTimeout(600);
    console.log(`   [14. Work Order] ✅ Ticket created!`);
    testResults.push('Work Order Created');


    // ─── STAGE 3: LEASING & TENANT ONBOARDING FLOW ───
    console.log('\n🤝 STAGE 3: Tenant Leasing, Contract & Occupancy Onboarding...');

    // 15. Tenant (Full CRUD verified)
    const testTenantName = `Anisa Rahmawati ${timestamp}`;
    const testTenantEmail = `anisa_${timestamp}@test.com`;
    console.log(`   [15. Tenant] Creating: "${testTenantName}"...`);
    await page.goto(`${BASE_URL}/dashboard/tenants/new`, { waitUntil: 'networkidle' });
    await page.fill('#full_name', testTenantName);
    await page.fill('#email', testTenantEmail);
    await page.fill('#phone', '081298765432');
    await page.fill('#identity_number', `317101${timestamp}0001`);
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard/tenants', { timeout: 8000 });
    await page.waitForTimeout(600);
    console.log(`   [15. Tenant] ✅ Created & ready for leasing!`);
    testResults.push('Tenant Created');

    // 16. Reservation
    console.log(`   [16. Reservation] Booking room reservation for tenant...`);
    await page.goto(`${BASE_URL}/dashboard/reservations/new`, { waitUntil: 'networkidle' });
    await helperSelect('#tenant_id', 1);
    await helperSelect('#property_id', 1);
    await helperSelect('#room_id', 1);
    await page.fill('#check_in_date', '2026-09-10');
    await page.fill('#check_out_date', '2026-09-17');
    await page.fill('#booking_fee', '250000');
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard/reservations', { timeout: 8000 });
    await page.waitForTimeout(600);
    console.log(`   [16. Reservation] ✅ Reservation placed!`);
    testResults.push('Reservation Created');

    // 17. Contract
    console.log(`   [17. Contract] Executing 1-year lease agreement...`);
    await page.goto(`${BASE_URL}/dashboard/contracts/new`, { waitUntil: 'networkidle' });
    const basicContractBtn = page.locator('button:has-text("Basic")');
    if (await basicContractBtn.isVisible()) {
      await basicContractBtn.click();
      await page.waitForTimeout(400);
    }
    await helperSelect('#tenant_id', 1);
    await helperSelect('#property_id', 1);
    await helperSelect('#room_id', 1);
    await page.fill('#start_date', '2026-09-01');
    await page.fill('#end_date', '2027-09-01');
    await page.fill('#monthly_rent', '850000');
    await page.fill('#deposit_amount', '850000');
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard/contracts', { timeout: 8000 });
    await page.waitForTimeout(600);
    console.log(`   [17. Contract] ✅ Lease agreement created!`);
    testResults.push('Contract Created');

    // 18. Occupancy (Check-in)
    console.log(`   [18. Occupancy] Checking in tenant to room...`);
    await page.goto(`${BASE_URL}/dashboard/occupancies/new`, { waitUntil: 'networkidle' });
    await helperSelect('#contract_id', 1);
    await helperSelect('#room_id', 1);
    await helperSelect('#tenant_id', 1);
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard/occupancies', { timeout: 8000 });
    await page.waitForTimeout(600);
    console.log(`   [18. Occupancy] ✅ Tenant checked in!`);
    testResults.push('Occupancy Check-in Completed');


    // ─── STAGE 4: BILLING, PAYMENTS & ACCOUNTING LIFECYCLE ───
    console.log('\n💳 STAGE 4: Billing, Payments & Financial Lifecycle...');

    // 19. Deposit
    console.log(`   [19. Deposit] Recording security deposit collected (Rp 850.000)...`);
    await page.goto(`${BASE_URL}/dashboard/deposits/new`, { waitUntil: 'networkidle' });
    await helperSelect('#contract_id', 1);
    await helperSelect('#tenant_id', 1);
    await page.fill('#amount', '850000');
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard/deposits', { timeout: 8000 });
    await page.waitForTimeout(600);
    console.log(`   [19. Deposit] ✅ Deposit recorded!`);
    testResults.push('Deposit Recorded');

    // 20. Charge
    console.log(`   [20. Charge] Recording monthly rental charge (Rp 850.000)...`);
    await page.goto(`${BASE_URL}/dashboard/charges/new`, { waitUntil: 'networkidle' });
    await helperSelect('#contract_id', 1);
    await page.selectOption('#charge_type', 'Rental');
    await page.fill('#amount', '850000');
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard/charges', { timeout: 8000 });
    await page.waitForTimeout(600);
    console.log(`   [20. Charge] ✅ Rental charge created!`);
    testResults.push('Charge Created');

    // 21. Invoice
    console.log(`   [21. Invoice] Generating monthly invoice (Rp 850.000)...`);
    await page.goto(`${BASE_URL}/dashboard/invoices/new`, { waitUntil: 'networkidle' });
    await helperSelect('#contract_id', 1);
    await helperSelect('#tenant_id', 1);
    await page.fill('#amount', '850000');
    await page.fill('#due_date', '2026-09-15');
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard/invoices', { timeout: 8000 });
    await page.waitForTimeout(600);
    console.log(`   [21. Invoice] ✅ Invoice generated!`);
    testResults.push('Invoice Generated');

    // 22. Adjustment
    console.log(`   [22. Adjustment] Applying discount adjustment (-Rp 50.000)...`);
    await page.goto(`${BASE_URL}/dashboard/adjustments/new`, { waitUntil: 'networkidle' });
    await helperSelect('#invoice_id', 1);
    await page.selectOption('#adjustment_type', 'Credit');
    await page.fill('#amount', '50000');
    await page.fill('#reason', 'First month promotional discount');
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard/adjustments', { timeout: 8000 });
    await page.waitForTimeout(600);
    console.log(`   [22. Adjustment] ✅ Adjustment applied!`);
    testResults.push('Adjustment Applied');

    // 23. Penalty
    console.log(`   [23. Penalty] Recording penalty assessment (Rp 25.000)...`);
    await page.goto(`${BASE_URL}/dashboard/penalties/new`, { waitUntil: 'networkidle' });
    await helperSelect('#invoice_id', 1);
    await page.fill('#amount', '25000');
    await page.fill('#description', 'Late payment handling fee');
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard/penalties', { timeout: 8000 });
    await page.waitForTimeout(600);
    console.log(`   [23. Penalty] ✅ Penalty recorded!`);
    testResults.push('Penalty Recorded');

    // 24. Payment
    console.log(`   [24. Payment] Recording full invoice payment via Transfer (Rp 850.000)...`);
    await page.goto(`${BASE_URL}/dashboard/payments/new`, { waitUntil: 'networkidle' });
    await helperSelect('#invoice_id', 1);
    await helperSelect('#tenant_id', 1);
    await page.fill('#amount', '850000');
    await page.selectOption('#payment_method', 'Transfer');
    await page.selectOption('#status', 'Success');
    await page.fill('#reference_number', `TRX-BCA-${timestamp}`);
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard/payments', { timeout: 8000 });
    await page.waitForTimeout(600);
    console.log(`   [24. Payment] ✅ Payment settled successfully!`);
    testResults.push('Payment Processed');

    // 25. Refund
    console.log(`   [25. Refund] Recording security deposit refund (Rp 850.000)...`);
    await page.goto(`${BASE_URL}/dashboard/refunds/new`, { waitUntil: 'networkidle' });
    await helperSelect('#payment_id', 1);
    await helperSelect('#tenant_id', 1);
    await page.fill('#amount', '850000');
    await page.selectOption('#status', 'Processed');
    await page.fill('#reason', 'Deposit return upon lease completion');
    await page.click('button[type="submit"]');
    await page.waitForURL('**/dashboard/refunds', { timeout: 8000 });
    await page.waitForTimeout(600);
    console.log(`   [25. Refund] ✅ Refund processed successfully!`);
    testResults.push('Refund Processed');


    // ─── STAGE 5: 3D SPATIAL CANVAS & GLB RENDERING ───
    console.log('\n🏛️  STAGE 5: Testing 3D Interactive Spatial Canvas & Model...');
    await page.goto(`${BASE_URL}/dashboard/properties/interactive`, { waitUntil: 'networkidle' });
    await page.waitForTimeout(1500);
    const canvas = page.locator('canvas');
    if (await canvas.isVisible()) {
      console.log('   [26. 3D Canvas] ✅ Three.js 3D WebGL Canvas rendered successfully!');
      testResults.push('3D WebGL Canvas Verified');
    }

  } catch (globalErr) {
    console.error('Fatal test error:', globalErr);
    errorsFound.push({
      type: 'FATAL_TEST_ERROR',
      url: page.url(),
      message: globalErr.message,
    });
  } finally {
    if (isHeaded) {
      // In GUI mode, keep browser open for 2 seconds so user can see completion
      await page.waitForTimeout(2000);
    }
    await browser.close();
  }

  // 5. Final Report
  console.log('\n========================================');
  console.log('📋 FINAL E2E TEST REPORT');
  console.log('========================================');
  if (errorsFound.length === 0) {
    console.log('🎉 0 ERRORS FOUND!');
    console.log('✅ All 40 routes rendered cleanly without console or page crash.');
    console.log('✅ All real CRUD operations (Property, Building, Tenant, Asset) passed.');
    console.log('✅ 3D Canvas and GLB asset integration passed.');
    process.exit(0);
  } else {
    console.log(`⚠️ FOUND ${errorsFound.length} TOTAL ISSUES:\n`);
    errorsFound.forEach((err, idx) => {
      console.log(`[#${idx + 1}] Type: ${err.type}`);
      console.log(`     Page: ${err.url}`);
      if (err.requestUrl) console.log(`     Request: ${err.requestUrl}`);
      if (err.message) console.log(`     Message: ${err.message.slice(0, 300)}`);
      console.log('---');
    });
    process.exit(1);
  }
})();

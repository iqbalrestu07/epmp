import { chromium } from "playwright-core";

const CHROME_PATH = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome";
const BASE_URL = process.env.BASE_URL || "http://localhost:3000";
const ts = Date.now();
const EMAIL = `debug.${ts}@example.com`;
const PASSWORD = "DebugPass123!";

const browser = await chromium.launch({ executablePath: CHROME_PATH, headless: true });
const page = await browser.newPage();

const logs = [];
page.on("console", (msg) => {
  logs.push({ type: msg.type(), text: msg.text() });
});
page.on("pageerror", (err) => {
  logs.push({ type: "pageerror", text: err.stack || err.message });
});
page.on("response", async (res) => {
  const url = res.url();
  if (url.includes("/api/v1/") && res.request().method() === "GET") {
    try {
      const body = await res.text();
      console.log(`\n[API RESPONSE] ${res.status()} ${url}\n${body.slice(0, 500)}`);
    } catch {}
  }
});

function flushLogs(label) {
  console.log(`\n--- Logs after: ${label} ---`);
  const relevant = logs.filter((l) => l.type === "error" || l.type === "pageerror" || l.type === "warning");
  if (relevant.length === 0) console.log("(none)");
  for (const l of relevant) console.log(`[${l.type}]`, l.text);
  logs.length = 0;
}

// 1. Sign up
console.log(`Signing up as ${EMAIL} ...`);
await page.goto(`${BASE_URL}/auth/signup`, { waitUntil: "networkidle" });
await page.fill('input[type="text"]', "Debug User");
await page.fill('input[type="email"]', EMAIL);
const passwordInputs = await page.locator('input[type="password"]').all();
await passwordInputs[0].fill(PASSWORD);
await passwordInputs[1].fill(PASSWORD);
await page.click('button[type="submit"]');
await page.waitForTimeout(2000);
flushLogs("signup");
console.log("URL after signup:", page.url());

// 2. Create organization (should be redirected here automatically)
if (page.url().includes("/organizations/new")) {
  console.log("Creating organization...");
  await page.fill("#name", "Debug Org");
  await page.click('button[type="submit"]');
  await page.waitForTimeout(2000);
  flushLogs("create organization");
  console.log("URL after org create:", page.url());
}

// 3. Create a property (needed for building/asset forms to have data)
console.log("Creating property...");
await page.goto(`${BASE_URL}/dashboard/properties/new`, { waitUntil: "networkidle" });
await page.waitForTimeout(500);
flushLogs("visit property create page");

const routesToCheck = process.argv.slice(2);
if (routesToCheck.length === 0) {
  routesToCheck.push(
    "/dashboard/buildings/new",
    "/dashboard/properties/interactive",
    "/dashboard/explorer",
    "/dashboard/assets",
  );
}

for (const route of routesToCheck) {
  console.log(`\n=== Visiting ${route} ===`);
  try {
    await page.goto(`${BASE_URL}${route}`, { waitUntil: "networkidle", timeout: 15000 });
  } catch (e) {
    console.log("Navigation error:", e.message);
  }
  await page.waitForTimeout(2500);
  flushLogs(route);
}

// Extra: click into "Basic View (Form)" tab on building create page specifically
console.log(`\n=== Visiting /dashboard/buildings/new and clicking Basic View tab ===`);
await page.goto(`${BASE_URL}/dashboard/buildings/new`, { waitUntil: "networkidle" });
await page.waitForTimeout(1000);
try {
  await page.click("text=Basic View (Form)");
} catch (e) {
  console.log("Could not click Basic View tab:", e.message);
}
await page.waitForTimeout(2000);
flushLogs("buildings/new basic view tab");

await browser.close();


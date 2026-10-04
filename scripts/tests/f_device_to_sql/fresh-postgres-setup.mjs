const CONFIG_PATH = '/api/v1/datalink/studio-v2/workspace/database-config';
const CONFIG_ROUTE = `**${CONFIG_PATH}**`;
const OWNED_FIELDS = ['kind', 'host', 'port', 'database', 'username', 'schema', 'table'];

function requestPath(request) {
  try {
    return new URL(request.url()).pathname;
  } catch {
    return '';
  }
}

function requestPayload(request) {
  try {
    const body = request.postDataJSON();
    if (!body || typeof body !== 'object' || Array.isArray(body)) return null;
    for (const key of ['config', 'connector', 'data']) {
      if (body[key] && typeof body[key] === 'object' && !Array.isArray(body[key])) return body[key];
    }
    return body;
  } catch {
    return null;
  }
}

function sameConfigValue(field, actual, target) {
  if (field === 'port') return Number(actual) === Number(target);
  return typeof actual === 'string' && actual === target;
}

/** Pure decision used by the browser route and its no-browser regression tests. */
export function inspectPostgresConfigRequest(request, target) {
  const method = request.method().toUpperCase();
  if (requestPath(request) !== CONFIG_PATH || !['POST', 'PUT'].includes(method)) {
    return { action: 'continue', blockedFields: [] };
  }
  const payload = requestPayload(request);
  const blockedFields = payload
    ? OWNED_FIELDS.filter((field) => !sameConfigValue(field, payload[field], field === 'kind' ? 'postgres' : target[field]))
    : ['body'];
  return blockedFields.length === 0
    ? { action: 'continue', blockedFields: [] }
    : { action: 'abort', blockedFields };
}

export function createOwnedPostgresRequestGuard(target) {
  const ownedTarget = Object.fromEntries(OWNED_FIELDS
    .filter((field) => field !== 'kind')
    .map((field) => [field, target[field]]));
  const blockedFields = new Set();
  let blockedCount = 0;
  return {
    inspect(request) {
      const decision = inspectPostgresConfigRequest(request, ownedTarget);
      if (decision.action === 'abort') {
        blockedCount += 1;
        decision.blockedFields.forEach((field) => blockedFields.add(field));
      }
      return decision;
    },
    snapshot() {
      return {
        blocked_count: blockedCount,
        blocked_fields: [...blockedFields].sort(),
      };
    },
  };
}

function validateTarget(target) {
  const required = ['host', 'port', 'database', 'username', 'schema', 'table', 'password'];
  if (!target || required.some((field) => target[field] === undefined || target[field] === null || target[field] === '')) {
    throw new Error('fresh PostgreSQL target is incomplete');
  }
  if (!Number.isInteger(Number(target.port)) || Number(target.port) < 1 || Number(target.port) > 65535) {
    throw new Error('fresh PostgreSQL target port is invalid');
  }
}

async function readPersistedConfig(base, timeoutMs) {
  const response = await fetch(`${base.replace(/\/+$/, '')}${CONFIG_PATH}`, {
    signal: AbortSignal.timeout(Math.max(300, timeoutMs)),
  });
  const body = await response.json();
  if (!response.ok) throw new Error('database config read failed');
  const config = body?.data ?? body;
  if (!config || typeof config !== 'object' || Array.isArray(config)) return config;
  const safeConfig = { ...config };
  for (const secretField of ['password', 'password_masked', 'dsn', 'connection_string']) delete safeConfig[secretField];
  return safeConfig;
}

function exactConfig(config, target) {
  return config?.save_state === 'saved' && config.kind === 'postgres' && config.name === 'F fresh PostgreSQL' &&
    config.host === target.host && Number(config.port) === Number(target.port) &&
    config.database === target.database && config.username === target.username &&
    config.schema === target.schema && config.table === target.table;
}

async function fillOwned(fields, target) {
  const values = [
    [fields.name, 'F fresh PostgreSQL'], [fields.host, target.host], [fields.port, String(target.port)],
    [fields.database, target.database], [fields.username, target.username],
    [fields.schema, target.schema], [fields.table, target.table],
    // Identity edits intentionally require a fresh password; enter it last.
    [fields.password, target.password],
  ];
  for (const [input, value] of values) {
    await input.fill(value);
    await input.press('Tab');
  }
}

async function assertSavedDOM(page, fields, target, timeoutMs) {
  const details = page.getByTestId('step4-advanced-recording');
  if ((await details.getAttribute('open')) === null) await details.locator('summary').click();
  await fields.name.waitFor({ state: 'visible', timeout: timeoutMs });
  const expected = [
    [fields.name, 'F fresh PostgreSQL'], [fields.host, target.host], [fields.port, String(target.port)],
    [fields.database, target.database], [fields.username, target.username],
    [fields.schema, target.schema], [fields.table, target.table],
  ];
  for (const [input, value] of expected) {
    if (await input.inputValue() !== value) throw new Error('PostgreSQL connector DOM did not retain saved values');
  }
}

/** Configure PostgreSQL through Step 4 and return safe persisted/guard evidence. */
export async function configureFreshPostgres(page, base, target) {
  validateTarget(target);
  const guard = createOwnedPostgresRequestGuard(target);
  const routeHandler = async (route) => {
    const decision = guard.inspect(route.request());
    if (decision.action === 'abort') return route.abort('blockedbyclient');
    return route.continue();
  };
  await page.route(CONFIG_ROUTE, routeHandler);
  const fields = {
    name: page.getByLabel(/Connection Name|連線名稱/i).first(),
    host: page.getByLabel(/Host IP Address|主機 IP/i).first(),
    port: page.getByLabel(/Port|連接埠/i).first(),
    database: page.getByLabel(/Database Name|資料庫名稱/i).first(),
    username: page.getByLabel(/Username|使用者名稱/i).first(),
    password: page.getByLabel(/Password|密碼/i).first(),
    schema: page.getByLabel(/^Schema$|^Schema（選填）$/i).first(),
    table: page.getByLabel(/Table Name|資料表名稱/i).first(),
  };
  const deadline = Date.now() + 30_000;
  let savedConfig;
  try {
    await page.getByTestId('step-nav-button-4').click();
    await page.getByRole('button', { name: /PostgreSQL/i }).click();
    await fields.name.waitFor({ state: 'visible', timeout: 10_000 });
    for (let attempt = 0; attempt < 8 && Date.now() < deadline; attempt += 1) {
      await fillOwned(fields, target);
      await page.waitForTimeout(250);
      try {
        const observed = await readPersistedConfig(base, Math.min(4_000, deadline - Date.now()));
        if (exactConfig(observed, target)) {
          savedConfig = observed;
          break;
        }
      } catch { /* bounded retry; details are deliberately not exposed */ }
      await page.waitForTimeout(Math.min(350, Math.max(0, deadline - Date.now())));
    }
    if (!savedConfig) {
      throw new Error(`fresh PostgreSQL config did not reach saved state; blocked fields: ${guard.snapshot().blocked_fields.join(',') || 'none'}`);
    }
    await page.reload();
    await page.getByTestId('step-nav-button-4').click();
    await page.getByTestId('basic-recording-panel').waitFor({ timeout: 20_000 });
    await assertSavedDOM(page, fields, target, Math.min(10_000, Math.max(300, deadline - Date.now())));
    savedConfig = await readPersistedConfig(base, Math.min(5_000, Math.max(300, deadline - Date.now())));
    if (!exactConfig(savedConfig, target)) throw new Error('fresh PostgreSQL config readback changed after reload');
    if (!await page.getByTestId('basic-recording-start').isEnabled()) {
      throw new Error('persisted PostgreSQL connector did not make the Basic Start control available');
    }
    return { config: savedConfig, request_guard: guard.snapshot() };
  } finally {
    await page.unroute(CONFIG_ROUTE, routeHandler).catch(() => undefined);
  }
}

export { CONFIG_PATH };

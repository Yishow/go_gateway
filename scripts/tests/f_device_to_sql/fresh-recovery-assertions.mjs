import assert from 'node:assert/strict';

/**
 * Assertions shared by the fresh F recovery witness. They deliberately accept
 * only structured facts returned by the real UI/API/SQL adapters; no assertion
 * manufactures a successful operation or target row.
 */

export function assertOperationIdentity(first, second, label) {
  assert.ok(first?.operation_id, `${label}: first operation id is missing`);
  assert.ok(second?.operation_id, `${label}: replay operation id is missing`);
  assert.equal(second.operation_id, first.operation_id, `${label}: replay changed operation identity`);
  if (first.action !== undefined || second.action !== undefined) {
    assert.equal(second.action, first.action, `${label}: replay changed operation action`);
  }
}

export function assertOneEffect(effects, label) {
  assert.ok(Array.isArray(effects), `${label}: effect observation is not an array`);
  const keys = effects.map((effect) => String(effect?.effect_key ?? effect?.operation_id ?? ''));
  assert.ok(keys.every(Boolean), `${label}: an effect has no durable identity`);
  assert.equal(new Set(keys).size, keys.length, `${label}: duplicate durable effect identity`);
  assert.equal(keys.length, 1, `${label}: expected one durable effect`);
}

export function assertAcquisitionProvenance(row, expected, label = 'delayed resend') {
  assert.ok(row, `${label}: SQL row is missing`);
  const acquisition = row.acquisition_time ?? row.bucket_start ?? row.observed_at;
  const acquisitionMillis = Date.parse(String(acquisition ?? ''));
  const expectedMillis = Date.parse(String(expected.observed_at ?? ''));
  assert.ok(Number.isFinite(acquisitionMillis) && Number.isFinite(expectedMillis), `${label}: acquisition time is not valid UTC`);
  assert.equal(acquisitionMillis, expectedMillis, `${label}: acquisition time changed`);
  const rawProvenance = row.provenance ?? row.prov;
  assert.ok(rawProvenance, `${label}: SQL row has no provenance evidence`);
  let provenance;
  try { provenance = typeof rawProvenance === 'string' ? JSON.parse(rawProvenance) : rawProvenance; } catch { provenance = null; }
  assert.ok(Array.isArray(provenance) && provenance.length > 0, `${label}: provenance evidence is not a non-empty array`);
  assert.ok(provenance.every((entry) => typeof entry?.observed_at === 'string' && Number.isFinite(Date.parse(entry.observed_at)) && /Z$/.test(entry.observed_at)), `${label}: provenance has invalid acquisition time`);
  if (expected.provenance !== undefined) assert.deepEqual(provenance, expected.provenance, `${label}: provenance changed during delayed resend`);
}

export function assertNeighborUnchanged(before, after, label = 'test-write neighbor') {
  assert.deepEqual(after, before, `${label}: non-owned neighbor changed`);
}

export function assertTargetRowsUnchanged(before, after, owner, ownerColumn = 'owner', label = 'test-write target') {
  const withoutOwner = (rows) => rows.filter((row) => String(row?.[ownerColumn] ?? row?.owner ?? row?.line ?? '') !== owner);
  assert.deepEqual(withoutOwner(after), withoutOwner(before), `${label}: non-owned target rows changed`);
}

export function assertTestWriteOutcome(operation, label = 'test-write') {
  assert.ok(operation, `${label}: operation result is missing`);
  assert.equal(operation.write_outcome, 'written_verified', `${label}: write was not verified`);
  assert.equal(operation.cleanup_status, 'cleaned', `${label}: owned row was not cleaned`);
  assert.ok(operation.operation_id, `${label}: operation identity is missing`);
}

export function assertNoOwnedRows(rows, owner, ownerColumn = 'owner', label = 'test-write cleanup') {
  assert.ok(Array.isArray(rows), `${label}: SQL row observation is not an array`);
  const remaining = rows.filter((row) => String(row?.[ownerColumn] ?? row?.owner ?? row?.line ?? '') === owner);
  assert.equal(remaining.length, 0, `${label}: owned row remains in target`);
}

export function sanitizeEvidence(value, depth = 0) {
  if (depth > 8) return '[depth-limited]';
  if (value === null || typeof value === 'boolean' || typeof value === 'number') return value;
  if (typeof value === 'string') {
    const scrubbed = value
      .replace(/\/(?:Users|private|tmp|var|home)\/[^\s"'`]+/g, '[private-value]')
      .replace(/(?:password|passwd|pwd)=[^&\s]+/gi, 'password=[redacted]')
      .replace(/postgres(?:ql)?:\/\/[^\s"'`]+/gi, '[private-value]');
    if (scrubbed.startsWith('/') || /^[A-Za-z]:[\\/]/.test(scrubbed) || /^(?:file:|postgres(?:ql)?:)/i.test(scrubbed)) return '[private-value]';
    return scrubbed;
  }
  if (Array.isArray(value)) return value.slice(0, 128).map((entry) => sanitizeEvidence(entry, depth + 1));
  if (typeof value !== 'object') return String(value);
  const output = {};
  for (const [key, entry] of Object.entries(value).slice(0, 128)) {
    if (/(?:password|secret|credential|authorization|cookie|dsn|token)/i.test(key)) {
      output[key] = '[redacted]';
    } else {
      output[key] = sanitizeEvidence(entry, depth + 1);
    }
  }
  return output;
}

export async function inspectBasicLayout(page) {
  return page.evaluate(() => {
    const basic = document.querySelector('[data-testid="basic-recording-panel"]');
    const overflowing = [...document.querySelectorAll('body *')].filter((element) => {
      if (typeof element.checkVisibility === 'function' && !element.checkVisibility()) return false;
      const box = element.getBoundingClientRect();
      return box.width > 0 && box.height > 0 && box.right > innerWidth + 1;
    });
    const describeOverflow = (elements) => elements.map((element) => ({ tag: element.tagName, test_id: element.dataset.testid, right: element.getBoundingClientRect().right, width: element.getBoundingClientRect().width }));
    const describe = (element) => element ? { left: element.getBoundingClientRect().left, right: element.getBoundingClientRect().right, width: element.getBoundingClientRect().width } : null;
    return { document_width: document.documentElement.scrollWidth, basic_bounds: describe(basic), basic_outside: describeOverflow(basic ? overflowing.filter((element) => basic.contains(element)) : []), legacy_outside: describeOverflow(basic ? overflowing.filter((element) => !basic.contains(element)) : overflowing) };
  });
}

export function assertBasicLayout(layout, width) {
  assert.ok(layout?.basic_bounds && layout.basic_bounds.left >= 0 && layout.basic_bounds.right <= width + 1 && layout.basic_outside.length === 0, `Basic controls overflow at ${width}px`);
}

export async function screenshotAnchor(page, testIds) {
  for (const testId of testIds) {
    const locator = page.getByTestId(testId);
    if (await locator.isVisible().catch(() => false)) return { testId, locator };
  }
  throw new Error(`no visible screenshot anchor: ${testIds.join(',')}`);
}

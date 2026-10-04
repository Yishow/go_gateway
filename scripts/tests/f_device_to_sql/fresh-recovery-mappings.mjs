import assert from 'node:assert/strict';
import { observeJSON } from './fresh-ui.mjs';

function mappingList(data) {
  return Array.isArray(data) ? data : data?.mappings ?? data?.items ?? [];
}

function mappingMutationPath(pathname) {
  return /^\/api\/v1\/datalink\/studio-v2\/workspace\/mappings(?:\/[^/]+)?$/.test(pathname);
}

function safeMutationBody(body) {
  if (!body || typeof body !== 'object') return body;
  return Object.fromEntries(['id', 'point_id', 'workspace_id', 'rule_id', 'device_id', 'address', 'tag_id', 'tag_key', 'display_name', 'target_type', 'unit', 'scale', 'offset', 'enabled'].filter((key) => key in body).map((key) => [key, body[key]]));
}

function installMappingMutationRecorder(page) {
  const events = [];
  const pending = new Set();
  let activeField = '';
  const onResponse = (response) => {
    const request = response.request();
    const url = new URL(response.url());
    if (!mappingMutationPath(url.pathname) || !['POST', 'PUT'].includes(request.method())) return;
    const task = (async () => {
      let responseBody = null;
      try {
        responseBody = safeMutationBody((await response.json())?.data);
      } catch {
        try { responseBody = (await response.text()).slice(0, 800); } catch { responseBody = '<unreadable response>'; }
      }
      let requestBody = null;
      try { requestBody = safeMutationBody(request.postDataJSON()); } catch { /* request may have no JSON body */ }
      events.push({ field: activeField, method: request.method(), path: url.pathname, status: response.status(), request: requestBody, response: responseBody });
    })();
    pending.add(task);
    void task.finally(() => pending.delete(task));
  };
  page.on('response', onResponse);
  return {
    events,
    setField(field) { activeField = field; },
    async flush() { await Promise.allSettled([...pending]); },
    async close() { await this.flush(); page.off('response', onResponse); },
  };
}

async function mappingDOMFacts(page, pointID, domField) {
  const field = page.getByTestId(`${domField === 'target-type' ? 'select-target-type' : `input-${domField}`}-${pointID}`);
  const badge = page.getByTestId(`mapping-save-state-${pointID}`);
  return {
    point_id: pointID,
    field: domField,
    input_value: await field.inputValue().catch(() => '<unreadable>'),
    input_visible: await field.isVisible().catch(() => false),
    input_count: await field.count().catch(() => 0),
    badge_text: await badge.innerText().catch(() => ''),
    badge_content: await badge.textContent().catch(() => ''),
    badge_visible: await badge.isVisible().catch(() => false),
    badge_count: await badge.count().catch(() => 0),
  };
}

async function waitForMappingField(page, base, pointID, tag, apiField, domField, expected, recorder) {
  const deadline = Date.now() + 15_000;
  let last = [];
  let domFacts = await mappingDOMFacts(page, pointID, domField);
  while (Date.now() < deadline) {
    last = mappingList(await observeJSON(base, '/studio-v2/workspace/mappings'));
    const mutation = recorder?.events.findLast((event) => event.response?.tag_key === tag) ?? null;
    const persistedPointID = mutation?.response?.point_id ?? '';
    const record = last.find((item) => item.tag_key === tag && (!persistedPointID || item.point_id === persistedPointID));
    domFacts = await mappingDOMFacts(page, pointID, domField);
    if (persistedPointID && record?.point_id === persistedPointID && record?.[apiField] === expected && domFacts.input_value === expected && /saved|已儲存/i.test(domFacts.badge_text)) return last;
    await new Promise((resolve) => setTimeout(resolve, 500));
  }
  await recorder?.flush();
  const mutation = recorder?.events.findLast((event) => event.response?.tag_key === tag) ?? null;
  const persistedPointID = mutation?.response?.point_id ?? '';
  const record = last.find((item) => item.tag_key === tag && (!persistedPointID || item.point_id === persistedPointID)) ?? null;
  throw new Error(`mapping field did not persist: ${JSON.stringify({
    point_id: pointID,
    tag_key: tag,
    api_field: apiField,
    expected,
    api_point_id: persistedPointID,
    api_record: safeMutationBody(record),
    api_records_for_tag: last.filter((item) => item.tag_key === tag).slice(0, 3).map(safeMutationBody),
    dom: domFacts,
    mutation_events: recorder?.events.slice(-8) ?? [],
  }).slice(0, 4000)}`);
}

export async function persistRecoveryMappings(page, base, { tagPrefix, targetTypes, expectedCount }) {
  const rows = page.locator('[data-testid^="mapping-row-"]');
  await rows.first().waitFor({ timeout: 20_000 });
  await page.waitForFunction((expected) => document.querySelectorAll('[data-testid^="mapping-row-"]').length >= expected, expectedCount, { timeout: 20_000 });
  const pointIDs = await rows.evaluateAll((elements) => elements.map((element) => element.getAttribute('data-testid')?.replace('mapping-row-', '')).filter(Boolean));
  assert.equal(pointIDs.length, expectedCount, `fresh mapping row count is ${pointIDs.length}, expected ${expectedCount}`);
  const tags = [];
  const recorder = installMappingMutationRecorder(page);
  try {
    for (let index = 0; index < pointIDs.length; index += 1) {
      const pointID = pointIDs[index];
      const tag = `${tagPrefix}.${index + 1}`;
      const targetType = targetTypes[index];
      let saved = false;
      for (let attempt = 0; attempt < 3 && !saved; attempt += 1) {
        try {
          const tagInput = page.getByTestId(`input-tag-key-${pointID}`);
          recorder.setField(`${pointID}.tag_key`);
          await tagInput.fill(tag);
          await tagInput.press('Tab');
          await waitForMappingField(page, base, pointID, tag, 'tag_key', 'tag-key', tag, recorder);
          const displayInput = page.getByTestId(`input-display-name-${pointID}`);
          recorder.setField(`${pointID}.display_name`);
          await displayInput.fill(tag);
          await displayInput.press('Tab');
          await waitForMappingField(page, base, pointID, tag, 'display_name', 'display-name', tag, recorder);
          if (targetType) {
            const selector = page.getByTestId(`select-target-type-${pointID}`);
            recorder.setField(`${pointID}.target_type`);
            await selector.selectOption(targetType);
            await selector.press('Tab');
            await waitForMappingField(page, base, pointID, tag, 'target_type', 'target-type', targetType, recorder);
          }
          saved = true;
        } catch (error) {
          if (attempt === 2) throw error;
          await new Promise((resolve) => setTimeout(resolve, 500));
        }
      }
      assert.ok(saved, `mapping was not saved for ${pointID}`);
      const finalRecord = mappingList(await observeJSON(base, '/studio-v2/workspace/mappings')).find((item) => item.tag_key === tag);
      assert.ok(finalRecord?.point_id, `mapping ${tag} did not expose persisted point identity`);
      tags.push({ point_id: pointID, persisted_point_id: finalRecord.point_id, tag_key: tag });
    }
    const final = mappingList(await observeJSON(base, '/studio-v2/workspace/mappings'));
    assert.ok(tags.every((tag) => final.some((item) => item.point_id === tag.persisted_point_id && item.tag_key === tag.tag_key)), 'not all fresh mappings are persisted');
    await page.getByTestId('btn-continue').click();
    await page.getByTestId('step-nav-button-4').waitFor({ timeout: 20_000 });
    return tags;
  } finally {
    await recorder.close();
  }
}

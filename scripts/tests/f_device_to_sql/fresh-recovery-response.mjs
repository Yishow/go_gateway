import assert from 'node:assert/strict';

function responsePath(request) {
  return new URL(request.url()).pathname;
}

function requestEvidence(request) {
  let body = null;
  try { body = request.postDataJSON(); } catch { /* GET or non-JSON request */ }
  if (!body || typeof body !== 'object') return null;
  const evidence = {
    request_id: body.request_id,
    operation_id: body.operation_id,
    workspace_id: body.workspace_id,
    group_id: body.group_id,
    expected_workspace_revision: body.expected_workspace_revision,
    expected_group_revision: body.expected_group_revision,
    expected_connector_revision: body.expected_connector_revision,
    groups: Array.isArray(body.groups) ? body.groups.map((group) => ({
      group_id: group.group_id,
      expected_group_revision: group.expected_group_revision,
      expected_connector_revision: group.expected_connector_revision,
    })) : undefined,
  };
  return Object.fromEntries(Object.entries(evidence).filter(([, value]) => value !== undefined));
}

/**
 * Drops exactly one browser response only after route.fetch() has completed.
 * The server therefore performed the real mutation before the browser loses
 * its acknowledgement. The fetched body is retained only as sanitized test
 * evidence and is never fed back to the page.
 */
export async function dropOneResponseAfterEffect(page, { method = 'POST', path, predicate } = {}) {
  assert.ok(page, 'response-drop requires a Playwright page');
  assert.ok(path || predicate, 'response-drop requires a path or predicate');
  let observed;
  let matched = false;
  let resolveDone;
  let rejectDone;
  const done = new Promise((resolve, reject) => {
    resolveDone = resolve;
    rejectDone = reject;
  });
  const handler = async (route) => {
    const request = route.request();
    const matches = !matched && request.method() === method &&
      (predicate ? predicate(request) : responsePath(request).includes(path));
    if (!matches) {
      await route.continue();
      return;
    }
    matched = true;
    try {
      const response = await route.fetch();
      const bodyText = await response.text();
      let body;
      try { body = JSON.parse(bodyText); } catch { body = { raw: bodyText.slice(0, 512) }; }
      observed = { status: response.status(), request: requestEvidence(request), body };
      // Keep the exact browser request available to the recovery harness so it
      // can exercise a real same-request retry after the acknowledgement drop.
      // It is deliberately non-enumerable: public evidence must retain only
      // the sanitized identity fields above.
      Object.defineProperty(observed, 'replayRequest', { value: (() => {
        try { return request.postDataJSON(); } catch { return undefined; }
      })(), enumerable: false });
      await route.abort('failed');
      resolveDone(observed);
    } catch (error) {
      rejectDone(error);
    }
  };
  await page.route('**/api/v1/datalink/**', handler);
  return {
    done,
    async close() {
      await page.unroute('**/api/v1/datalink/**', handler);
    },
    get matched() { return matched; },
  };
}

export async function observeResponse(page, { method = 'POST', path }, action) {
  const responsePromise = page.waitForResponse((response) => {
    const request = response.request();
    return request.method() === method && responsePath(request).includes(path);
  });
  await action();
  const response = await responsePromise;
  const body = await response.json();
  return { status: response.status(), request: requestEvidence(response.request()), body };
}

// Report cleanup before publishing success; stop only this run's child PIDs.
import { stopProcesses } from './lib.mjs';

function ownedChildren(children) {
  if (!children) return [];
  const values = Array.isArray(children) ? children : Object.values(children);
  return values.filter((child) => child && Number.isInteger(child.pid));
}

function boundedError(error) {
  return String(error instanceof Error ? error.message : error).slice(0, 256);
}

function markCleanupFailure(outcome) {
  if (outcome?.passed === false) process.exitCode = 1;
  return outcome;
}

/** Keep diagnostics from turning a cleanup failure into an unreported leak. */
export function safeRead(label, read, fallback, errors = []) {
  try { return read(); }
  catch (error) {
    errors.push({ label, reason: 'diagnostic-read-failed', message: boundedError(error) });
    return fallback;
  }
}

export async function safeReadAsync(label, read, fallback, errors = []) {
  try { return await read(); }
  catch (error) {
    errors.push({ label, reason: 'diagnostic-read-failed', message: boundedError(error) });
    return fallback;
  }
}

export async function cleanupFixture(browser, children, waitForPortFree, permissions, targets) {
  const steps = {};
  const check = async (name, action) => {
    try {
      const detail = await action();
      markCleanupFailure(detail);
      steps[name] = { passed: true, ...(detail && typeof detail === 'object' ? detail : {}) };
    }
    catch { steps[name] = { passed: false, reason: 'fixture-cleanup-failed' }; process.exitCode = 1; }
  };
  if (typeof browser?.close === 'function') await check('browser_closed', () => browser.close());
  else steps.browser_closed = { passed: true, reason: 'not-created' };
  const childrenToStop = ownedChildren(children);
  await check('owned_children_stopped', async () => {
    const statuses = await stopProcesses(childrenToStop);
    const failed = statuses.filter((status) => status.state !== 'exited'
      || (status.exit_code !== null && status.exit_code !== undefined && status.exit_code !== 0)
      || Object.entries(status).some(([key, value]) => key.endsWith('_error') && value));
    return {
      passed: failed.length === 0,
      statuses,
      ...(failed.length ? { reason: 'owned fixture child did not stop' } : {}),
    };
  });
  if (typeof waitForPortFree === 'function') await check('gateway_port_released', waitForPortFree);
  else steps.gateway_port_released = { passed: true, reason: 'not-created' };
  if (typeof permissions?.cleanup === 'function') await check('permission_role_removed', () => permissions.cleanup());
  else steps.permission_role_removed = { passed: true, reason: 'not-created' };
  if (typeof targets?.cleanup === 'function') await check('owned_targets_removed', () => targets.cleanup());
  else steps.owned_targets_removed = { passed: true, reason: 'not-created' };
  return markCleanupFailure({ passed: Object.values(steps).every((step) => step.passed), steps });
}

/** Cleanup is always completed before report assembly or publication is attempted. */
export async function publishAfterCleanup({ cleanup, buildReport, writeReport }) {
  let fixtureCleanup;
  try {
    fixtureCleanup = markCleanupFailure(await cleanup());
  } catch (error) {
    process.exitCode = 1;
    fixtureCleanup = { passed: false, steps: { cleanup: { passed: false, reason: boundedError(error) } } };
  }

  let report;
  try {
    report = await buildReport(fixtureCleanup);
  } catch (error) {
    process.exitCode = 1;
    report = {
      passed: false,
      failure: { stage: 'report', reason: 'report-build-failed', message: boundedError(error) },
      fixture_cleanup: fixtureCleanup,
    };
  }
  if (fixtureCleanup?.passed === false) {
    process.exitCode = 1;
    report = report && typeof report === 'object'
      ? { ...report, passed: false }
      : { passed: false, fixture_cleanup: fixtureCleanup };
  }
  try { await writeReport(report); }
  catch (error) { process.exitCode = 1; console.error(`fixture report write failed: ${boundedError(error)}`); }
  return report;
}

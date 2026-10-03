// Container entrypoint: copy only Git-visible regular files from read-only /repo.
import assert from 'node:assert/strict';
import { createHash } from 'node:crypto';
import { execFileSync } from 'node:child_process';
import { chmodSync, copyFileSync, existsSync, lstatSync, mkdirSync, readFileSync, rmSync, writeFileSync } from 'node:fs';
import { dirname, join, resolve, sep } from 'node:path';

assert.equal(process.platform, 'linux');
const source = '/work/source';
const output = '/output';
mkdirSync(source, { recursive: true }); mkdirSync(output, { recursive: true });
const sourceSHA = execFileSync('git', ['-C', '/repo', 'rev-parse', 'HEAD'], { encoding: 'utf8' }).trim();
const dirtyWorktree = execFileSync('git', ['-C', '/repo', 'status', '--porcelain'], { encoding: 'utf8' }).trim() !== '';
const files = execFileSync('git', ['-C', '/repo', 'ls-files', '--cached', '--others', '--exclude-standard', '-z'], { encoding: 'utf8' }).split('\0').filter(Boolean);
const manifest = [];
for (const relative of files) {
  const from = resolve('/repo', relative); const to = resolve(source, relative);
  assert.ok(from.startsWith(`/repo${sep}`) && to.startsWith(`${source}${sep}`));
  if (!existsSync(from) || !lstatSync(from).isFile()) continue;
  mkdirSync(dirname(to), { recursive: true }); copyFileSync(from, to); chmodSync(to, lstatSync(from).mode & 0o777);
  manifest.push({ path: relative, sha256: createHash('sha256').update(readFileSync(to)).digest('hex') });
}
const metadata = {
  image_id: process.env.F_LINUX_IMAGE_ID,
  platform: process.platform, architecture: process.arch,
  kernel: execFileSync('uname', ['-sr'], { encoding: 'utf8' }).trim(),
  environment: 'Docker Desktop Linux; amd64 userland under emulation on arm64 host',
  source_sha: sourceSHA,
  dirty_worktree: dirtyWorktree,
  source_manifest_sha256: createHash('sha256').update(JSON.stringify(manifest)).digest('hex'),
  commands: ['npm ci', 'npm run build', 'node scripts/tests/f_device_to_sql/run.mjs sqlite'],
  limits: ['container witness; no native Linux hardware, Windows, ARM deployment, embedded browser, LAN, PLC or SCADA acceptance'],
};

function retainedWorkFromWitness(witnessPath) {
  if (!existsSync(witnessPath)) return { path: null, error: 'device witness is missing' };
  let payload;
  try {
    payload = JSON.parse(readFileSync(witnessPath, 'utf8'));
  } catch (error) {
    return { path: null, error: `device witness is not valid JSON: ${error.message.split('\n')[0]}` };
  }
  const workPath = payload?.cleanup?.retained_work;
  if (typeof workPath !== 'string' || resolve(workPath) !== workPath || !/^\/tmp\/gw-f-[A-Za-z0-9_-]+$/.test(workPath)) {
    return { path: null, error: 'device witness cleanup.retained_work is not an owned /tmp/gw-f-* path' };
  }
  return { path: workPath };
}

function captureSQLiteFailureArtifacts(retainedWork) {
  if (!retainedWork.path) return { retained_work: null, error: retainedWork.error, files: [] };
  const workPath = retainedWork.path;
  const artifactNames = [];
  for (const name of ['gateway.db', 'destination.db', 'gateway.log', 'sim-a.log', 'sim-b.log']) {
    const from = join(workPath, name);
    if (!existsSync(from)) continue;
    try {
      copyFileSync(from, join(output, `linux-${name}`));
      artifactNames.push(`linux-${name}`);
    } catch (error) {
      artifactNames.push({ name: `linux-${name}`, error: error.message.split('\n')[0] });
    }
  }
  const gatewayDB = join(workPath, 'gateway.db');
  if (!existsSync(gatewayDB)) return { retained_work: workPath, files: artifactNames };
  const queries = {
    'linux-gateway-outbox.json': `SELECT effect_key, record_id, workspace_id, group_id, group_revision, entity_key,
      bucket_start, partition_key, connector_id, connector_revision, table_name, dedupe_capability,
      payload_digest, state, retry_count, next_retry_at, last_error_code, claim_owner,
      claim_expires_at, claim_epoch, committed_at, created_at, updated_at
      FROM wg_delivery_outbox ORDER BY group_id, bucket_start, effect_key;`,
    'linux-gateway-outbox-summary.json': `SELECT state, last_error_code, COUNT(*) AS count,
      MIN(bucket_start) AS first_bucket, MAX(bucket_start) AS last_bucket
      FROM wg_delivery_outbox GROUP BY state, last_error_code ORDER BY state, last_error_code;`,
    'linux-gateway-buckets.json': `SELECT group_id, group_revision, entity_key, bucket_start, kind,
      reason, record_id, members, created_at FROM wg_delivery_buckets
      ORDER BY group_id, bucket_start, entity_key;`,
    'linux-gateway-samples.json': `SELECT group_id, group_revision, member_key, observed_at,
      bucket_start, payload_digest, consumed, created_at FROM wg_delivery_samples
      ORDER BY group_id, bucket_start, member_key, observed_at;`,
    'linux-gateway-receipts.json': `SELECT effect_key, payload_digest, committed_at
      FROM wg_delivery_receipts ORDER BY effect_key;`,
  };
  for (const [name, sql] of Object.entries(queries)) {
    try {
      const raw = execFileSync('sqlite3', ['-cmd', '.timeout 15000', '-json', gatewayDB, sql], { encoding: 'utf8' }).trim();
      writeFileSync(join(output, name), `${JSON.stringify(raw ? JSON.parse(raw) : [])}\n`);
      artifactNames.push(name);
    } catch (error) {
      writeFileSync(join(output, name), `${JSON.stringify({ error: error.message.split('\n')[0] })}\n`);
      artifactNames.push(name);
    }
  }
  return { retained_work: workPath, files: artifactNames };
}

// A copied macOS witness must never be exported as a newly executed Linux run.
for (const name of ['device-to-sqlite-mixed.json', 'device-to-sqlite-step4.png']) {
  rmSync(join(source, 'docs/plans/studio-v2-write-groups/evidence-f', name), { force: true });
}
try {
  process.env.F_SOURCE_SHA = sourceSHA;
  process.env.F_DIRTY_WORKTREE = String(dirtyWorktree);
  execFileSync('npm', ['ci'], { cwd: join(source, 'frontend'), stdio: 'inherit' });
  execFileSync('npm', ['run', 'build'], { cwd: join(source, 'frontend'), stdio: 'inherit' });
  execFileSync('node', ['scripts/tests/f_device_to_sql/run.mjs', 'sqlite'], { cwd: source, stdio: 'inherit' });
  metadata.passed = true;
} catch (error) {
  metadata.passed = false; metadata.exit_code = error.status ?? 1;
  metadata.failure = 'Linux acceptance command failed; see container log and saved device witness';
  process.exitCode = 1;
} finally {
  const evidence = join(source, 'docs/plans/studio-v2-write-groups/evidence-f');
  for (const [from, to] of [
    ['device-to-sqlite-mixed.json', 'device-to-sqlite-mixed-linux.json'],
    ['device-to-sqlite-step4.png', 'device-to-sqlite-step4-linux.png'],
  ]) if (existsSync(join(evidence, from))) copyFileSync(join(evidence, from), join(output, to));
  const retainedWork = retainedWorkFromWitness(join(evidence, 'device-to-sqlite-mixed.json'));
  metadata.retained_work = retainedWork.path;
  if (metadata.passed === false) metadata.diagnostic_artifacts = captureSQLiteFailureArtifacts(retainedWork);
  writeFileSync(join(output, 'linux-environment.json'), JSON.stringify(metadata) + '\n');
  writeFileSync(join(output, 'linux-source-manifest.json'), JSON.stringify(manifest) + '\n');
}

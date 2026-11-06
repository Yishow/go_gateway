import { spawn, type ChildProcessByStdio } from 'node:child_process';
import { createHash, randomBytes } from 'node:crypto';
import { appendFile, cp, mkdtemp, readFile, readdir, rm, symlink, writeFile } from 'node:fs/promises';
import { createServer, type AddressInfo } from 'node:net';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';
import type { Readable } from 'node:stream';

type OwnedProcess = ChildProcessByStdio<null, Readable, Readable>;
export const repoRoot = process.cwd().endsWith('frontend') ? resolve(process.cwd(), '..') : process.cwd();
const digest = (bytes: Buffer) => createHash('sha256').update(bytes).digest('hex');

export async function runProcess(command: string, args: string[], cwd: string, env = process.env): Promise<string> {
  const child = spawn(command, args, { cwd, env, stdio: ['ignore', 'pipe', 'pipe'], windowsHide: true });
  const output: string[] = [];
  child.stdout.on('data', (chunk: Buffer) => output.push(chunk.toString()));
  child.stderr.on('data', (chunk: Buffer) => output.push(chunk.toString()));
  return new Promise((resolveRun, reject) => {
    child.once('error', reject);
    child.once('close', code => code === 0 ? resolveRun(output.join('')) : reject(new Error(`${command} exited ${code}\n${output.join('').slice(-8000)}`)));
  });
}

export async function embeddedPort(): Promise<number> {
  const probe = createServer();
  await new Promise<void>((done, reject) => { probe.once('error', reject); probe.listen(0, '127.0.0.1', done); });
  const port = (probe.address() as AddressInfo).port;
  await new Promise<void>((done, reject) => probe.close(error => error ? reject(error) : done()));
  return port;
}

export interface EmbeddedBuild {
  root: string;
  binary: string;
  readiness: { path: string; body: string };
  assets: string[];
  identity: { sha: string; lockSHA256: string; go: string; node: string; npm: string; variant: string | null; assetSHA256: Record<string, string>; assetBytes: Record<string, number>; htmlSHA256: string; binarySHA256: string };
  dispose: () => Promise<void>;
}

/** Archive only tracked source. Never copy local dotenv, databases or generated assets. */
export async function buildEmbedded(options: { revision?: string; cacheProtocolVariant?: string } = {}): Promise<EmbeddedBuild> {
  const root = await mkdtemp(join(tmpdir(), 'gateway-embedded-build-'));
  try {
    const sha = (await runProcess('git', ['rev-parse', '--verify', `${options.revision ?? 'HEAD'}^{commit}`], repoRoot)).trim();
    await runProcess('git', ['archive', '--format=tar', '-o', join(root, 'source.tar'), sha], repoRoot);
    await runProcess('tar', ['-xf', 'source.tar'], root);
    await rm(join(root, 'source.tar'));
    const lock = await readFile(join(root, 'frontend/package-lock.json'));
    if (!lock.equals(await readFile(join(repoRoot, 'frontend/package-lock.json')))) throw new Error('Archived lockfile differs; shared node_modules reuse refused');
    // A controlled cache-protocol variant is NOT a before/after product benchmark.
    if (options.cacheProtocolVariant) await appendFile(join(root, 'frontend/src/main.tsx'), `\nObject.defineProperty(globalThis, '__embeddedCacheProtocolVersion', { value: ${JSON.stringify(options.cacheProtocolVariant)} });\n`);
    await symlink(join(repoRoot, 'frontend/node_modules'), join(root, 'frontend/node_modules'), 'junction');
    const cleanBuildEnv = { ...Object.fromEntries(Object.entries(process.env).filter(([key]) => !key.startsWith('VITE_'))), VITE_API_BASE_URL: '/api/v1' };
    await runProcess(process.platform === 'win32' ? 'npm.cmd' : 'npm', ['run', 'build'], join(root, 'frontend'), cleanBuildEnv);
    const staticRoot = join(root, 'cmd/test_ui/static');
    await rm(staticRoot, { recursive: true, force: true });
    await cp(join(root, 'frontend/dist'), staticRoot, { recursive: true });
    // Test-only proof served by the existing static handler, never imported by the app.
    const readiness = { path: `assets/embedded-test-identity-${randomBytes(24).toString('hex')}.txt`, body: `owned-build:${randomBytes(32).toString('hex')}` };
    await writeFile(join(staticRoot, readiness.path), readiness.body);
    const binary = join(root, process.platform === 'win32' ? 'test-ui.exe' : 'test-ui');
    await runProcess('go', ['build', '-buildvcs=false', '-o', binary, './cmd/test_ui'], root);
    const assets = (await readdir(join(staticRoot, 'assets'))).map(name => `assets/${name}`).sort();
    const assetSHA256: Record<string, string> = {};
    const assetBytes: Record<string, number> = {};
    for (const path of assets) {
      const bytes = await readFile(join(staticRoot, path));
      assetSHA256[path] = digest(bytes);
      assetBytes[path] = bytes.length;
    }
    return { root, binary, readiness, assets, identity: { sha, lockSHA256: digest(lock), variant: options.cacheProtocolVariant ?? null, assetSHA256, assetBytes,
      htmlSHA256: digest(await readFile(join(staticRoot, 'index.html'))), binarySHA256: digest(await readFile(binary)),
      go: (await runProcess('go', ['version'], root)).trim(), node: process.version,
      npm: (await runProcess(process.platform === 'win32' ? 'npm.cmd' : 'npm', ['--version'], root)).trim() },
    dispose: () => rm(root, { recursive: true, force: true }) };
  } catch (error) { await rm(root, { recursive: true, force: true }); throw error; }
}

/** Fresh disposable SQLite; loopback only; no inherited service/device configuration. */
export async function startEmbedded(build: EmbeddedBuild, requestedPort?: number) {
  if (requestedPort !== undefined && (!Number.isInteger(requestedPort) || requestedPort < 1 || requestedPort > 65535)) {
    throw new Error('Requested embedded port must be an integer from 1 to 65535');
  }
  const port = requestedPort ?? await embeddedPort();
  // Refuse an already occupied origin before spawning. Identity below also closes
  // the release-to-child-bind race without ever signalling an unrelated server.
  const reservation = createServer();
  await new Promise<void>((done, reject) => {
    reservation.once('error', reject);
    reservation.listen(port, '127.0.0.1', done);
  });
  await new Promise<void>((done, reject) => reservation.close(error => error ? reject(error) : done()));
  const runtimeRoot = await mkdtemp(join(tmpdir(), 'gateway-embedded-runtime-'));
  const baseURL = `http://127.0.0.1:${port}`;
  const logs: string[] = [];
  const child: OwnedProcess = spawn(build.binary, [], { cwd: runtimeRoot, env: {
    PATH: process.env.PATH, HOME: runtimeRoot, TMPDIR: runtimeRoot, SystemRoot: process.env.SystemRoot,
    PORT: String(port), HOST: '127.0.0.1', API_BASE_PATH: '/api/v1',
    GATEWAY_DB_PATH: join(runtimeRoot, 'datalink.db'), LOG_FILE: join(runtimeRoot, 'gateway-runtime.jsonl'),
    LOG_OUTPUT: 'console', AUTO_OPEN_BROWSER: 'false',
  }, stdio: ['ignore', 'pipe', 'pipe'], windowsHide: true });
  child.stdout.on('data', (chunk: Buffer) => logs.push(chunk.toString()));
  child.stderr.on('data', (chunk: Buffer) => logs.push(chunk.toString()));
  let spawnError: Error | undefined;
  child.once('error', error => { spawnError = error; });
  const stop = async () => {
    if (child.exitCode === null && child.signalCode === null && !spawnError) {
      await new Promise<void>((done, reject) => {
        const timer = setTimeout(() => reject(new Error('Owned server did not stop; preserve runtime for inspection')), 15_000);
        child.once('close', () => { clearTimeout(timer); done(); });
        child.kill('SIGTERM');
      });
    }
    await rm(runtimeRoot, { recursive: true, force: true });
  };
  try {
    const deadline = Date.now() + 120_000;
    let ready = false;
    while (Date.now() < deadline) {
      if (spawnError) throw spawnError;
      if (child.exitCode !== null || child.signalCode !== null) throw new Error(`Owned server exited\n${logs.join('').slice(-8000)}`);
      let response: Response | undefined;
      try { response = await fetch(`${baseURL}/${build.readiness.path}`, { signal: AbortSignal.timeout(2000), redirect: 'manual', cache: 'no-store' }); } catch { /* not listening yet */ }
      if (response) {
        await assertEmbeddedIdentity(response, build.readiness.body);
        if (child.exitCode !== null || child.signalCode !== null || spawnError) throw new Error('Owned server exited during readiness');
        ready = true;
        break;
      }
      await new Promise(done => setTimeout(done, 100));
    }
    if (!ready) throw new Error(`Owned server readiness timeout\n${logs.join('').slice(-8000)}`);
    return { baseURL, port, runtimeRoot, stop, logs: () => logs.join('') };
  } catch (error) { await stop(); throw error; }
}

/** Bound the response read: a different service must never turn into an accepted origin. */
export async function assertEmbeddedIdentity(response: Response, expected: string): Promise<void> {
  const reader = response.body?.getReader();
  if (response.status !== 200 || !reader) {
    await reader?.cancel();
    throw new Error('Embedded readiness identity mismatch');
  }
  const expectedBytes = Buffer.from(expected);
  const chunks: Uint8Array[] = [];
  let length = 0;
  try {
    while (true) {
      const { done, value } = await reader.read();
      if (done) break;
      length += value.length;
      if (length > expectedBytes.length) throw new Error('Embedded readiness identity mismatch');
      chunks.push(value);
    }
    if (!Buffer.concat(chunks).equals(expectedBytes)) throw new Error('Embedded readiness identity mismatch');
  } finally { await reader.cancel(); }
}

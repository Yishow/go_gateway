// @vitest-environment node
import { afterEach, describe, expect, it } from 'vitest';
import { createServer, type Server } from 'node:http';
import { mkdtemp, writeFile, rm } from 'node:fs/promises';
import { join } from 'node:path';
import { tmpdir } from 'node:os';
import type { AddressInfo } from 'node:net';
import { assertEmbeddedIdentity, startEmbedded, type EmbeddedBuild } from '../../e2e/helpers/embedded-server';

const roots: string[] = [];
const stubs: Server[] = [];
afterEach(async () => {
  for (const server of stubs.splice(0)) await new Promise<void>((done, reject) => server.close(error => error ? reject(error) : done()));
  for (const root of roots.splice(0)) await rm(root, { recursive: true, force: true });
});
async function nonListeningBuild() {
  const root = await mkdtemp(join(tmpdir(), 'embedded-ownership-test-'));
  roots.push(root);
  const binary = join(root, 'delayed-child');
  await writeFile(binary, '#!/usr/bin/env node\nsetTimeout(() => {}, 30000);\n', { mode: 0o700 });
  return { root, binary, readiness: { path: 'assets/expected.txt', body: 'owned-build:expected' }, identity: {}, assets: [], dispose: async () => {} } as unknown as EmbeddedBuild;
}

describe('owned embedded readiness', () => {
  it.each([0, -1, 65536, 1.5, Number.NaN, Number.POSITIVE_INFINITY])('rejects impossible requested port %s before spawning', async port => {
    await expect(startEmbedded(await nonListeningBuild(), port)).rejects.toThrow('Requested embedded port');
  });
  it('accepts only the exact bounded build marker, never HTML/another build/missing/oversized responses', async () => {
    await expect(assertEmbeddedIdentity(new Response('owned-build:expected'), 'owned-build:expected')).resolves.toBeUndefined();
    for (const response of [new Response('<html>other service</html>'), new Response('owned-build:other'), new Response('', { status: 404 }), new Response('x'.repeat(10000))]) {
      await expect(assertEmbeddedIdentity(response, 'owned-build:expected')).rejects.toThrow('identity mismatch');
    }
  });
  // POSIX executable fixture only; Windows launch behavior needs its own runner.
  it.skipIf(process.platform === 'win32')('rejects an impostor response after bind preflight and terminates only its spawned child', async () => {
    const build = await nonListeningBuild();
    await writeFile(build.binary, '#!/usr/bin/env node\nrequire("node:http").createServer((_req,res)=>res.end("other-build")).listen(Number(process.env.PORT),"127.0.0.1");\n', { mode: 0o700 });
    await expect(startEmbedded(build)).rejects.toThrow('identity mismatch');
  });
  it.skipIf(process.platform === 'win32')('rejects a redirect without requesting a valid marker at a different origin', async () => {
    let redirectedRequests = 0;
    const target = createServer((_request, response) => { redirectedRequests++; response.end('owned-build:expected'); });
    stubs.push(target);
    await new Promise<void>(done => target.listen(0, '127.0.0.1', done));
    const redirectPort = (target.address() as AddressInfo).port;
    const build = await nonListeningBuild();
    await writeFile(build.binary, `#!/usr/bin/env node\nrequire("node:http").createServer((_req,res)=>{res.writeHead(302,{Location:"http://127.0.0.1:${redirectPort}/valid-marker"});res.end();}).listen(Number(process.env.PORT),"127.0.0.1");\n`, { mode: 0o700 });
    await expect(startEmbedded(build)).rejects.toThrow('identity mismatch');
    expect(redirectedRequests).toBe(0);
    expect(target.listening).toBe(true);
  });
  it('rejects an occupied HTTP200 origin instead of returning it for future seed writes', async () => {
    const stub = createServer((_request, response) => { response.writeHead(200); response.end('unrelated local service'); });
    stubs.push(stub);
    await new Promise<void>(done => stub.listen(0, '127.0.0.1', done));
    const port = (stub.address() as AddressInfo).port;
    let returned: Awaited<ReturnType<typeof startEmbedded>> | undefined;
    let failure: unknown;
    try { returned = await startEmbedded(await nonListeningBuild(), port); } catch (error) { failure = error; }
    finally { if (returned) await returned.stop(); }
    expect(failure, 'must not hand out an unrelated origin').toBeInstanceOf(Error);
    expect(returned).toBeUndefined();
    expect((await fetch(`http://127.0.0.1:${port}`)).status).toBe(200);
  });
});

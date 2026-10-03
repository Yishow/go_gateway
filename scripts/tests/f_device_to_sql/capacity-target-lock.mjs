import assert from 'node:assert/strict';
import { spawn } from 'node:child_process';

// A real reserved SQLite write lock keeps the target metadata readable while
// production INSERTs wait. No row is inserted or changed by this process.
export async function holdTargetWrites(path) {
  assert.match(path, /^\/tmp\/gw-f-capacity-\d+\/a\.db\.away$/);
  const script = `import sqlite3,sys,urllib.parse
db=sqlite3.connect('file:'+urllib.parse.quote(sys.argv[1],safe='/')+'?mode=rw',uri=True,timeout=2)
db.execute('BEGIN IMMEDIATE')
print('READY',flush=True)
sys.stdin.buffer.read(1)
db.rollback()
db.close()
`;
  const child = spawn('python3', ['-u', '-c', script, path], { stdio: ['pipe', 'pipe', 'pipe'] });
  const exited = new Promise((resolve) => child.once('exit', resolve));
  const stopOwnedChild = () => child.kill('SIGTERM');
  process.once('exit', stopOwnedChild);
  let stderr = '';
  child.stderr.on('data', (value) => { stderr += value.toString(); });
  await new Promise((resolve, reject) => {
    const timer = setTimeout(() => reject(new Error(`owned SQLite lock was not ready: ${stderr}`)), 5000);
    child.once('error', (error) => { clearTimeout(timer); reject(error); });
    child.once('exit', (code) => { clearTimeout(timer); reject(new Error(`owned SQLite lock exited ${code}: ${stderr}`)); });
    child.stdout.once('data', (value) => {
      clearTimeout(timer);
      if (value.toString().trim() === 'READY') resolve(); else reject(new Error('owned SQLite lock returned an unexpected readiness message'));
    });
  });
  return {
    pid: child.pid, sql: 'BEGIN IMMEDIATE; ROLLBACK',
    close: async () => {
      child.stdin.end();
      const timer = setTimeout(stopOwnedChild, 3000);
      const code = await exited; clearTimeout(timer); process.off('exit', stopOwnedChild);
      assert.equal(code, 0, `owned SQLite lock cleanup: ${stderr}`);
    },
  };
}

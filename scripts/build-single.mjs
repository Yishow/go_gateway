// Shared by start.sh and start.ps1; all frontend assets are embedded by Go.
import { spawnSync } from 'node:child_process';
import { cpSync, existsSync, mkdirSync, readdirSync, renameSync, rmSync, statSync } from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const repoRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..');

export function runCommand(command, args, options = {}) {
  const result = spawnSync(command, args, {
    stdio: 'inherit',
    // npm is a .cmd wrapper on Windows; arguments here are fixed script names.
    shell: process.platform === 'win32' && command === 'npm',
    ...options,
  });
  if (result.error) throw result.error;
  if (result.status !== 0) throw new Error(`${command} ${args.join(' ')} failed (${result.status ?? result.signal})`);
}

export function resolveTarget(os = process.platform, arch = process.arch) {
  const targetOS = ({ mac: 'darwin', macos: 'darwin', win32: 'windows' })[os] ?? os;
  const targetArch = arch === 'x64' ? 'amd64' : arch;
  if (!['windows', 'darwin', 'linux'].includes(targetOS)) throw new Error(`不支援的平台: ${os}`);
  if (!['amd64', 'arm64'].includes(targetArch)) throw new Error(`不支援的架構: ${arch}`);
  return { targetOS, targetArch };
}

export function prepareFrontend({ root = repoRoot, run = runCommand } = {}) {
  const frontend = path.join(root, 'frontend');
  const dist = path.join(frontend, 'dist');
  const staticDir = path.join(root, 'cmd/test_ui/static');
  for (const input of ['package.json', 'package-lock.json']) {
    if (!existsSync(path.join(frontend, input))) throw new Error(`缺少 frontend/${input}`);
  }

  const hasCompiler = ['typescript', '@typescript/native'].some(
    (name) => existsSync(path.join(frontend, 'node_modules', name, 'bin/tsc')),
  );
  if (!existsSync(path.join(frontend, 'node_modules/vite/bin/vite.js')) || !hasCompiler) {
    run('npm', ['ci'], { cwd: frontend });
  }
  // A build that exits without emitting output must not reuse an old dist.
  rmSync(dist, { recursive: true, force: true });
  run('npm', ['run', 'build'], { cwd: frontend });
  const index = path.join(dist, 'index.html');
  if (!existsSync(index) || !statSync(index).isFile() || statSync(index).size === 0) {
    throw new Error('frontend/dist/index.html 不存在或為空，已停止後端建置');
  }

  mkdirSync(staticDir, { recursive: true });
  for (const entry of readdirSync(staticDir)) {
    if (entry !== 'embed-placeholder.txt') rmSync(path.join(staticDir, entry), { recursive: true, force: true });
  }
  for (const entry of readdirSync(dist)) {
    cpSync(path.join(dist, entry), path.join(staticDir, entry), { recursive: true });
  }
}

export function buildSingle({ root = repoRoot, os, arch, run = runCommand } = {}) {
  const { targetOS, targetArch } = resolveTarget(os, arch);
  prepareFrontend({ root, run });
  const bin = path.join(root, 'bin');
  mkdirSync(bin, { recursive: true });
  const output = path.join(bin, `gateway-${targetOS}-${targetArch}${targetOS === 'windows' ? '.exe' : ''}`);
  const pending = `${output}.pending-${process.pid}`;
  console.log(`建置單一執行檔: ${targetOS}/${targetArch}`);
  try {
    run('go', ['build', '-trimpath', '-ldflags', '-s -w', '-o', pending, './cmd/test_ui'], {
      cwd: root,
      env: { ...process.env, GOOS: targetOS, GOARCH: targetArch, CGO_ENABLED: '0' },
    });
    renameSync(pending, output);
  } finally {
    rmSync(pending, { force: true });
  }
  console.log(`完成: ${output}`);
  return output;
}

if (process.argv[1] && path.resolve(process.argv[1]) === fileURLToPath(import.meta.url)) {
  try {
    const options = {};
    const args = process.argv.slice(2);
    if (args.length === 1 && args[0] === '--prepare-dev') {
      prepareFrontend();
    } else {
      for (let i = 0; i < args.length; i += 2) {
        if (!['--os', '--arch'].includes(args[i]) || !args[i + 1]) throw new Error('用法: node scripts/build-single.mjs [--os windows|mac|linux] [--arch amd64|arm64]');
        options[args[i].slice(2)] = args[i + 1];
      }
      buildSingle(options);
    }
  } catch (error) {
    console.error(error.message);
    process.exitCode = 1;
  }
}

#!/usr/bin/env python3
"""Create/check recording E2E evidence; never operates a service or database."""
import argparse
import hashlib
import json
import subprocess
import sys
import uuid
from datetime import datetime, timezone
from pathlib import Path

SOURCE_STAGES = (
    'device-create', 'readiness', 'range-create-edit', 'range-activation',
    'memory-views-values', 'source-failure-recovery',
)
RECORDING_STAGES = (
    'record-create', 'schema-consent', 'test-write', 'record-start',
    'sql-readback', 'reload', 'record-stop',
)
PROTOCOLS = ('mc_3e', 'modbus_tcp')
DATABASES = ('sqlite', 'postgres', 'mysql', 'mariadb', 'sqlserver')


def digest(data):
    return hashlib.sha256(data).hexdigest()


def encoded(value):
    return json.dumps(value, sort_keys=True, separators=(',', ':')).encode()


def git(repo, *args):
    return subprocess.check_output(
        ['git', '--no-optional-locks', '-C', str(repo), *args], stderr=subprocess.PIPE,
    )


def snapshot(repo):
    head = git(repo, 'rev-parse', 'HEAD').decode().strip()
    diff = git(repo, 'diff', '--binary', 'HEAD', '--')
    tracked = {}
    names = git(repo, 'ls-files', '--cached', '-z')
    for raw in names.split(b'\0'):
        if raw:
            name = raw.decode()
            file = repo / name
            tracked[name] = (
                {'sha256': digest(file.read_bytes()), 'executable': file.stat().st_mode & 0o111}
                if file.is_file() else 'not-file'
            )
            if file.is_symlink():
                tracked[name] = [tracked[name], str(file.readlink())]
    names = git(repo, 'ls-files', '--others', '--exclude-standard', '-z')
    untracked = {}
    untracked_executable = {}
    untracked_symlinks = {}
    for raw in names.split(b'\0'):
        if raw:
            name = raw.decode()
            file = repo / name
            untracked[name] = digest(file.read_bytes()) if file.is_file() else 'not-file'
            if file.is_file():
                untracked_executable[name] = file.stat().st_mode & 0o111
            if file.is_symlink():
                untracked_symlinks[name] = str(file.readlink())
    return {
        'head': head, 'diff_sha256': digest(diff),
        'tracked_sha256': digest(encoded(tracked)), 'untracked_sha256': untracked,
        'untracked_executable': untracked_executable, 'untracked_symlinks': untracked_symlinks,
    }


def scope_valid(scope):
    if not isinstance(scope, dict) or set(scope) != {'protocols', 'databases'}:
        raise ValueError('invalid scope')
    for name, allowed in [('protocols', PROTOCOLS), ('databases', DATABASES)]:
        values = scope[name]
        if not isinstance(values, list) or not values or not all(isinstance(x, str) for x in values):
            raise ValueError('invalid ' + name)
        if len(values) != len(set(values)) or not set(values) <= set(allowed):
            raise ValueError('unsupported or duplicate ' + name)


def required(scope):
    scope_valid(scope)
    result = ['build-provenance', 'home-stability', 'cleanup']
    for protocol in scope['protocols']:
        result.extend(f'source/{protocol}/{stage}' for stage in SOURCE_STAGES)
        for database in scope['databases']:
            result.extend(f'recording/{protocol}/{database}/{stage}' for stage in RECORDING_STAGES)
    return result


def load(path):
    return json.loads(path.read_text())


def write(path, value):
    path.write_text(json.dumps(value, ensure_ascii=False, indent=2) + '\n')


def initialize(args):
    repo = Path(args.repo).expanduser().resolve()
    if Path(git(repo, 'rev-parse', '--show-toplevel').decode().strip()).resolve() != repo:
        raise ValueError('--repo must be the repository root')
    root = Path(args.out).expanduser().resolve()
    if root == repo or repo in root.parents:
        raise ValueError('run directory must be outside the candidate repo')
    scope = {'protocols': args.protocol or list(PROTOCOLS), 'databases': args.database or ['sqlite']}
    ids = required(scope)
    source = snapshot(repo)
    root.mkdir(parents=True, exist_ok=False)
    manifest = {
        'schema_version': 1, 'run_id': str(uuid.uuid4()),
        'created_at': datetime.now(timezone.utc).isoformat(), 'repo': str(repo),
        'scope': scope, 'source': source,
    }
    write(root / 'manifest.json', manifest)
    write(root / '.recording-e2e-owner.json', {
        'run_id': manifest['run_id'], 'scope_sha256': digest(encoded(scope)),
        'source_sha256': digest(encoded(source)),
    })
    write(root / 'checks.json', [
        {'id': name, 'status': 'NOT_RUN', 'detail': '', 'evidence': []} for name in ids
    ])
    print(json.dumps({'run': str(root), 'required_checks': len(ids), 'product_e2e': 'NOT_RUN'}))
    return 0


def check(args):
    root = Path(args.run).expanduser().resolve()
    manifest = load(root / 'manifest.json')
    owner = load(root / '.recording-e2e-owner.json')
    if manifest.get('schema_version') != 1:
        raise ValueError('unsupported manifest version')
    if manifest['run_id'] != owner['run_id']:
        raise ValueError('run identity mismatch')
    if digest(encoded(manifest['scope'])) != owner['scope_sha256']:
        raise ValueError('declared scope changed; retain this run and start a new one')
    if digest(encoded(manifest['source'])) != owner['source_sha256']:
        raise ValueError('recorded source identity changed')
    expected = required(manifest['scope'])
    snapshot_fields = {'tracked_sha256', 'untracked_executable', 'untracked_symlinks'}
    if not isinstance(manifest['source'], dict) or not snapshot_fields <= manifest['source'].keys():
        raise ValueError('source snapshot predates this tool; retain this run and start a new run')
    if snapshot(Path(manifest['repo'])) != manifest['source']:
        raise ValueError('candidate source changed; rebuild and start a new run')
    rows = load(root / 'checks.json')
    if not isinstance(rows, list) or not all(isinstance(row, dict) for row in rows):
        raise ValueError('checks must be a list of objects')
    ids = [row.get('id') for row in rows]
    if not all(isinstance(name, str) for name in ids):
        raise ValueError('check IDs must be strings')
    if len(ids) != len(set(ids)) or set(ids) != set(expected):
        missing = sorted(set(expected) - set(ids))
        raise ValueError('missing, extra or duplicate checks; missing=' + ','.join(missing))
    failures = []
    for row in rows:
        name = row['id']
        if row.get('status') != 'PASS':
            failures.append(name + ': ' + str(row.get('status', 'missing status')))
            continue
        if not isinstance(row.get('detail'), str) or not row['detail'].strip():
            failures.append(name + ': missing observation detail')
        refs = row.get('evidence')
        if not isinstance(refs, list) or not refs:
            failures.append(name + ': missing evidence')
            continue
        for ref in refs:
            if not isinstance(ref, str) or not ref.strip():
                failures.append(name + ': invalid evidence reference')
                continue
            relative = Path(ref)
            file = (root / relative).resolve()
            if relative.is_absolute() or '..' in relative.parts or root not in file.parents:
                failures.append(name + ': evidence escapes run directory')
            elif not file.is_file() or file.stat().st_size == 0:
                failures.append(name + ': missing or empty evidence ' + ref)
    result = {
        'run_id': manifest['run_id'], 'evidence_gate': 'FAIL' if failures else 'PASS',
        'failures': failures,
        'limit': 'Structural/source checks only; agent must review actual UI/API/protocol/SQL evidence.',
    }
    print(json.dumps(result, ensure_ascii=False, indent=2))
    return 1 if failures else 0


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    actions = parser.add_subparsers(dest='action', required=True)
    init = actions.add_parser('init', help='create a new NOT_RUN checklist')
    init.add_argument('--repo', required=True)
    init.add_argument('--out', required=True)
    init.add_argument('--protocol', action='append', choices=PROTOCOLS)
    init.add_argument('--database', action='append', choices=DATABASES)
    verify = actions.add_parser('check', help='check completeness and source identity, not business correctness')
    verify.add_argument('--run', required=True)
    args = parser.parse_args()
    try:
        return initialize(args) if args.action == 'init' else check(args)
    except (OSError, ValueError, KeyError, TypeError, subprocess.CalledProcessError) as error:
        print(json.dumps({'evidence_gate': 'ERROR', 'reason': str(error)}, ensure_ascii=False), file=sys.stderr)
        return 1


if __name__ == '__main__':
    sys.exit(main())

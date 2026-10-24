#!/usr/bin/env python3
"""Opt-in native smoke for Bash dev reload using a disposable source copy/DB.

Prerequisites: frontend dependencies and bin/tools/air prepared by start.sh.
Runs Go/frontend compilation, real local HTTP/WebSocket probes, and cleanup.
"""
import os, pathlib, shutil, signal, socket, subprocess, tempfile, time, urllib.request
root=pathlib.Path(__file__).resolve().parents[2]
opener=urllib.request.build_opener(urllib.request.ProxyHandler({}))
def free_port():
    with socket.socket() as sock:
        sock.bind(('127.0.0.1',0)); return sock.getsockname()[1]
def listener(port):
    p=subprocess.run(['lsof','-tiTCP:'+str(port),'-sTCP:LISTEN'],capture_output=True,text=True)
    return p.stdout.strip()
def wait_for(fn, seconds, message):
    deadline=time.monotonic()+seconds
    while time.monotonic()<deadline:
        v=fn()
        if v: return v
        time.sleep(.2)
    raise RuntimeError(message)
with tempfile.TemporaryDirectory(prefix='gateway dev smoke ') as temp:
    repo=pathlib.Path(temp)/'repo'; repo.mkdir()
    for name in ['cmd','internal','lib','frontend']:
        shutil.copytree(root/name,repo/name,ignore=shutil.ignore_patterns('node_modules','dist','static','.git'))
    shutil.copytree(root/'docs/swagger',repo/'docs/swagger')
    for name in ['go.mod','go.sum','.air.toml','start.sh']: shutil.copy2(root/name,repo/name)
    (repo/'scripts').mkdir()
    for name in ['build-single.mjs','load-env.sh','start-process-utils.sh']: shutil.copy2(root/'scripts'/name,repo/'scripts'/name)
    (repo/'cmd/test_ui/static').mkdir()
    shutil.copy2(root/'cmd/test_ui/static/embed-placeholder.txt',repo/'cmd/test_ui/static/embed-placeholder.txt')
    (repo/'frontend/node_modules').symlink_to(root/'frontend/node_modules',target_is_directory=True)
    (repo/'bin/tools').mkdir(parents=True)
    (repo/'bin/tools/air').symlink_to(root/'bin/tools/air')
    port,front=free_port(),free_port()
    env={**os.environ,'HOST':'127.0.0.1','PORT':str(port),'FRONTEND_DEV_HOST':'127.0.0.1','FRONTEND_DEV_PORT':str(front),'GATEWAY_DB_PATH':str(repo/'smoke.db'),'AUTO_OPEN_BROWSER':'false','LOG_OUTPUT':'console','GOOS':'windows','GOARCH':'amd64','CGO_ENABLED':'1'}
    log=pathlib.Path(temp)/'launcher.log'
    with log.open('w') as out: proc=subprocess.Popen(['/bin/bash',str(repo/'start.sh'),'--dev'],cwd='/tmp',env=env,stdout=out,stderr=subprocess.STDOUT)
    try:
        def ready():
            if proc.poll() is not None: raise RuntimeError('launcher exited: '+str(proc.returncode))
            try:
                with opener.open(f'http://127.0.0.1:{front}/api/v1/datalink/health',timeout=1) as r: return r.status==200
            except OSError: return False
        wait_for(ready,120,'dev startup timeout')
        assert (repo/'cmd/test_ui/static/index.html').stat().st_size>0
        print('fresh bootstrap + Vite API proxy + host Go target: PASS',flush=True)
        old=listener(port)
        lib=repo/'lib/hsllogic/address_parser.go'
        original=lib.read_bytes(); lib.write_bytes(original+b'\n// launcher smoke reload\n')
        wait_for(lambda: listener(port) and listener(port)!=old,60,'lib reload did not replace backend PID')
        wait_for(ready,30,'backend not ready after lib reload')
        print('lib Go edit automatically replaces backend PID: PASS',flush=True)
        lib.write_bytes(original+b'\nthis is an intentional compile failure\n')
        wait_for(lambda: 'failed to build' in log.read_text().lower() or 'syntax error' in log.read_text().lower(),30,'compile failure not visible')
        wait_for(lambda: not listener(port),30,'old backend still running after failed compile')
        lib.write_bytes(original)
        wait_for(ready,60,'saving fixed source did not recover backend')
        print('compile error stops stale backend; corrected source recovers: PASS',flush=True)
        js=r'''
const fs = require('node:fs');
(async () => {
 const base=process.argv[1], app=process.argv[2];
 const client=await (await fetch(base+'/@vite/client')).text();
 const token=client.match(/const wsToken = "([^"]+)"/)[1];
 await (await fetch(base+'/src/App.tsx')).text();
 const ws=new WebSocket(base.replace('http','ws')+'/?token='+token,'vite-hmr');
 const timeout=setTimeout(()=>{console.error('HMR timeout');process.exit(1)},15000);
 ws.addEventListener('message',event=>{
  const data=JSON.parse(event.data);
  if(data.type==='connected') setTimeout(()=>fs.appendFileSync(app,'\n// launcher HMR smoke\n'),300);
  if(data.type==='update' && data.updates.some(u=>u.path.includes('App.tsx'))) {
   clearTimeout(timeout);ws.close();console.log('Vite HMR update over WebSocket: PASS');
  }
 });
})().catch(e=>{console.error(e.message);process.exit(1)});
'''
        subprocess.run(['node','-e',js,f'http://127.0.0.1:{front}',str(repo/'frontend/src/App.tsx')],check=True,timeout=25)
        proc.send_signal(signal.SIGTERM);proc.wait(timeout=20)
        assert proc.returncode==143,proc.returncode
        assert not listener(port) and not listener(front)
        assert not (repo/f'bin/tmp/start-{port}.lock').exists()
        print('TERM owned process cleanup + ports/lock released: PASS',flush=True)
    except BaseException:
        print(log.read_text()[-5500:],flush=True)
        raise
    finally:
        if proc.poll() is None:
            proc.send_signal(signal.SIGTERM)
            try:proc.wait(timeout=20)
            except subprocess.TimeoutExpired:proc.kill();proc.wait()

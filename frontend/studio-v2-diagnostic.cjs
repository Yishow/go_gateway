const { spawnSync } = require('node:child_process');
const { randomBytes } = require('node:crypto');

const session = `gwdiag-${randomBytes(6).toString('hex')}`;
const url = process.env.STUDIO_DIAGNOSTIC_URL || 'http://127.0.0.1:4174/studio/v2';
const screenshot = process.env.STUDIO_DIAGNOSTIC_SCREENSHOT || '/private/tmp/studio-v2-diagnostic.png';

function browser(...args) {
  const command = spawnSync('agent-browser', ['--session', session, '--headed', 'false', '--json', ...args], {
    encoding: 'utf8',
    timeout: 40000,
  });
  if (command.error) throw command.error;
  let response;
  try {
    response = JSON.parse(command.stdout);
  } catch {
    throw new Error(command.stderr.trim() || `agent-browser returned invalid JSON (exit ${command.status})`);
  }
  if (command.status !== 0 || !response.success) {
    throw new Error(response.error || command.stderr.trim() || 'agent-browser command failed');
  }
  return response.data;
}

const inspectPage = `(() => {
  const count = (selector) => document.querySelectorAll(selector).length;
  const countText = (text) => Array.from(document.querySelectorAll('body *')).filter((element) => {
    const contains = (node) => (node.innerText || '').includes(text);
    return contains(element) && !Array.from(element.children).some(contains);
  }).length;
  return {
    title: document.title,
    url: location.href,
    step1: count('[data-testid="step1-device-view"]'),
    inputHost: count('[data-testid="input-host"]'),
    runTest: count('[data-testid="run-test-button"]'),
    signin: countText('登入'),
    notFound: countText('404'),
    testids: Array.from(document.querySelectorAll('[data-testid]'))
      .slice(0, 40).map((element) => element.getAttribute('data-testid')).filter(Boolean),
    body: (document.body?.innerText || '').replace(/\\s+/g, ' ').slice(0, 1800),
  };
})()`;

try {
  browser('open', url);
  browser('wait', '3000');
  const { result } = browser('eval', '-b', Buffer.from(inspectPage).toString('base64'));
  console.log('TITLE=', result.title);
  console.log('URL=', result.url);
  console.log('HAS_step1=', result.step1);
  console.log('HAS_input_host=', result.inputHost);
  console.log('HAS_run_test=', result.runTest);
  console.log('HAS_signin=', result.signin);
  console.log('HAS_404=', result.notFound);
  console.log('TESTID_HEAD=', result.testids.join(','));
  console.log('BODY=', result.body);
  console.log('BROWSER_CONSOLE=', JSON.stringify(browser('console')));
  console.log('BROWSER_PAGE_ERRORS=', JSON.stringify(browser('errors')));
  browser('screenshot', screenshot, '--full');
} catch (error) {
  console.error('DIAGNOSTIC_ERROR=', error.message);
  process.exitCode = 1;
} finally {
  try {
    browser('close');
  } catch (error) {
    console.error('BROWSER_CLOSE_ERROR=', error.message);
    process.exitCode = 1;
  }
}

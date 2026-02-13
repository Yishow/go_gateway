/* eslint-disable no-console */
const isDev = import.meta.env.DEV;

function output(method: 'log' | 'warn' | 'error', ...args: unknown[]) {
  if (!isDev) {
    return;
  }
  console[method](...args);
}

export const logger = {
  log: (...args: unknown[]) => output('log', ...args),
  warn: (...args: unknown[]) => output('warn', ...args),
  error: (...args: unknown[]) => output('error', ...args),
};


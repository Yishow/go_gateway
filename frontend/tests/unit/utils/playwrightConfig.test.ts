// @vitest-environment node
import { expect, it } from 'vitest';
import config from '../../../playwright.config';

it('serves isolated E2E HTML fixtures without the bundled app entrypoint', () => {
  const server = Array.isArray(config.webServer) ? config.webServer[0] : config.webServer;
  expect(server?.env?.VITE_BUNDLED_DEV).toBe('false');
});

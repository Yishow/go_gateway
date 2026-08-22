// @vitest-environment node

import { mkdtempSync, rmSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { afterEach, describe, expect, it } from "vitest";
import type { ConfigEnv, UserConfig } from "vite";
import viteConfig from "../../../vite.config";

const buildConfig = viteConfig as (env: ConfigEnv) => UserConfig;

const originalPort = process.env.PORT;
const originalProxyTarget = process.env.VITE_API_PROXY_TARGET;
const originalDevPort = process.env.VITE_DEV_PORT;
const originalWorkingDirectory = process.cwd();

function getApiProxyTarget(config: UserConfig): string | undefined {
  const proxy = config.server?.proxy;
  if (!proxy || Array.isArray(proxy)) {
    return undefined;
  }

  const apiProxy = proxy["/api"];
  if (!apiProxy) {
    return undefined;
  }

  if (typeof apiProxy === "string") {
    return apiProxy;
  }

  const target = apiProxy.target;
  return typeof target === "string" ? target : target?.toString();
}

function buildConfigWithoutAmbientEnv(env: ConfigEnv): UserConfig {
  const previousWorkingDirectory = process.cwd();
  const isolatedWorkingDirectory = mkdtempSync(join(tmpdir(), "go-gateway-vite-config-"));

  try {
    process.chdir(isolatedWorkingDirectory);
    return buildConfig(env);
  } finally {
    process.chdir(previousWorkingDirectory);
    rmSync(isolatedWorkingDirectory, { force: true, recursive: true });
  }
}

describe("vite proxy target", () => {
  afterEach(() => {
    if (originalPort === undefined) {
      delete process.env.PORT;
    } else {
      process.env.PORT = originalPort;
    }

    if (originalProxyTarget === undefined) {
      delete process.env.VITE_API_PROXY_TARGET;
    } else {
      process.env.VITE_API_PROXY_TARGET = originalProxyTarget;
    }

    if (originalDevPort === undefined) {
      delete process.env.VITE_DEV_PORT;
    } else {
      process.env.VITE_DEV_PORT = originalDevPort;
    }

    process.chdir(originalWorkingDirectory);
  });

  it("uses PORT when VITE_API_PROXY_TARGET is not set", () => {
    process.env.PORT = "3333";
    delete process.env.VITE_API_PROXY_TARGET;

    const config = buildConfigWithoutAmbientEnv({ command: "serve", mode: "test" });

    expect(getApiProxyTarget(config)).toBe("http://127.0.0.1:3333");
  });

  it("prefers VITE_API_PROXY_TARGET over PORT", () => {
    process.env.PORT = "3333";
    process.env.VITE_API_PROXY_TARGET = "http://127.0.0.1:9090";

    const config = buildConfigWithoutAmbientEnv({ command: "serve", mode: "test" });

    expect(getApiProxyTarget(config)).toBe("http://127.0.0.1:9090");
  });

  it("uses VITE_DEV_PORT for the frontend dev server port", () => {
    process.env.VITE_DEV_PORT = "4173";

    const config = buildConfigWithoutAmbientEnv({ command: "serve", mode: "test" });

    expect(config.server?.port).toBe(4173);
  });
});

import { useSyncExternalStore } from 'react';
import { gatewayAdapter, type ExpertDraft, type QuickDraft } from './gatewayAdapter';

const QUICK_CONFIG_KEYS = ['host', 'port', 'unitID', 'station', 'route', 'auth'] as const;

export interface GatewayQuickWizardDraft {
  sourcePlan: string;
  tagName: string;
  securityReviewed: boolean;
  validationConfirmed: boolean;
  submitConfirmed: boolean;
}

export interface GatewayExpertRouteDraft {
  id: string;
  method: 'GET' | 'POST' | 'PUT' | 'DELETE';
  path: string;
  target: string;
  enabled: boolean;
}

export interface GatewayExpertPluginDraft {
  requestId: boolean;
  cors: boolean;
  rateLimit: boolean;
}

export type ExpertQuickCompatibility = 'compatible' | 'unsupported' | 'invalid';

export interface GatewayDraftState {
  quickDraft: QuickDraft;
  quickWizardDraft: GatewayQuickWizardDraft;
  quickStep: number;
  expertDraft: ExpertDraft;
  expertBuilderRoutes: GatewayExpertRouteDraft[];
  expertBuilderPlugins: GatewayExpertPluginDraft;
  expertQuickCompatibility: ExpertQuickCompatibility;
  expertUnsupportedKeys: string[];
}

const DEFAULT_QUICK_DRAFT: QuickDraft = {
  protocol: 'modbus-tcp',
  host: '192.168.1.100',
  port: 502,
  unitID: 1,
  route: '/api/v1/data',
  auth: false,
};

const DEFAULT_QUICK_WIZARD_DRAFT: GatewayQuickWizardDraft = {
  sourcePlan: 'line-a-source',
  tagName: '',
  securityReviewed: false,
  validationConfirmed: false,
  submitConfirmed: false,
};

const DEFAULT_EXPERT_BUILDER_ROUTES: GatewayExpertRouteDraft[] = [
  {
    id: 'route-1',
    method: 'GET',
    path: '/api/v1/metrics',
    target: 'http://collector.internal:9000/metrics',
    enabled: true,
  },
  {
    id: 'route-2',
    method: 'POST',
    path: '/api/v1/write',
    target: 'http://collector.internal:9000/write',
    enabled: true,
  },
];

const DEFAULT_EXPERT_BUILDER_PLUGINS: GatewayExpertPluginDraft = {
  requestId: true,
  cors: true,
  rateLimit: false,
};

function deriveExpertDraftFromQuickDraft(quickDraft: QuickDraft): ExpertDraft {
  const payload = gatewayAdapter.quickToPayload(quickDraft);
  return {
    protocol: payload.protocol,
    configJson: JSON.stringify(payload.config, null, 2),
  };
}

function extractQuickDraftFromPayload(payload: { protocol: string; config: Record<string, unknown> }, fallback: QuickDraft): QuickDraft {
  const quickDraft: QuickDraft = {
    ...fallback,
    protocol: payload.protocol || fallback.protocol,
  };

  QUICK_CONFIG_KEYS.forEach((key) => {
    const value = payload.config[key];
    if (key === 'host' || key === 'route') {
      if (typeof value === 'string') quickDraft[key] = value;
    } else if (key === 'port' || key === 'unitID' || key === 'station') {
      if (typeof value === 'number') quickDraft[key] = value;
    } else if (key === 'auth') {
      if (typeof value === 'boolean') quickDraft[key] = value;
    }
  });

  return quickDraft;
}

function resolveExpertCompatibility(expertDraft: ExpertDraft, fallbackQuickDraft: QuickDraft): {
  compatibility: ExpertQuickCompatibility;
  unsupportedKeys: string[];
  quickDraft: QuickDraft | null;
} {
  try {
    const payload = gatewayAdapter.expertToPayload(expertDraft);
    if (!payload.config || typeof payload.config !== 'object' || Array.isArray(payload.config)) {
      return { compatibility: 'invalid', unsupportedKeys: [], quickDraft: null };
    }

    const config = payload.config as Record<string, unknown>;
    const unsupportedKeys = Object.keys(config).filter((key) => {
      if (!QUICK_CONFIG_KEYS.includes(key as (typeof QUICK_CONFIG_KEYS)[number])) {
        return true;
      }
      const value = config[key];
      return typeof value === 'object' || Array.isArray(value);
    });

    if (unsupportedKeys.length > 0) {
      return { compatibility: 'unsupported', unsupportedKeys, quickDraft: null };
    }

    return {
      compatibility: 'compatible',
      unsupportedKeys: [],
      quickDraft: extractQuickDraftFromPayload({ protocol: payload.protocol, config }, fallbackQuickDraft),
    };
  } catch {
    return { compatibility: 'invalid', unsupportedKeys: [], quickDraft: null };
  }
}

function createInitialState(): GatewayDraftState {
  return {
    quickDraft: { ...DEFAULT_QUICK_DRAFT },
    quickWizardDraft: { ...DEFAULT_QUICK_WIZARD_DRAFT },
    quickStep: 1,
    expertDraft: deriveExpertDraftFromQuickDraft(DEFAULT_QUICK_DRAFT),
    expertBuilderRoutes: DEFAULT_EXPERT_BUILDER_ROUTES.map((route) => ({ ...route })),
    expertBuilderPlugins: { ...DEFAULT_EXPERT_BUILDER_PLUGINS },
    expertQuickCompatibility: 'compatible',
    expertUnsupportedKeys: [],
  };
}

let state = createInitialState();
const listeners = new Set<() => void>();

function emitStoreChange() {
  listeners.forEach((listener) => listener());
}

function setState(updater: (current: GatewayDraftState) => GatewayDraftState) {
  const nextState = updater(state);
  if (Object.is(nextState, state)) return;
  state = nextState;
  emitStoreChange();
}

export function useGatewayDraftStore<T>(selector: (current: GatewayDraftState) => T): T {
  return useSyncExternalStore(
    (listener) => {
      listeners.add(listener);
      return () => listeners.delete(listener);
    },
    () => selector(state),
    () => selector(state),
  );
}

export function getGatewayDraftState(): GatewayDraftState {
  return state;
}

export function resetGatewayDraftStore() {
  state = createInitialState();
  emitStoreChange();
}

export function updateGatewayQuickDraft(patch: Partial<QuickDraft>) {
  setState((current) => {
    const nextQuickDraft = { ...current.quickDraft, ...patch };
    return {
      ...current,
      quickDraft: nextQuickDraft,
      expertDraft: deriveExpertDraftFromQuickDraft(nextQuickDraft),
      expertQuickCompatibility: 'compatible',
      expertUnsupportedKeys: [],
    };
  });
}

export function updateGatewayQuickWizardDraft(patch: Partial<GatewayQuickWizardDraft>) {
  setState((current) => ({
    ...current,
    quickWizardDraft: { ...current.quickWizardDraft, ...patch },
  }));
}

export function setGatewayQuickStep(step: number) {
  const safeStep = Math.min(Math.max(step, 1), 5);
  setState((current) => ({ ...current, quickStep: safeStep }));
}

export function updateGatewayExpertDraft(patch: Partial<ExpertDraft>) {
  setState((current) => {
    const nextExpertDraft: ExpertDraft = {
      protocol: patch.protocol ?? current.expertDraft.protocol,
      configJson: patch.configJson ?? current.expertDraft.configJson,
    };
    const compatibilityResult = resolveExpertCompatibility(nextExpertDraft, current.quickDraft);

    return {
      ...current,
      expertDraft: nextExpertDraft,
      quickDraft: compatibilityResult.quickDraft ?? current.quickDraft,
      expertQuickCompatibility: compatibilityResult.compatibility,
      expertUnsupportedKeys: compatibilityResult.unsupportedKeys,
    };
  });
}

export function updateGatewayExpertBuilderRoutes(routes: GatewayExpertRouteDraft[]) {
  setState((current) => ({
    ...current,
    expertBuilderRoutes: routes,
  }));
}

export function updateGatewayExpertBuilderPlugins(patch: Partial<GatewayExpertPluginDraft>) {
  setState((current) => ({
    ...current,
    expertBuilderPlugins: {
      ...current.expertBuilderPlugins,
      ...patch,
    },
  }));
}

export function confirmGatewayExpertDowngradeToQuick(): boolean {
  let didDowngrade = false;
  setState((current) => {
    if (current.expertQuickCompatibility !== 'unsupported') {
      return current;
    }

    const compatibilityResult = resolveExpertCompatibility(current.expertDraft, current.quickDraft);
    if (compatibilityResult.compatibility !== 'unsupported') {
      return {
        ...current,
        quickDraft: compatibilityResult.quickDraft ?? current.quickDraft,
        expertQuickCompatibility: compatibilityResult.compatibility,
        expertUnsupportedKeys: compatibilityResult.unsupportedKeys,
      };
    }

    let nextQuickDraft = { ...current.quickDraft };
    try {
      const payload = gatewayAdapter.expertToPayload(current.expertDraft);
      const config = payload.config as Record<string, unknown>;
      nextQuickDraft = extractQuickDraftFromPayload({ protocol: payload.protocol, config }, current.quickDraft);
    } catch {
      return current;
    }

    didDowngrade = true;
    return {
      ...current,
      quickDraft: nextQuickDraft,
      expertDraft: deriveExpertDraftFromQuickDraft(nextQuickDraft),
      expertQuickCompatibility: 'compatible',
      expertUnsupportedKeys: [],
    };
  });

  return didDowngrade;
}
